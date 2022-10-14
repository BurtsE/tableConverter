package main

import (
	"regexp"
	"testing"
)

func TestCreateCell(t *testing.T) {

	data := []byte("abcф12")
	t.Logf("Initialized cell info: %s", data)
	re, err := regexp.Compile(`^\s*<td>.*</td>\s*$`)
	if err != nil {
		t.Fatal(err)
	}
	str := CreateTableCell(data)
	t.Logf("Input string: %s", data)
	if re.FindString(str) == "" {
		t.Logf("Expected output:\t<td>%s </td>", data)
		t.Errorf("Given output: %s", str)
	}
}

func TestSepId(t *testing.T) {
	str := "абвгд  абвqwe123  "
	t.Logf("")
	sep := SepTd(str, 0, 13)
	if sep != "абвгд  абвqwe" {
		t.Errorf("Wrong separation: %s", sep)
	}
}

func TestParseCsv(t *testing.T) {
	var backupData []byte
	var countErrors int64
	testData := []byte("First       test string\n" +
		"Second    string here  ")
	columnLen = [6]int{10, 7, 6}
	result := ParseCsv(testData)

	for _, row := range result {
		if len(row) == 0 {
			t.Error("Empty rows in result")
		}
		for _, cell := range row {
			if len(cell) == 0 {
				t.Error("Empty cells in result")
			}
			backupData = append(backupData, cell...)
		}
		backupData = append(backupData[:len(backupData)-1], '\n')
	}
	for i := range testData {
		if testData[i] != backupData[i] {
			countErrors++
		}
	}
	if countErrors > 0 {
		t.Errorf("data changed while it was parsed. Byte mismatches: %d", countErrors)
	}
}
