// this file is lexer for pdf content syntax

package pdfextract

type PDFTextToken int

const (
	TokenEOF PDFTextToken = iota
	TokenFont
	TokenText
	TokenPosition
	TokenMatrix
)

type PDFCMapToken int

const (
	TokenCMapEOF PDFCMapToken = iota
	TokenCMapError
	TokenCMapName
	TokenCMapBeginDict
	TokenCMapEndDict
	TokenCMapBeginCodeSpaceRange
	TokenCMapEndCodeSpaceRange
	TokenCMapBeginBFRange
	TokenCMapEndBFRange
	TokenCMapBeginBFChar
	TokenCMapEndBFChar
	TokenCMapHexValue
)
