package pdfcpu

import (
	"bufio"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mrrizkin/omniscan/pkg/pdf/provider/pdfcpu/encoder"
)

var (
	fontRegex      = regexp.MustCompile(`\/(F\d+) (\d+\.?\d*) Tf`)
	textRegex      = regexp.MustCompile(`\(([^)]+)\)\s*Tj`)
	textArrayRegex = regexp.MustCompile(`\[*((?:\(([^)]+)\)+)+)\]*\s*TJ`)
	posRegex       = regexp.MustCompile(`(\d+\.?\d*) (\d+\.?\d*) Td`)
	tmRegex        = regexp.MustCompile(`(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+Tm`)
	bfCharRegex    = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>`)
	bfRangeRegex   = regexp.MustCompile(`<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>\s*<([0-9A-Fa-f]+)>`)
)

type parserState struct {
	State             string
	CurrentFont       *fontObject
	CurrentFontSize   float64
	CurrentTextObject TextObject
	TextObjects       Content
}

func (p *PDFCPU) parse(r io.Reader) (Content, error) {
	scanner := bufio.NewScanner(r)

	var result Content
	state := parserState{
		State: "INITIAL",
	}

	textFound := 0
	for scanner.Scan() {
		line := scanner.Text()
		switch state.State {
		case "INITIAL":
			if strings.HasPrefix(line, "BT") { // Begin Text block
				state.State = "TEXT_BLOCK"
				state.CurrentTextObject = TextObject{}

				if state.CurrentFont != nil {
					state.CurrentTextObject.FontName = state.CurrentFont.FontName
					state.CurrentTextObject.FontSize = state.CurrentFontSize
				}
			}

		case "TEXT_BLOCK":
			// Parse font name
			if fontMatch := fontRegex.FindStringSubmatch(line); len(fontMatch) > 2 {
				if font, ok := p.fonts.Get(fontMatch[1]); ok {
					state.CurrentTextObject.FontName = font.FontName
					state.CurrentFont = font
				}
				fontSize, _ := strconv.ParseFloat(fontMatch[2], 64)
				state.CurrentTextObject.FontSize = fontSize
				state.CurrentFontSize = fontSize
			}

			// Parse text content Tj
			if textMatch := textRegex.FindStringSubmatch(line); len(textMatch) > 1 {
				state.CurrentTextObject.Text += decodeText(line, textMatch[1], state.CurrentFont)
				textFound++
			}

			// Parse text content TJ
			if textArrayMatch := textArrayRegex.FindStringSubmatch(line); len(textArrayMatch) > 2 {
				state.CurrentTextObject.Text += decodeText(line, textArrayMatch[2], state.CurrentFont)
				textFound++
			}

			// Parse position
			if posMatch := posRegex.FindStringSubmatch(line); len(posMatch) > 2 {
				x, _ := strconv.ParseFloat(posMatch[1], 64)
				y, _ := strconv.ParseFloat(posMatch[2], 64)
				if state.CurrentTextObject.Position.Y != 0 &&
					math.Abs(state.CurrentTextObject.Position.Y-y) > state.CurrentFontSize {
					state.CurrentTextObject.Text += "\n"
				}
				state.CurrentTextObject.Position.X += x
				state.CurrentTextObject.Position.Y += y
			}

			// Parse position using Tm operator
			if tmMatch := tmRegex.FindStringSubmatch(line); len(tmMatch) > 6 {
				x, _ := strconv.ParseFloat(tmMatch[5], 64)
				y, _ := strconv.ParseFloat(tmMatch[6], 64)
				if state.CurrentTextObject.Position.Y != 0 &&
					math.Abs(state.CurrentTextObject.Position.Y-y) > state.CurrentFontSize {
					state.CurrentTextObject.Text += "\n"
				}
				state.CurrentTextObject.Position = Position{X: x, Y: y}
			}

			// End of text block
			if strings.HasPrefix(line, "ET") {
				state.State = "INITIAL"
				textFound = 0
				if state.CurrentTextObject.Text != "" {
					result = append(result, state.CurrentTextObject)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Sort(result)
	return result, nil
}

func sanitizeIdentityH(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "(") {
		text = text[1:]
	}
	if strings.HasSuffix(text, "Tj") {
		text = text[:len(text)-2]
	}
	text = strings.TrimSpace(text)
	if strings.HasSuffix(text, ")") {
		text = text[:len(text)-1]
	}

	tmp := []byte("")
	reader := strings.NewReader(text)
	depth := 1
Loop:
	for {
		if reader.Len() == 0 {
			break
		}

		c, err := reader.ReadByte()
		if err != nil {
			break
		}

		switch c {
		default:
			tmp = append(tmp, c)
		case '(':
			depth++
			tmp = append(tmp, c)
		case ')':
			if depth--; depth == 0 {
				break Loop
			}
			tmp = append(tmp, c)
		case '\\':
			c, err := reader.ReadByte()
			if err != nil {
				break
			}
			switch c {
			default:
				reader.UnreadByte()
				tmp = append(tmp, '\\', c)
			case 'n':
				tmp = append(tmp, '\n')
			case 'r':
				tmp = append(tmp, '\r')
			case 'b':
				tmp = append(tmp, '\b')
			case 't':
				tmp = append(tmp, '\t')
			case 'f':
				tmp = append(tmp, '\f')
			case '(', ')', '\\':
				tmp = append(tmp, c)
			case '\r':
				if innerC, err := reader.ReadByte(); err == nil {
					if innerC != '\n' {
						reader.UnreadByte()
					}
				}
				fallthrough
			case '\n':
				// no append
			case '0', '1', '2', '3', '4', '5', '6', '7':
				x := int(c - '0')
				for i := 0; i < 2; i++ {
					if c, err := reader.ReadByte(); err == nil {
						if c < '0' || c > '7' {
							reader.UnreadByte()
							break
						}
						x = x*8 + int(c-'0')
					}
				}
				if x > 255 {
					reader.UnreadByte()
					break
				}
				tmp = append(tmp, byte(x))
			}
		}
	}

	return string(tmp)
}

func decodeText(raw, text string, font *fontObject) string {
	if encoder.IsPDFDocEncoded(text) {
		return encoder.PdfDocDecode(text)
	} else if encoder.IsUTF16(text) {
		return encoder.Utf16Decode(text[2:])
	} else {
		if font != nil {
			return font.Decode(sanitizeIdentityH(raw))
		} else {
			return text
		}
	}
}
