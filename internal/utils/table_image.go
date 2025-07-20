package utils

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
	"github.com/vsespontanno/tgbot_fschedule/internal/types"
)

func TableImage(data []types.Standing) (*bytes.Buffer, error) {
	// Константы
	const (
		width        = 780
		height       = 920
		padding      = 10
		rowHeight    = 40
		headerHeight = 60
		fontSize     = 20
		numCols      = 10
	)

	// Динамически определяем ширину столбцов
	colWidths := []int{
		50,  // # - Номер позиции
		280, // Команда - Уменьшено с 300 до 280, чтобы освободить место для других колонок
		55,  // И - Увеличено с 50 до 55
		55,  // В - Увеличено с 50 до 55
		55,  // Н - Увеличено с 50 до 55
		55,  // П - Увеличено с 50 до 55
		55,  // ГЗ - Увеличено с 50 до 55
		55,  // ГП - Увеличено с 50 до 55
		55,  // РГ - Увеличено с 50 до 55
		55,  // О - Увеличено с 50 до 55
	}

	dc := gg.NewContext(width, height)

	// Задаем фон
	dc.SetColor(color.RGBA{18, 18, 18, 255}) // Dark background
	dc.Clear()

	// Задаем цвет текста
	dc.SetColor(color.RGBA{230, 230, 230, 255}) // Light gray text for better contrast

	// Загружаем шрифт
	if err := dc.LoadFontFace("assets/NotoSans-Regular.ttf", fontSize); err != nil {
		fmt.Println("Error loading font:", err)

		return nil, fmt.Errorf("error loading font: %v", err)
	}

	// Рисуем заголовок таблицы
	dc.SetColor(color.RGBA{255, 255, 255, 255}) // White text for header
	dc.DrawStringAnchored("Турнирная таблица", float64(width/2), float64(padding)+20, 0.5, 0.5)

	// Рисуем шапку таблицы
	headers := []string{"#", "Команда", "И", "В", "Н", "П", "ГЗ", "ГП", "РГ", "О"}
	x := padding
	y := headerHeight + padding
	dc.SetColor(color.RGBA{40, 40, 40, 255}) // Slightly lighter background for header
	dc.DrawRectangle(0, float64(headerHeight), float64(width), float64(rowHeight)-5)
	dc.Fill()
	dc.SetColor(color.RGBA{230, 230, 230, 255}) // Light gray text for headers

	for i, header := range headers {
		dc.DrawStringAnchored(header, float64(x+colWidths[i]/2), float64(y)+7.5, 0.5, 0.5)
		x += colWidths[i]
	}

	// Рисуем таблицу
	y += rowHeight
	for i, row := range data {
		x = padding
		if i%2 == 1 {
			dc.SetColor(color.RGBA{30, 30, 30, 255}) // Slightly lighter background for alternating rows
			dc.DrawRectangle(0, float64(y-rowHeight/2), float64(width), float64(rowHeight))
			dc.Fill()
			dc.SetColor(color.RGBA{230, 230, 230, 255}) // Light gray text for content
		}
		cells := []string{
			fmt.Sprintf("%d", row.Position),
			row.Team.Name,
			fmt.Sprintf("%d", row.PlayedGames),
			fmt.Sprintf("%d", row.Won),
			fmt.Sprintf("%d", row.Draw),
			fmt.Sprintf("%d", row.Lost),
			fmt.Sprintf("%d", row.GoalsFor),
			fmt.Sprintf("%d", row.GoalsAgainst),
			fmt.Sprintf("%d", row.GoalDifference),
			fmt.Sprintf("%d", row.Points),
		}
		for j, cell := range cells {
			// Обработка длинных названий команд
			if j == 1 {
				// Разбиваем длинное название на несколько строк
				maxWidth := float64(colWidths[j]) - padding*2
				words := strings.Fields(cell)
				var lines []string
				currentLine := ""
				for _, word := range words {
					testLine := currentLine
					if currentLine != "" {
						testLine += " "
					}
					testLine += word
					w, _ := dc.MeasureString(testLine)
					if w > maxWidth {
						lines = append(lines, currentLine)
						currentLine = word
					} else {
						currentLine = testLine
					}
				}
				lines = append(lines, currentLine)
				// Рисуем строки
				for k, line := range lines {
					dc.DrawStringAnchored(line, float64(x+colWidths[j]/2), float64(y)+float64(k*fontSize), 0.5, 0.5)
				}
				// Увеличиваем высоту строки, если название длинное
				if len(lines) > 1 {
					y += int(fontSize * (len(lines) - 1))
				}
			} else {
				dc.DrawStringAnchored(cell, float64(x+colWidths[j]/2), float64(y), 0.5, 0.5)
			}
			x += colWidths[j]
		}
		y += rowHeight
	}
	// Сохраняем изображение в буфер

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}
	return buf, nil
}
