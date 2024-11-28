package pdfextract

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

var (
	fontRegex    = regexp.MustCompile(`/F\d+ \d+ Tf`)
	textRegex    = regexp.MustCompile(`\(([^)]+)\)\s*Tj`)
	posRegex     = regexp.MustCompile(`(\d+\.?\d*) (\d+\.?\d*) Td`)
	tmRegex      = regexp.MustCompile(`(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+Tm`)
	bfCharRegex  = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>`)
	bfRangeRegex = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>`)
)

type parserState struct {
	State             string
	CurrentTextObject TextObject
	TextObjects       []TextObject
}

func (p *PDFReader) parse(r io.Reader) ([]TextObject, error) {
	scanner := bufio.NewScanner(r)

	var result []TextObject
	state := parserState{
		State: "INITIAL",
	}

	textFound := 0
	for scanner.Scan() {
		line := scanner.Text()
		switch state.State {
		case "INITIAL":
			if strings.Contains(line, "BT") { // Begin Text block
				state.State = "TEXT_BLOCK"
				state.CurrentTextObject = TextObject{}
			}

		case "TEXT_BLOCK":
			// Parse font name
			if fontMatch := fontRegex.FindString(line); fontMatch != "" {
				parts := strings.Split(fontMatch, " ")
				if len(parts) >= 2 {
					if font, ok := p.fonts.Get(parts[0]); ok {
						state.CurrentTextObject.ResourceName = parts[0]
						state.CurrentTextObject.FontName = font.FontName
						if encoding, ok := font.FontDict.Find("Encoding"); ok {
							state.CurrentTextObject.Encoding = encoding.String()
						}
					}
					fontSize, _ := strconv.ParseFloat(parts[1], 64)
					state.CurrentTextObject.FontSize = fontSize
				}
			}

			// Parse text content
			if textMatch := textRegex.FindStringSubmatch(line); len(textMatch) > 1 {
				if isPDFDocEncoded(textMatch[1]) {
					state.CurrentTextObject.Text = pdfDocDecode(textMatch[1])
				} else if isUTF16(textMatch[1]) {
					state.CurrentTextObject.Text = utf16Decode(textMatch[1][2:])
				} else {
					switch state.CurrentTextObject.Encoding {
					case "WinAnsiEncoding":
						encoder := &byteEncoder{&winAnsiEncoding}
						state.CurrentTextObject.Text += encoder.Decode(textMatch[1])
					case "MacRomanEncoding":
						encoder := &byteEncoder{&macRomanEncoding}
						state.CurrentTextObject.Text += encoder.Decode(textMatch[1])
					case "Identity-H":
						var (
							toUnicode       types.Object
							toUnicodeStream *types.StreamDict
							encoder         *ToUnicodeDecoder
							err             error
							valid           bool
							ok              bool

							text = textMatch[1]
						)

						font, ok := p.fonts.Get(state.CurrentTextObject.ResourceName)
						if !ok {
							goto AssignText
						}

						toUnicode, ok = font.FontDict.Find("ToUnicode")
						if !ok {
							goto AssignText
						}

						toUnicodeStream, valid, err = p.ctx.DereferenceStreamDict(toUnicode)
						if err != nil {
							goto AssignText
						}

						if !valid {
							goto AssignText
						}

						err = toUnicodeStream.Decode()
						if err != nil {
							goto AssignText
						}

						encoder, err = NewToUnicodeDecoder(toUnicodeStream.Content)
						if err != nil {
							goto AssignText
						}
						text = encoder.Decode([]byte(text))

					AssignText:
						state.CurrentTextObject.Text += text
					default:
						state.CurrentTextObject.Text += textMatch[1]
					}
				}

				textFound++
			}

			// Parse position
			if posMatch := posRegex.FindStringSubmatch(line); len(posMatch) > 2 {
				if textFound == 0 {
					x, _ := strconv.ParseFloat(posMatch[1], 64)
					y, _ := strconv.ParseFloat(posMatch[2], 64)
					state.CurrentTextObject.Position.X += x
					state.CurrentTextObject.Position.Y += y
				}
			}

			// Parse position using Tm operator
			if tmMatch := tmRegex.FindStringSubmatch(line); len(tmMatch) > 6 {
				// tmMatch[1] through tmMatch[4] contain the transformation matrix values (a,b,c,d)
				// tmMatch[5] and tmMatch[6] contain the translation values (e,f) which give us x,y
				if textFound == 0 {
					x, _ := strconv.ParseFloat(tmMatch[5], 64)
					y, _ := strconv.ParseFloat(tmMatch[6], 64)
					state.CurrentTextObject.Position = Position{X: x, Y: y}
				}
			}

			// End of text block
			if strings.Contains(line, "ET") {
				if state.CurrentTextObject.Text != "" {
					result = append(result, state.CurrentTextObject)
				}
				state.State = "INITIAL"
				textFound = 0
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

type byteEncoder struct {
	table *[256]rune
}

func (e *byteEncoder) Decode(raw string) (text string) {
	r := make([]rune, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		r = append(r, e.table[raw[i]])
	}
	return string(r)
}

const noRune = unicode.ReplacementChar

func isPDFDocEncoded(s string) bool {
	if isUTF16(s) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if pdfDocEncoding[s[i]] == noRune {
			return false
		}
	}
	return true
}

func pdfDocDecode(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 || pdfDocEncoding[s[i]] != rune(s[i]) {
			goto Decode
		}
	}
	return s

Decode:
	r := make([]rune, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = pdfDocEncoding[s[i]]
	}
	return string(r)
}

func isUTF16(s string) bool {
	return len(s) >= 2 && s[0] == 0xfe && s[1] == 0xff && len(s)%2 == 0
}

func utf16Decode(s string) string {
	var u []uint16
	for i := 0; i < len(s); i += 2 {
		u = append(u, uint16(s[i])<<8|uint16(s[i+1]))
	}
	return string(utf16.Decode(u))
}

// See PDF 32000-1:2008, Table D.2
var pdfDocEncoding = [256]rune{
	noRune, noRune, noRune, noRune, noRune, noRune, noRune, noRune,
	noRune, 0x0009, 0x000a, noRune, noRune, 0x000d, noRune, noRune,
	noRune, noRune, noRune, noRune, noRune, noRune, noRune, noRune,
	0x02d8, 0x02c7, 0x02c6, 0x02d9, 0x02dd, 0x02db, 0x02da, 0x02dc,
	0x0020, 0x0021, 0x0022, 0x0023, 0x0024, 0x0025, 0x0026, 0x0027,
	0x0028, 0x0029, 0x002a, 0x002b, 0x002c, 0x002d, 0x002e, 0x002f,
	0x0030, 0x0031, 0x0032, 0x0033, 0x0034, 0x0035, 0x0036, 0x0037,
	0x0038, 0x0039, 0x003a, 0x003b, 0x003c, 0x003d, 0x003e, 0x003f,
	0x0040, 0x0041, 0x0042, 0x0043, 0x0044, 0x0045, 0x0046, 0x0047,
	0x0048, 0x0049, 0x004a, 0x004b, 0x004c, 0x004d, 0x004e, 0x004f,
	0x0050, 0x0051, 0x0052, 0x0053, 0x0054, 0x0055, 0x0056, 0x0057,
	0x0058, 0x0059, 0x005a, 0x005b, 0x005c, 0x005d, 0x005e, 0x005f,
	0x0060, 0x0061, 0x0062, 0x0063, 0x0064, 0x0065, 0x0066, 0x0067,
	0x0068, 0x0069, 0x006a, 0x006b, 0x006c, 0x006d, 0x006e, 0x006f,
	0x0070, 0x0071, 0x0072, 0x0073, 0x0074, 0x0075, 0x0076, 0x0077,
	0x0078, 0x0079, 0x007a, 0x007b, 0x007c, 0x007d, 0x007e, noRune,
	0x2022, 0x2020, 0x2021, 0x2026, 0x2014, 0x2013, 0x0192, 0x2044,
	0x2039, 0x203a, 0x2212, 0x2030, 0x201e, 0x201c, 0x201d, 0x2018,
	0x2019, 0x201a, 0x2122, 0xfb01, 0xfb02, 0x0141, 0x0152, 0x0160,
	0x0178, 0x017d, 0x0131, 0x0142, 0x0153, 0x0161, 0x017e, noRune,
	0x20ac, 0x00a1, 0x00a2, 0x00a3, 0x00a4, 0x00a5, 0x00a6, 0x00a7,
	0x00a8, 0x00a9, 0x00aa, 0x00ab, 0x00ac, noRune, 0x00ae, 0x00af,
	0x00b0, 0x00b1, 0x00b2, 0x00b3, 0x00b4, 0x00b5, 0x00b6, 0x00b7,
	0x00b8, 0x00b9, 0x00ba, 0x00bb, 0x00bc, 0x00bd, 0x00be, 0x00bf,
	0x00c0, 0x00c1, 0x00c2, 0x00c3, 0x00c4, 0x00c5, 0x00c6, 0x00c7,
	0x00c8, 0x00c9, 0x00ca, 0x00cb, 0x00cc, 0x00cd, 0x00ce, 0x00cf,
	0x00d0, 0x00d1, 0x00d2, 0x00d3, 0x00d4, 0x00d5, 0x00d6, 0x00d7,
	0x00d8, 0x00d9, 0x00da, 0x00db, 0x00dc, 0x00dd, 0x00de, 0x00df,
	0x00e0, 0x00e1, 0x00e2, 0x00e3, 0x00e4, 0x00e5, 0x00e6, 0x00e7,
	0x00e8, 0x00e9, 0x00ea, 0x00eb, 0x00ec, 0x00ed, 0x00ee, 0x00ef,
	0x00f0, 0x00f1, 0x00f2, 0x00f3, 0x00f4, 0x00f5, 0x00f6, 0x00f7,
	0x00f8, 0x00f9, 0x00fa, 0x00fb, 0x00fc, 0x00fd, 0x00fe, 0x00ff,
}

var winAnsiEncoding = [256]rune{
	0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007,
	0x0008, 0x0009, 0x000a, 0x000b, 0x000c, 0x000d, 0x000e, 0x000f,
	0x0010, 0x0011, 0x0012, 0x0013, 0x0014, 0x0015, 0x0016, 0x0017,
	0x0018, 0x0019, 0x001a, 0x001b, 0x001c, 0x001d, 0x001e, 0x001f,
	0x0020, 0x0021, 0x0022, 0x0023, 0x0024, 0x0025, 0x0026, 0x0027,
	0x0028, 0x0029, 0x002a, 0x002b, 0x002c, 0x002d, 0x002e, 0x002f,
	0x0030, 0x0031, 0x0032, 0x0033, 0x0034, 0x0035, 0x0036, 0x0037,
	0x0038, 0x0039, 0x003a, 0x003b, 0x003c, 0x003d, 0x003e, 0x003f,
	0x0040, 0x0041, 0x0042, 0x0043, 0x0044, 0x0045, 0x0046, 0x0047,
	0x0048, 0x0049, 0x004a, 0x004b, 0x004c, 0x004d, 0x004e, 0x004f,
	0x0050, 0x0051, 0x0052, 0x0053, 0x0054, 0x0055, 0x0056, 0x0057,
	0x0058, 0x0059, 0x005a, 0x005b, 0x005c, 0x005d, 0x005e, 0x005f,
	0x0060, 0x0061, 0x0062, 0x0063, 0x0064, 0x0065, 0x0066, 0x0067,
	0x0068, 0x0069, 0x006a, 0x006b, 0x006c, 0x006d, 0x006e, 0x006f,
	0x0070, 0x0071, 0x0072, 0x0073, 0x0074, 0x0075, 0x0076, 0x0077,
	0x0078, 0x0079, 0x007a, 0x007b, 0x007c, 0x007d, 0x007e, 0x007f,
	0x20ac, noRune, 0x201a, 0x0192, 0x201e, 0x2026, 0x2020, 0x2021,
	0x02c6, 0x2030, 0x0160, 0x2039, 0x0152, noRune, 0x017d, noRune,
	noRune, 0x2018, 0x2019, 0x201c, 0x201d, 0x2022, 0x2013, 0x2014,
	0x02dc, 0x2122, 0x0161, 0x203a, 0x0153, noRune, 0x017e, 0x0178,
	0x00a0, 0x00a1, 0x00a2, 0x00a3, 0x00a4, 0x00a5, 0x00a6, 0x00a7,
	0x00a8, 0x00a9, 0x00aa, 0x00ab, 0x00ac, 0x00ad, 0x00ae, 0x00af,
	0x00b0, 0x00b1, 0x00b2, 0x00b3, 0x00b4, 0x00b5, 0x00b6, 0x00b7,
	0x00b8, 0x00b9, 0x00ba, 0x00bb, 0x00bc, 0x00bd, 0x00be, 0x00bf,
	0x00c0, 0x00c1, 0x00c2, 0x00c3, 0x00c4, 0x00c5, 0x00c6, 0x00c7,
	0x00c8, 0x00c9, 0x00ca, 0x00cb, 0x00cc, 0x00cd, 0x00ce, 0x00cf,
	0x00d0, 0x00d1, 0x00d2, 0x00d3, 0x00d4, 0x00d5, 0x00d6, 0x00d7,
	0x00d8, 0x00d9, 0x00da, 0x00db, 0x00dc, 0x00dd, 0x00de, 0x00df,
	0x00e0, 0x00e1, 0x00e2, 0x00e3, 0x00e4, 0x00e5, 0x00e6, 0x00e7,
	0x00e8, 0x00e9, 0x00ea, 0x00eb, 0x00ec, 0x00ed, 0x00ee, 0x00ef,
	0x00f0, 0x00f1, 0x00f2, 0x00f3, 0x00f4, 0x00f5, 0x00f6, 0x00f7,
	0x00f8, 0x00f9, 0x00fa, 0x00fb, 0x00fc, 0x00fd, 0x00fe, 0x00ff,
}

var macRomanEncoding = [256]rune{
	0x0000, 0x0001, 0x0002, 0x0003, 0x0004, 0x0005, 0x0006, 0x0007,
	0x0008, 0x0009, 0x000a, 0x000b, 0x000c, 0x000d, 0x000e, 0x000f,
	0x0010, 0x0011, 0x0012, 0x0013, 0x0014, 0x0015, 0x0016, 0x0017,
	0x0018, 0x0019, 0x001a, 0x001b, 0x001c, 0x001d, 0x001e, 0x001f,
	0x0020, 0x0021, 0x0022, 0x0023, 0x0024, 0x0025, 0x0026, 0x0027,
	0x0028, 0x0029, 0x002a, 0x002b, 0x002c, 0x002d, 0x002e, 0x002f,
	0x0030, 0x0031, 0x0032, 0x0033, 0x0034, 0x0035, 0x0036, 0x0037,
	0x0038, 0x0039, 0x003a, 0x003b, 0x003c, 0x003d, 0x003e, 0x003f,
	0x0040, 0x0041, 0x0042, 0x0043, 0x0044, 0x0045, 0x0046, 0x0047,
	0x0048, 0x0049, 0x004a, 0x004b, 0x004c, 0x004d, 0x004e, 0x004f,
	0x0050, 0x0051, 0x0052, 0x0053, 0x0054, 0x0055, 0x0056, 0x0057,
	0x0058, 0x0059, 0x005a, 0x005b, 0x005c, 0x005d, 0x005e, 0x005f,
	0x0060, 0x0061, 0x0062, 0x0063, 0x0064, 0x0065, 0x0066, 0x0067,
	0x0068, 0x0069, 0x006a, 0x006b, 0x006c, 0x006d, 0x006e, 0x006f,
	0x0070, 0x0071, 0x0072, 0x0073, 0x0074, 0x0075, 0x0076, 0x0077,
	0x0078, 0x0079, 0x007a, 0x007b, 0x007c, 0x007d, 0x007e, 0x007f,
	0x00c4, 0x00c5, 0x00c7, 0x00c9, 0x00d1, 0x00d6, 0x00dc, 0x00e1,
	0x00e0, 0x00e2, 0x00e4, 0x00e3, 0x00e5, 0x00e7, 0x00e9, 0x00e8,
	0x00ea, 0x00eb, 0x00ed, 0x00ec, 0x00ee, 0x00ef, 0x00f1, 0x00f3,
	0x00f2, 0x00f4, 0x00f6, 0x00f5, 0x00fa, 0x00f9, 0x00fb, 0x00fc,
	0x2020, 0x00b0, 0x00a2, 0x00a3, 0x00a7, 0x2022, 0x00b6, 0x00df,
	0x00ae, 0x00a9, 0x2122, 0x00b4, 0x00a8, 0x2260, 0x00c6, 0x00d8,
	0x221e, 0x00b1, 0x2264, 0x2265, 0x00a5, 0x00b5, 0x2202, 0x2211,
	0x220f, 0x03c0, 0x222b, 0x00aa, 0x00ba, 0x03a9, 0x00e6, 0x00f8,
	0x00bf, 0x00a1, 0x00ac, 0x221a, 0x0192, 0x2248, 0x2206, 0x00ab,
	0x00bb, 0x2026, 0x00a0, 0x00c0, 0x00c3, 0x00d5, 0x0152, 0x0153,
	0x2013, 0x2014, 0x201c, 0x201d, 0x2018, 0x2019, 0x00f7, 0x25ca,
	0x00ff, 0x0178, 0x2044, 0x20ac, 0x2039, 0x203a, 0xfb01, 0xfb02,
	0x2021, 0x00b7, 0x201a, 0x201e, 0x2030, 0x00c2, 0x00ca, 0x00c1,
	0x00cb, 0x00c8, 0x00cd, 0x00ce, 0x00cf, 0x00cc, 0x00d3, 0x00d4,
	0xf8ff, 0x00d2, 0x00da, 0x00db, 0x00d9, 0x0131, 0x02c6, 0x02dc,
	0x00af, 0x02d8, 0x02d9, 0x02da, 0x00b8, 0x02dd, 0x02db, 0x02c7,
}

type ToUnicodeDecoder struct {
	charMappings map[uint32]rune
}

// NewToUnicodeDecoder creates a new decoder from a ToUnicode CMap
func NewToUnicodeDecoder(toUnicodeCMap []byte) (*ToUnicodeDecoder, error) {
	decoder := &ToUnicodeDecoder{
		charMappings: make(map[uint32]rune),
	}

	// Parse the CMap
	err := decoder.parseCMap(toUnicodeCMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ToUnicode CMap: %v", err)
	}

	return decoder, nil
}

// parseCMap handles parsing the ToUnicode CMap
func (d *ToUnicodeDecoder) parseCMap(cmapData []byte) error {
	reader := bytes.NewReader(cmapData)
	scanner := bufio.NewScanner(reader)

	// Parsing state variables
	var (
		inBfChar  bool
		inBfRange bool
	)

	for scanner.Scan() {
		lineStr := scanner.Text()

		// Detect CMap sections
		switch {
		case strings.Contains(lineStr, "beginbfchar"):
			inBfChar = true
			continue
		case strings.Contains(lineStr, "endbfchar"):
			inBfChar = false
			continue
		case strings.Contains(lineStr, "beginbfrange"):
			inBfRange = true
			continue
		case strings.Contains(lineStr, "endbfrange"):
			inBfRange = false
			continue
		}

		// Process different sections
		switch {
		case inBfChar:
			// Direct character mappings
			if match := bfCharRegex.FindStringSubmatch(lineStr); match != nil {
				code, _ := strconv.ParseUint(match[1], 16, 32)
				unicode, _ := strconv.ParseUint(match[2], 16, 32)
				d.charMappings[uint32(code)] = rune(unicode)
			}
		case inBfRange:
			// Range-based character mappings
			if match := bfRangeRegex.FindStringSubmatch(lineStr); match != nil {
				startCode, _ := strconv.ParseUint(match[1], 16, 32)
				endCode, _ := strconv.ParseUint(match[2], 16, 32)
				targetCode, _ := strconv.ParseUint(match[3], 16, 32)

				for i := startCode; i <= endCode; i++ {
					d.charMappings[uint32(i)] = rune(targetCode + (i - startCode))
				}
			}
		}
	}

	return nil
}

// Decode converts Identity-H encoded bytes to a Unicode string
func (d *ToUnicodeDecoder) Decode(input []byte) string {
	var result []rune

	// Process input in 2-byte chunks for Identity-H
	for i := 0; i < len(input); i += 2 {
		if i+1 >= len(input) {
			break
		}

		// Convert 2 bytes to uint32
		code := binary.BigEndian.Uint16(input[i : i+2])

		// Look up Unicode mapping
		if unicode, ok := d.charMappings[uint32(code)]; ok {
			result = append(result, unicode)
		} else {
			// Fallback: use the code point itself if no mapping found
			result = append(result, rune(code))
		}
	}

	return string(result)
}
