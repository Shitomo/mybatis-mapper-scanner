package mapperxml

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetMapperXmlFiles(dirPath string) ([]string, error) {
	var xmlFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if strings.Contains(path, ".xml") {
			xmlFiles = append(xmlFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return xmlFiles, nil
}
