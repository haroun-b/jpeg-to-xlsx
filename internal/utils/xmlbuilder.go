package utils

import (
	"fmt"
	"image"
	"math"
	"strings"
)

const colorChannelCount = 3
const baseFillsCount = 2
const baseXfsCount = 1

func ImgToXMLs(img image.Image) (sheet, styles string) {
	bounds := img.Bounds()
	sheetWidth, height := bounds.Max.X, bounds.Max.Y
	sheetHeight := height * colorChannelCount // 3 rows (RGB) per pixel

	xfsIds := make(map[string]int)
	colors := make([]string, 0, 255*colorChannelCount) // a max of 255 colors per channel
	table := make([]string, 0, sheetHeight)

	colNames := make([]string, 0, sheetWidth)
	for i := 1; i < sheetWidth+1; i++ {
		colNames = append(colNames, numberToAlpha(i))
	}

	rgbRows := [colorChannelCount][]string{
		make([]string, sheetWidth),
		make([]string, sheetWidth),
		make([]string, sheetWidth),
	}
	latestRow := 1

	progress := 0
	for y := 0; y < height; y++ {
		progress = int(math.Ceil(float64(y) / float64(height) * 100))
		fmt.Printf("\r[%-100s] %d%%", strings.Repeat("#", progress), progress)

		for x := 0; x < sheetWidth; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			rgbHex := [colorChannelCount]string{
				fmt.Sprintf("ff%s0000", colorChannelToHex(r)),
				fmt.Sprintf("ff00%s00", colorChannelToHex(g)),
				fmt.Sprintf("ff0000%s", colorChannelToHex(b)),
			}

			for i, colorX := range rgbHex {
				colorIdx, found := xfsIds[colorX]

				if !found {
					colorIdx = len(colors) + baseXfsCount
					colors = append(colors, colorX)
					xfsIds[colorX] = colorIdx
				}

				rgbRows[i][x] = fmt.Sprintf(
					`<c r="%s%d" s="%d" />`, colNames[x], latestRow+i, colorIdx,
				)
			}
		}

		for i, row := range rgbRows {
			table = append(table,
				fmt.Sprintf(
					`<row r="%d" spans="1:%d">%s</row>`,
					latestRow+i,
					sheetWidth,
					strings.Join(row, ""),
				),
			)
		}

		latestRow += colorChannelCount
	}
	fmt.Println()

	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetFormatPr baseColWidth="10" defaultRowHeight="16" /><sheetData>%s</sheetData></worksheet>`,
		strings.Join(table, ""),
	), colorsToStylesXML(colors)
}

func colorsToStylesXML(colors []string) string {
	fills := make([]string, 0, len(colors))
	xfs := make([]string, 0, len(colors))

	for i, color := range colors {
		fills = append(
			fills,
			fmt.Sprintf(
				`<fill><patternFill patternType="solid"><fgColor rgb="%s" /><bgColor indexed="64" /></patternFill></fill>`,
				color,
			),
		)

		xfs = append(
			xfs,
			fmt.Sprintf(
				`<xf numFmtId="0" fillId="%d" borderId="0" applyFill="1" />`,
				i+baseFillsCount,
			),
		)
	}

	return fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><fonts count="1"><font><sz val="12" /><name val="Arial" /><family val="2" /></font></fonts><fills count="%d"><fill><patternFill patternType="none" /></fill><fill><patternFill patternType="gray125" /></fill>%v</fills><borders count="1"><border><left /><right /><top /><bottom /><diagonal /></border></borders><cellXfs count="%d"><xf numFmtId="0" fillId="0" borderId="0" />%v</cellXfs></styleSheet>`,
		len(fills)+baseFillsCount,
		strings.Join(fills, ""),
		len(xfs)+baseXfsCount,
		strings.Join(xfs, ""),
	)
}

func numberToAlpha(num int) string {
	result := ""

	for num > 0 {
		num--
		result = string(rune('A'+num%26)) + result
		num /= 26
	}

	return result
}

func colorChannelToHex(ch uint32) string {
	hex := fmt.Sprintf("%x", ch>>8) // 16-bit color channel to 8-bit

	if len(hex) == 1 {
		return "0" + hex
	} else {
		return hex
	}
}
