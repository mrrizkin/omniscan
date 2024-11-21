package pdfextract

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	fontRegex = regexp.MustCompile(`/F\d+ \d+ Tf`)
	textRegex = regexp.MustCompile(`\(([^)]+)\)\s*Tj`)
	posRegex  = regexp.MustCompile(`(\d+\.?\d*) (\d+\.?\d*) Td`)
	tmRegex   = regexp.MustCompile(`(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+(\d+\.?\d*)\s+Tm`)
)

type parserState struct {
	State             string
	CurrentTextObject TextObject
	TextObjects       []TextObject
}

func parser(f *os.File) ([]TextObject, error) {
	scanner := bufio.NewScanner(f)

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
					state.CurrentTextObject.FontName = parts[0]
					fontSize, _ := strconv.ParseFloat(parts[1], 64)
					state.CurrentTextObject.FontSize = fontSize
				}
			}

			// Parse text content
			if textMatch := textRegex.FindStringSubmatch(line); len(textMatch) > 1 {
				state.CurrentTextObject.Text += textMatch[1]
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
