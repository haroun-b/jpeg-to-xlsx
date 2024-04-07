package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/haroun-b/jpeg-to-xlsx/internal/utils"
	"github.com/haroun-b/jpeg-to-xlsx/internal/xlsxstarter"
)

const outputExt = ".xlsx"

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Fatalf(
			"Wrong usage!\nExample usage: %s <./dir/source-img.jpeg> [<./dir/output%s>]",
			os.Args[0],
			outputExt,
		)
	}

	imgPath, outputPath := os.Args[1], fmt.Sprintf("./output%s", outputExt)
	if len(os.Args) == 3 {
		pathWithoutExt, _ := strings.CutSuffix(os.Args[2], outputExt)
		outputPath = fmt.Sprintf("%s%s", pathWithoutExt, outputExt)
	}

	file, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("Error opening image file %s\n%v", imgPath, err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding JPEG\n%v", err)
	}

	sheet1XML, stylesXML := utils.ImgToXMLs(img)

	outputDir, outputFileName := filepath.Dir(outputPath), filepath.Base(outputPath)

	if err := os.MkdirAll(outputDir, 0744); err != nil {
		log.Fatalf("Error creating output directory %s\n%v", outputDir, err)
	}

	tempDir, err := os.MkdirTemp(outputDir, fmt.Sprintf(".$%s", outputFileName))
	if err != nil {
		log.Fatalf("Error creating temp directory %s\n%v", tempDir, err)
	}

	err = xlsxstarter.CreateXLSXStarter(tempDir)
	if err != nil {
		log.Fatalf("Error creating XLSX starter files\n%v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "xl", "worksheets", "sheet1.xml"), []byte(sheet1XML), 0644)
	if err != nil {
		log.Fatalf("Error creating sheet1.xml\n%v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, "xl", "styles.xml"), []byte(stylesXML), 0644)
	if err != nil {
		log.Fatalf("Error creating styles.xml\n%v", err)
	}

	err = utils.BundleXLSX(tempDir, outputPath)
	if err != nil {
		log.Fatalf("Error bundling XLSX file\n%v", err)
	}

	os.RemoveAll(tempDir)
	log.Printf("XLSX file created => %s", outputPath)
}
