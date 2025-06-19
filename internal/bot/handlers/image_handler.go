package handlers

import (
	"bytes"
	"fmt"
	"football_tgbot/internal/types"
	"image/color"
	"strings"

	"github.com/fogleman/gg"
)

func GenerateTableImage(data []types.Standing, filename string) error {
	fmt.Printf("Generating image for гафоиафгаиыфиалфи %d standings\n", len(data))
	if len(data) == 0 {
		return fmt.Errorf("no standings data provided")
	}

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
	if err := dc.LoadFontFace("app/fonts/NotoSans-Regular.ttf", fontSize); err != nil {
		fmt.Println("Error loading font:", err)

		return fmt.Errorf("error loading font: %v", err)
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

	// Сохраняем изображение
	return dc.SavePNG(filename)
}

// функция для генерации изображения расписания матчей
func GenerateScheduleImage(matches []types.Match) (*bytes.Buffer, error) {

	const (
		width        = 780
		height       = 920
		padding      = 10
		headerHeight = 60
		fontSize     = 20
		lineWidth    = 1.5
		rowHeight    = 60
	)

	var (
		backgroundColor   = color.RGBA{18, 18, 18, 255}
		textColor         = color.RGBA{230, 230, 230, 255}
		headerTextColor   = color.RGBA{255, 255, 255, 255}
		headerBgColor     = color.RGBA{40, 40, 40, 255}
		alternateRowColor = color.RGBA{30, 30, 30, 255}
		lineColor         = color.RGBA{60, 60, 60, 255}
	)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches data provided")
	}

	dc := gg.NewContext(width, height)

	// Задаем фон
	dc.SetColor(backgroundColor)
	dc.Clear()

	// Загружаем шрифт
	if err := dc.LoadFontFace("app/fonts/NotoSans-Regular.ttf", fontSize); err != nil {
		fmt.Println("Error loading font:", err)

		return nil, fmt.Errorf("error loading font: %v", err)
	}

	// Рисуем заголовок
	dc.SetColor(headerTextColor)
	dc.DrawStringAnchored("Расписание матчей", float64(width/2), float64(padding)+20, 0.5, 0.5)

	// Определяем заголовки и ширину колонок
	headers := []string{"Дата", "Время", "Матч"}
	colWidths := []int{120, 80, 560}

	y := headerHeight + padding

	// Рисуем шапку таблицы
	dc.SetColor(headerBgColor)
	dc.DrawRectangle(0, float64(headerHeight)-10, float64(width), float64(rowHeight)-5)
	dc.Fill()

	// Рисуем линии
	dc.SetColor(lineColor)
	dc.SetLineWidth(lineWidth)
	dc.DrawLine(0, float64(headerHeight+rowHeight-15), float64(width), float64(headerHeight+rowHeight-15))
	dc.Stroke()

	// Рисуем заголовки
	dc.SetColor(textColor)
	currentX := padding
	for i, header := range headers {
		dc.DrawStringAnchored(header, float64(currentX+colWidths[i]/2), float64(y)+5, 0.5, 0.5)
		currentX += colWidths[i]

		if i < len(headers)-1 {
			dc.SetColor(lineColor)
			dc.DrawLine(float64(currentX), float64(headerHeight), float64(currentX), float64(height))
			dc.Stroke()
			dc.SetColor(textColor)
		}
	}

	// Рисуем строки с матчами
	y += rowHeight
	for i, match := range matches {
		currentX = padding
		if i%2 == 1 {
			dc.SetColor(alternateRowColor)
			dc.DrawRectangle(0, float64(y-rowHeight/2), float64(width), float64(rowHeight))
			dc.Fill()
		}

		dc.SetColor(lineColor)
		dc.DrawLine(0, float64(y+rowHeight/2), float64(width), float64(y+rowHeight/2))
		dc.Stroke()

		dc.SetColor(textColor)

		cells := []string{
			match.UTCDate[0:10],  // Дата
			match.UTCDate[11:16], // Время
			fmt.Sprintf("%s - %s", match.HomeTeam.Name, match.AwayTeam.Name),
		}

		for j, cell := range cells {
			if j == 2 {
				// Обработка длинных названий команд
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

				for k, line := range lines {
					dc.DrawStringAnchored(line, float64(currentX+colWidths[j]/2), float64(y)+float64(k*fontSize), 0.5, 0.5)
				}
				if len(lines) > 1 {
					y += int(fontSize * (len(lines) - 1))
				}
			} else {
				dc.DrawStringAnchored(cell, float64(currentX+colWidths[j]/2), float64(y), 0.5, 0.5)
			}
			currentX += colWidths[j]
		}
		y += rowHeight
	}

	// Возвращаем буфер с изображением
	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		return nil, fmt.Errorf("error encoding image: %v", err)
	}
	return buf, nil
}
