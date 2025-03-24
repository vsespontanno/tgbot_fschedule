package handlers

import (
	"fmt"
	"football_tgbot/types"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

func GenerateTableImage(data []types.Standing, filename string) error {
	// Константы
	const (
		width        = 780
		height       = 920
		padding      = 10
		rowHeight    = 40
		headerHeight = 50
		fontSize     = 20
		numCols      = 10
	)

	// Динамически определяем ширину столбцов
	colWidths := []int{
		50,  // #
		300, // Команда
		50,  // И
		50,  // В
		50,  // Н
		50,  // П
		50,  // ГЗ
		50,  // ГП
		50,  // РГ
		50,  // О
	}

	dc := gg.NewContext(width, height)

	// Задаем фон
	dc.SetColor(color.White)
	dc.Clear()

	// Задаем цвет текста
	dc.SetColor(color.Black)

	// Загружаем шрифт
	if err := dc.LoadFontFace("/usr/share/fonts/noto/NotoSans-Regular.ttf", fontSize); err != nil {
		fmt.Println("Error loading font:", err)
		return err
	}

	// Рисуем заголовок таблицы
	dc.SetColor(color.RGBA{0, 0, 0, 255})
	dc.DrawStringAnchored("Турнирная таблица", float64(width/2), float64(padding)+20, 0.5, 0.5)

	// Рисуем шапку таблицы
	headers := []string{"#", "Команда", "И", "В", "Н", "П", "ГЗ", "ГП", "РГ", "О"}
	x := padding
	y := headerHeight + padding
	dc.SetColor(color.RGBA{200, 200, 200, 255})
	dc.DrawRectangle(0, float64(headerHeight), float64(width), float64(rowHeight))
	dc.Fill()
	dc.SetColor(color.Black)

	for i, header := range headers {
		dc.DrawStringAnchored(header, float64(x+colWidths[i]/2), float64(y), 0.5, 0.5)
		x += colWidths[i]
	}

	// Рисуем таблицу
	y += rowHeight
	for i, row := range data {
		x = padding
		if i%2 == 1 {
			dc.SetColor(color.RGBA{240, 240, 240, 255})
			dc.DrawRectangle(0, float64(y-rowHeight/2), float64(width), float64(rowHeight))
			dc.Fill()
			dc.SetColor(color.Black)
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

	// Сохраняем изображение
	return dc.SavePNG(filename)
}
