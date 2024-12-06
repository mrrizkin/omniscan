package encoder

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/mrrizkin/omniscan/pkg/pdf-extract/utils"
)

var (
	registryRe   = regexp.MustCompile(`/Registry\s*\((.*?)\)`)
	orderingRe   = regexp.MustCompile(`/Ordering\s*\((.*?)\)`)
	supplementRe = regexp.MustCompile(`/Supplement\s*(\d+)`)
	nameRe       = regexp.MustCompile(`/CMapName\s*/(.*?)\s`)
	typeRe       = regexp.MustCompile(`/CMapType\s*(\d+)`)
	hexRe        = regexp.MustCompile(`^<([0-9A-Fa-f]+)>$`)
)

type (
	byteRange struct {
		low  string
		high string
	}

	bfChar struct {
		orig string
		repl string
	}

	bfRange struct {
		lo  string
		hi  string
		dst interface{}
	}

	CMap struct {
		space   [4][]byteRange
		bfchar  []bfChar
		bfrange []bfRange
	}
)

func (cm *CMap) Merge(other *CMap) {
	cm.bfchar = append(cm.bfchar, other.bfchar...)
	cm.bfrange = append(cm.bfrange, other.bfrange...)
	for i := 0; i < 4; i++ {
		cm.space[i] = append(cm.space[i], other.space[i]...)
	}
}

// switch c = b.readByte(); c {
// 		default:
// 			b.errorf("invalid escape sequence \\%c", c)
// 			tmp = append(tmp, '\\', c)
// 		case 'n':
// 			tmp = append(tmp, '\n')
// 		case 'r':
// 			tmp = append(tmp, '\r')
// 		case 'b':
// 			tmp = append(tmp, '\b')
// 		case 't':
// 			tmp = append(tmp, '\t')
// 		case 'f':
// 			tmp = append(tmp, '\f')
// 		case '(', ')', '\\':
// 			tmp = append(tmp, c)
// 		case '\r':
// 			if b.readByte() != '\n' {
// 				b.unreadByte()
// 			}
// 			fallthrough
// 		case '\n':
// 			// no append
// 		case '0', '1', '2', '3', '4', '5', '6', '7':
// 			x := int(c - '0')
// 			for i := 0; i < 2; i++ {
// 				c = b.readByte()
// 				if c < '0' || c > '7' {
// 					b.unreadByte()
// 					break
// 				}
// 				x = x*8 + int(c-'0')
// 			}
// 			if x > 255 {
// 				b.errorf("invalid octal escape \\%03o", x)
// 			}
// 			tmp = append(tmp, byte(x))
// 		}

func (cm *CMap) Decode(original string) string {
	if strings.Contains(original, "\\") {
		if strings.Contains(original, "\\\\") {
			original = strings.ReplaceAll(original, "\\\\", "\\")
		} else if strings.Contains(original, "\\\r\n") {
			original = strings.ReplaceAll(original, "\\\r\n", "\r\n")
		} else if strings.Contains(original, "\\\n") {
			original = strings.ReplaceAll(original, "\\\n", "\n")
		} else if strings.Contains(original, "\\n") {
			original = strings.ReplaceAll(original, "\\n", "\n")
		} else if strings.Contains(original, "\\r") {
			original = strings.ReplaceAll(original, "\\r", "\r")
		} else if strings.Contains(original, "\\b") {
			original = strings.ReplaceAll(original, "\\b", "\b")
		} else if strings.Contains(original, "\\t") {
			original = strings.ReplaceAll(original, "\\t", "\t")
		} else if strings.Contains(original, "\\f") {
			original = strings.ReplaceAll(original, "\\f", "\f")
		} else if strings.Contains(original, "\\(") {
			original = strings.ReplaceAll(original, "\\(", "(")
		} else if strings.Contains(original, "\\)") {
			original = strings.ReplaceAll(original, "\\)", ")")
		} else {
			original = strings.ReplaceAll(original, "\\", "")
		}
	}

	raw := original
	var r []rune

Decode:
	for len(raw) > 0 {
		for n := 1; n <= 4 && n <= len(raw); n++ { // number of digits in character replacement (1-4 possible)
			for _, space := range cm.space[n-1] { // find matching codespace Ranges for number of digits
				if space.low <= raw[:n] && raw[:n] <= space.high { // see if value is in range
					text := raw[:n]
					raw = raw[n:]
					for _, bfchar := range cm.bfchar { // check for matching bfchar
						if len(bfchar.orig) == n && bfchar.orig == text {
							r = append(r, []rune(Utf16Decode(bfchar.repl))...)
							continue Decode
						}
					}
					for _, bfrange := range cm.bfrange { // check for matching bfrange
						if len(bfrange.lo) == n && bfrange.lo <= text && text <= bfrange.hi {
							if dst, ok := bfrange.dst.(string); ok {
								s := dst
								if bfrange.lo != text { // value isn't at the beginning of the range so scale result
									b := []byte(s)
									b[len(b)-1] += text[len(text)-1] - bfrange.lo[len(bfrange.lo)-1] // increment last byte by difference
									s = string(b)
								}
								r = append(r, []rune(Utf16Decode(s))...)
								continue Decode
							}
							r = append(r, NoRune)
							continue Decode
						}
					}
					r = append(r, NoRune)
					continue Decode
				}
			}
		}

		r = append(r, NoRune)
		raw = raw[1:]
	}

	// check if there's any string that non utf8
	if !utils.IsUTF8(string(r)) {
		return string(r)
	}

	return string(r)
}

