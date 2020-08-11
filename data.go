package rpt

import (
	"encoding/json"
	"fmt"
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

	return ToJSON(dbds)
}

type DBQueryDataSet struct {
	Name  string
	Query string
}

func (dbqds *DBQueryDataSet) ToJson() []byte {
	output := ToJSON(dbqds)
	return output
}

type DataTable struct {
	Columns     map[string]DataColumn
	ColSlice    []*DataColumn `json:"-"`
	Rows        []*DataRow
	Delimiter   string
	Constraints []string
}

type DataColumn struct {
	Header       string
	DataType     string
	Constraints  []string
	DefaultValue interface{}
}

type DataRow string

func (dr *DataRow) Print() string {
	return fmt.Sprint(dr)
}

func (dr *DataRow) Parse(del string) []string {
	return strings.Split(dr.Print(), del)
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
