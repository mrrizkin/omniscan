package encoder

import (
	"unicode/utf16"
)

func IsPDFDocEncoded(s string) bool {
	if IsUTF16(s) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if PdfDocEncoding[s[i]] == NoRune {
			return false
		}
	}
	return true
}

func PdfDocDecode(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 || PdfDocEncoding[s[i]] != rune(s[i]) {
			goto Decode
		}
	}
	return s

Decode:
	r := make([]rune, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = PdfDocEncoding[s[i]]
	}
	return string(r)
}

func IsUTF16(s string) bool {
	return len(s) >= 2 && s[0] == 0xfe && s[1] == 0xff && len(s)%2 == 0
}

func Utf16Decode(s string) string {
	var u []uint16
	for i := 0; i < len(s); i += 2 {
		u = append(u, uint16(s[i])<<8|uint16(s[i+1]))
	}
	return string(utf16.Decode(u))
}