type StructedCMap struct {
	MapCount int
	Mapping  []string
}

func (s *StructedCMap) Push(hex string) {
	s.Mapping = append(s.Mapping, hex)
}

func ParseCmap(cmapBytes []byte) (*CMap, error) {
	cmap := &CMap{}
	tokenizer := NewTokenizer(bytes.NewReader(cmapBytes))
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, errors.New("no tokens")
	}

	var currentIntParam int
	var currentSection string

	mapping := make(map[string]*StructedCMap)

	for _, token := range tokens {
		switch token.Type {
		case "INT":
			i, err := strconv.Atoi(token.Value)
			if err != nil {
				currentIntParam = 0
				continue
			}

			currentIntParam = i
		case "KEYWORD":
			currentSection = token.Value
			switch currentSection {
			case "endcodespacerange":
				codespacerange, ok := mapping["codespacerange"]
				if !ok {
					continue
				}

				size := len(codespacerange.Mapping) / codespacerange.MapCount
				chunk := utils.Chunk(codespacerange.Mapping, size)

				for _, m := range chunk {
					if size != 2 {
						continue
					}

					low, _ := utils.Hex2Bytes(m[0])
					high, _ := utils.Hex2Bytes(m[1])

					cmap.space[len(low)-1] = append(cmap.space[len(low)-1], byteRange{string(low), string(high)})
				}
			case "endbfrange":
				bfrange, ok := mapping["bfrange"]
				if !ok {
					continue
				}

				size := len(bfrange.Mapping) / bfrange.MapCount
				chunk := utils.Chunk(bfrange.Mapping, size)

				for _, m := range chunk {
					if size != 2 && size != 3 {
						continue
					}

					// lo, hi, dst := m[0], m[0], m[1]
					lo, _ := utils.Hex2Bytes(m[0])
					hi := lo
					dst, _ := utils.Hex2Bytes(m[1])
					if size == 3 {
						hi = dst
						dst, _ = utils.Hex2Bytes(m[2])
					}

					cmap.bfrange = append(cmap.bfrange, bfRange{string(lo), string(hi), string(dst)})
				}

			case "endbfchar":
				bfchar, ok := mapping["bfchar"]
				if !ok {
					continue
				}

				size := len(bfchar.Mapping) / bfchar.MapCount
				chunk := utils.Chunk(bfchar.Mapping, size)

				for _, m := range chunk {
					if size != 2 {
						continue
					}

					orig, _ := utils.Hex2Bytes(m[0])
					repl, _ := utils.Hex2Bytes(m[1])
					cmap.bfchar = append(cmap.bfchar, bfChar{string(orig), string(repl)})
				}
			}
		case "HEX":
			switch currentSection {
			case "begincodespacerange":
				if _, ok := mapping["codespacerange"]; !ok {
					mapping["codespacerange"] = &StructedCMap{
						MapCount: currentIntParam,
						Mapping:  make([]string, 0),
					}
				}

				mapping["codespacerange"].Push(token.Value)
			case "beginbfrange":
				if _, ok := mapping["bfrange"]; !ok {
					mapping["bfrange"] = &StructedCMap{
						MapCount: currentIntParam,
						Mapping:  make([]string, 0),
					}
				}

				mapping["bfrange"].Push(token.Value)
			case "beginbfchar":
				if _, ok := mapping["bfchar"]; !ok {
					mapping["bfchar"] = &StructedCMap{
						MapCount: currentIntParam,
						Mapping:  make([]string, 0),
					}
				}

				mapping["bfchar"].Push(token.Value)
			}
		}
	}

	return cmap, nil
}

type Tokenizer struct {
	scanner *bufio.Scanner
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{scanner: bufio.NewScanner(r)}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0)
	for t.scanner.Scan() {
		line := t.scanner.Text()
		tokens = append(tokens, tokenizeLine(line)...)
	}
	if err := t.scanner.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

type Token struct {
	Type  string
	Value string
}

var (
	keyRegex = regexp.MustCompile(`(begincmap|endcmap|begincodespacerange|endcodespacerange|beginbfrange|endbfrange|beginbfchar|endbfchar|defineresource|def|CMapName|CMapType)`)
	cmdRegex = regexp.MustCompile(`\/([a-zA-Z]+)`)
	hexRegex = regexp.MustCompile(`<([0-9A-Fa-f]+)>`)
	intRegex = regexp.MustCompile(`([\d]+) `)
)

func tokenizeLine(line string) []Token {
	var tokens []Token
	matches := hexRegex.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		tokens = append(tokens, Token{Type: "HEX", Value: match[1]})
	}
	matches = cmdRegex.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		tokens = append(tokens, Token{Type: "CMD", Value: match[1]})
	}
	matches = intRegex.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		tokens = append(tokens, Token{Type: "INT", Value: match[1]})
	}

	matches = keyRegex.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		tokens = append(tokens, Token{Type: "KEYWORD", Value: match[1]})
	}

	return tokens
}
