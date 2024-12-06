package pdfextract

import (
	"sort"

	"github.com/mrrizkin/omniscan/pkg/pdf-extract/utils"
)

type Page struct {
	content []TextObject
}

func (p *Page) GetTextByRow(tolerance float64) (Rows, error) {
	row := make(Rows, 0)
	currentPosition := 0.0
	rowIndex := -1
	for _, object := range p.content {
		if !utils.IsEqualTolerance(object.Position.Y, currentPosition, tolerance) {
			row = append(row, Row{
				Content: []Text{{
					Font:     object.FontName,
					FontSize: object.FontSize,
					X:        object.Position.X,
					Y:        object.Position.Y,
					S:        object.Text,
				}},
				Position: object.Position.Y,
			})
			currentPosition = object.Position.Y
			rowIndex++
		} else {
			if rowIndex == -1 {
				continue
			}
			row[rowIndex].Content = append(row[rowIndex].Content, Text{
				Font:     object.FontName,
				FontSize: object.FontSize,
				X:        object.Position.X,
				Y:        object.Position.Y,
				S:        object.Text,
			})
		}
	}

	sort.Sort(row)
	return row, nil
}
