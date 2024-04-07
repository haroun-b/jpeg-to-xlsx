package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Fatalf(
			"Wrong usage!\nExample usage: %s <./dir/source-img.jpeg> [<./dir/output.xml>]",
			os.Args[0],
		)
	}

	imgPath, outputPath := os.Args[1], "./output.xml"
	if len(os.Args) == 3 {
		pathWithoutExt, _ := strings.CutSuffix(os.Args[2], ".xml")
		outputPath = fmt.Sprintf("%s.xml", pathWithoutExt)
	}

	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("Error Opening File %s\n%v", imgPath, err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatalf("Error Decoding JPEG\n%v", err)
	}

	bounds := img.Bounds()
	sheetWidth, height := bounds.Max.X, bounds.Max.Y
	sheetHeight := height * 3 // 1 row for each color channel

	colors := make(map[string]bool)
	table := make([]string, 0, sheetHeight)

	redRow := make([]string, sheetWidth)
	greenRow := make([]string, sheetWidth)
	blueRow := make([]string, sheetWidth)

	for y := 0; y < height; y++ {
		for x := 0; x < sheetWidth; x++ {
			color := img.At(x, y)

			r, g, b, _ := color.RGBA()
			rx, gx, bx := fmt.Sprintf("%x0000", r>>8), fmt.Sprintf("00%x00", g>>8), fmt.Sprintf("0000%x", b>>8)

			colors[rx] = true
			colors[gx] = true
			colors[bx] = true

			redRow[x] = fmt.Sprintf("<Cell ss:StyleID=\"%s\" />", rx)
			greenRow[x] = fmt.Sprintf("<Cell ss:StyleID=\"%s\" />", gx)
			blueRow[x] = fmt.Sprintf("<Cell ss:StyleID=\"%s\" />", bx)
		}

		table = append(table,
			fmt.Sprintf("<Row ss:AutoFitHeight=\"0\">\n%s\n</Row>", strings.Join(redRow, " ")),
			fmt.Sprintf("<Row ss:AutoFitHeight=\"0\">\n%s\n</Row>", strings.Join(greenRow, " ")),
			fmt.Sprintf("<Row ss:AutoFitHeight=\"0\">\n%s\n</Row>", strings.Join(blueRow, " ")),
		)

	}

	var styleBuilder strings.Builder
	for color := range colors {
		styleBuilder.WriteString(fmt.Sprintf("<Style ss:ID=\"%s\">", color))
		styleBuilder.WriteString(fmt.Sprintf("<Interior ss:Color=\"#%s\" ss:Pattern=\"Solid\" />", color))
		styleBuilder.WriteString("</Style>\n")

	}

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0744); err != nil {
		log.Fatalf("Error Creating Output Directory %s\n%v", outputDir, err)
	}

	err = os.WriteFile(outputPath, []byte(fmt.Sprintf(`<?xml version="1.0"?>
<?mso-application progid="Excel.Sheet"?>
<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet"
	xmlns:o="urn:schemas-microsoft-com:office:office"
	xmlns:x="urn:schemas-microsoft-com:office:excel"
	xmlns:ss="urn:schemas-microsoft-com:office:spreadsheet"
	xmlns:html="http://www.w3.org/TR/REC-html40">
<Styles>
%s
</Styles>
<Worksheet ss:Name="Book1">
<Table ss:ExpandedColumnCount="%d" ss:ExpandedRowCount="%d" x:FullColumns="1" x:FullRows="1" ss:DefaultColumnWidth="65" ss:DefaultRowHeight="16">
%s
</Table>
</Worksheet>
</Workbook>`, styleBuilder.String(), sheetWidth, sheetHeight, strings.Join(table, "\n"))), 0644)

	if err != nil {
		log.Fatalf("Error Writing output file %s\n%v", outputPath, err)
	}
}
