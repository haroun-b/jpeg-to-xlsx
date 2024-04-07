package xlsxstarter

import (
	"os"
	"path/filepath"
)

/*
the structure of an uncompressed xlsx file is as follows:

|- [Content_Types].xml
|- _rels
	|- .rels
|- xl
  |- styles.xml
	|- workbook.xml
	|- worksheets
		|- sheet1.xml
	|- _rels
		|- workbook.xml.rels

the only files we need to modify are xl/styles.xml and xl/worksheets/sheet1.xml. the reset can be copied as is.
*/

func CreateXLSXStarter(path string) error {

	err := os.Mkdir(filepath.Join(path, "_rels"), 0744)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(path, "xl", "_rels"), 0744)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(path, "xl", "worksheets"), 0744)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "[Content_Types].xml"), []byte(contentTypesXML), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "_rels", ".rels"), []byte(dotRelsXML), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "xl", "_rels", "workbook.xml.rels"), []byte(workbookRelsXML), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "xl", "workbook.xml"), []byte(workbookXML), 0644)
	if err != nil {
		return err
	}

	return nil
}
