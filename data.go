package rpt

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type DataSet interface {
	ToJson() []byte
}

type DBDataSet struct {
	Tables map[string]DataTable
	Name   string
}

func (dbds *DBDataSet) ToJson() []byte {

	output, _ := json.MarshalIndent(dbds, "", "  ")

	return output
}

type DataTable struct {
	Columns map[string]DataColumn
	Row     []*DataRow
}

type DataColumn struct {
	Header       string
	DataType     string
	Nullable     bool
	Delimiter    string
	DefaultValue interface{}
}

type DataRow struct {
	Data string
}

func (dr *DataRow) Parse(del string) []string {
	return strings.Split(dr.Data, del)
}

func ImportDBDataSet(filePath string) (*DBDataSet, []error) {

	errs := []error{}
	ds := &DBDataSet{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		errs = append(errs, err)
	}
	defer jsonFile.Close()

	content, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		errs = append(errs, err)
	}

	err = json.Unmarshal(content, ds)
	if err != nil {
		errs = append(errs, err)
	}

	return ds, errs
}
