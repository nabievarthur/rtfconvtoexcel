package service

import (
	"log"
	"os"
	"regexp"

	"github.com/j45k4/rtf"
	"github.com/xuri/excelize/v2"
)

type RtfConverter struct{}

func NewRtfConverter() *RtfConverter {
	return &RtfConverter{}
}

// Регулярка для удаления шапки
var headerRe = regexp.MustCompile(`(?m)^<1><3>--   --     -   <4>   <12>.*?<13>\. ч`)

func (r *RtfConverter) ConvertRtfToExcel(fileName string, excelFileName string) bool {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return false
	}

	text := rtf.StripRichTextFormat(string(b))
	line := headerRe.ReplaceAllString(text, "")

	rows := r.splitDataRows(line)

	f := excelize.NewFile()
	sheet := "Sheet1"

	// Заголовки
	header := []string{"<1>", "<3>", "--", "--", "-", "<4>", "<12>", "<Фабула>", "<13>", ".", "ч"}
	for colIdx, value := range header {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		f.SetCellValue(sheet, cell, value)
	}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheet, "A1", "K1", headerStyle)

	for rowIdx, row := range rows {
		for colIdx, value := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	// Автоширина
	cols := []string{"D", "E", "F", "G", "H"}
	widths := make(map[int]float64)
	for i := 3; i <= 7; i++ {
		widths[i] = float64(len(header[i]))
	}
	for _, row := range rows {
		for i := 3; i <= 7; i++ {
			if len(row) > i {
				thisLen := len(row[i])
				calculated := float64(thisLen)*1.2 + 8
				if calculated > widths[i] {
					widths[i] = calculated
				}
			}
		}
	}
	for i, col := range cols {
		w := widths[i+3]
		// минимальная ширина
		if w < 10 {
			w = 10
		}
		// максимальная ширина
		if w > 100 {
			w = 100
		}
		f.SetColWidth(sheet, col, col, w)
	}

	for i := 1; i <= len(header); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		if col == "D" || col == "E" || col == "F" || col == "G" || col == "H" {
			continue
		}
		f.SetColWidth(sheet, col, col, 5)
	}

	// Альбомная ориентация
	orientation := "landscape"
	pageLayout := &excelize.PageLayoutOptions{
		Orientation: &orientation,
	}
	if err := f.SetPageLayout(sheet, pageLayout); err != nil {
		log.Println("Error setting page layout:", err)
	}

	// Поля 1см
	l, rt, t, bt, h, ft := 0.3937, 0.3937, 0.3937, 0.3937, 0.0, 0.0
	opts := &excelize.PageLayoutMarginsOptions{
		Left:   &l,
		Right:  &rt,
		Top:    &t,
		Bottom: &bt,
		Header: &h,
		Footer: &ft,
	}
	if err := f.SetPageMargins(sheet, opts); err != nil {
		log.Println("Error setting margins:", err)
	}

	if err := f.SaveAs(excelFileName); err != nil {
		return false
	}
	return true
}

// Регулярка
func (r *RtfConverter) splitDataRows(data string) [][]string {
	re := regexp.MustCompile(`(\d{2})\s*(\d)\s*(\d{2})\s*(\d{8})\s*(\d{6})\s*(\d{3})\s*(\d{10})\s*([^\d]+)\s*(\d{3})\s*(\d{2})\s*(\d)`)
	matches := re.FindAllStringSubmatch(data, -1)
	result := [][]string{}
	for _, match := range matches {
		if len(match) == 12 {
			result = append(result, match[1:])
		}
	}
	return result
}
