/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Shitomo/mybatis-mapper-scanner/mapperxml"
	"github.com/spf13/cobra"
)

type Mapper struct {
	Selects []Select `xml:"select"`
}
type Select struct {
	ID  string `xml:"id,attr"`
	SQL string `xml:",innerxml"`
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "get mapper method which access specify tables",
	Run: func(cmd *cobra.Command, args []string) {
		tablesStr, err := cmd.Flags().GetString("tables")
		if err != nil || tablesStr == "" {
			log.Printf("tables is required")
		}

		dirPath, err := cmd.Flags().GetString("dir-path")
		if err != nil || dirPath == "" {
			log.Panicf("dir-path is required")
		}

		tables := strings.Split(tablesStr, ",")

		// get all mapper file
		xmlFiles, err := mapperxml.GetMapperXmlFiles(dirPath)
		if err != nil {
			log.Panic(err)
		}

		// parse mapper
		var mappers []Mapper
		for _, file := range xmlFiles {
			file, err := os.Open(file)
			if err != nil {
				log.Panic(err)
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				log.Panic(err)
			}

			mapper := Mapper{}
			err = xml.Unmarshal(data, &mapper)
			if err != nil {
				log.Panic(err)
			}

			mappers = append(mappers, mapper)

			file.Close()
		}

		var result []string
		for _, m := range mappers {
			for _, s := range m.Selects {
				for _, t := range tables {
					if strings.Contains(s.SQL, t) {
						result = append(result, s.ID)
						continue
					}
				}
			}
		}

		log.Printf("result: %v", result)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().String("tables", "", "target tables")
	scanCmd.Flags().String("dir-path", "", "mapper file path")
}
