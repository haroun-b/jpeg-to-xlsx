package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func BundleXLSX(sourceDir, outputPath string) error {
	xlsxFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer xlsxFile.Close()

	xlsxWriter := zip.NewWriter(xlsxFile)
	defer xlsxWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)
		header.Method = zip.Store // No compression

		writer, err := xlsxWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
