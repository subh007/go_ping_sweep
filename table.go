package go_ping_sweep

import (
	"errors"
	"fmt"
	"strings"
)

// struct to create the table
type Table struct {
	Title  string     //table title
	Header []string   // table header
	Data   [][]string // table row * column
	tail   int        // hold the index for the end data row
}

// function to set the table title
func (t *Table) SetTitle(title string) {
	t.Title = title
}

// function to set the table header.
func (t *Table) SetHeader(header ...string) {
	t.Header = header
}

// function to add the data to the table.
func (t *Table) AddData(data ...string) error {

	// if the header is not set return err
	if t.Header == nil {
		return errors.New("Header is not set")
	}

	if t.Data == nil {
		t.Data = make([][]string, 10, 100)
	}

	// add the data to the table data
	if data != nil {
		t.Data[t.tail] = make([]string, len(t.Header))
		for i := 0; i < len(data) && i < len(t.Header); i++ {
			t.Data[t.tail][i] = data[i]
		}
		t.tail++
	}
	return nil
}

//funtion print the data in table format
func (t *Table) CreateTable() {

	// print title
	if t.Title != "" {
		fmt.Println("=========" + t.Title + "==============")
	}
	// print the header first
	fmt.Println("")
	for i := 0; i < len(t.Header); i++ {
		fmt.Print(t.Header[i] + " | ")
	}

	fmt.Println("")

	for i := 0; i < t.tail; i++ {
		for j := 0; j < len(t.Header); j++ {
			fmt.Print(t.Data[i][j] + " | ")
		}
		fmt.Println("")
	}
	fmt.Println("")
}

// function print the data in table
func (t *Table) PrintTable() {
	for i := 0; i < t.tail; i++ {
		for j := 0; j < len(t.Header); j++ {
			fmt.Print(t.Data[i][j])
		}
	}
}

// function print the header data
func (t *Table) PrintHeader() {
	for i := 0; i < len(t.Header); i++ {
		fmt.Println(t.Header[i])
	}
}

// return the header indes if exist
// ex if table is:
// TimePing | DataSize | PacketSize | status |
// then getColumnIndes("Status") will be 4 otherwise return -1.
func (t *Table) getCoulumnIndex(header string) int {
	if t != nil {
		for i := 0; i < t.tail; i++ {
			if strings.Compare(t.Header[i], header) == 0 {
				return i
			}
		}
	}
	return -1
}

// function returns the the column.
func (t *Table) GetColumn(columHeader string) []string {
	headerIndex := t.getCoulumnIndex(columHeader)

	columnData := make([]string, 10, 100)

	if headerIndex != -1 {
		// retrieve the all the element from the column.
		for i:=0; i < t.tail; i++ {
			columnData[i] = t.Data[i][headerIndex]
			//fmt.Print(t.Data[i][headerIndex])
		}
	}
	return columnData
}
