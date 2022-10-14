package main

import (
	"fmt"
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

// TestSepId выполняет проверку корректной обработки строки формата utf-8
func TestSepId(t *testing.T) {
	t.Log("")
	str := "абвгд  абвqwe123  "
	t.Logf("")
	sep := SepTd(str, 3, 13)
	if sep != "гд  абвqwe" {
		t.Errorf("Wrong separation: %s", sep)
	}
}

// TestParseCsv Выполняет проверку отсутствия пустых полей после парсинга, а
// также проводится проверка то, были ли данные изменены в процессе парсинга
func TestParsePrn(t *testing.T) {
	var backupData []byte
	var countErrors int64
	testData := []byte("First       test string\n" +
		"Second    string here  ")
	columnLen = [6]int{10, 7, 6}
	result := ParseCsv(testData)

	t.Logf("Test Data: %s", testData)
	t.Log("Iterating through result")
	for _, row := range result {
		if len(row) == 0 {
			t.Fatal("Empty rows in result")
		}
		for _, cell := range row {
			if len(cell) == 0 {
				t.Fatal("Empty cells in result")
			}
			backupData = append(backupData, cell...)
		}
		backupData = append(backupData, '\n')
	}
	t.Logf("generated backup: %s", backupData)
	t.Log("Counting mismatches")
	for i := range testData {
		if testData[i] != backupData[i] {
			countErrors++
		}
	}
	if countErrors > 0 {
		t.Errorf("Data changed while it was parsed. Byte mismatches: %d", countErrors)
	}
}

func TestParseCsv(t *testing.T) {
	testData := []byte("\"Jack\",\"Via Rocco Chinnici 4d\",3423 ba,0313-111475,22,05/04/1984\n" +
		"\"Charlie\",\"Via Aldo Moro, 7\",3209 DD,30-34563332,343.8,04/10/1954")

	result := ParseCsv(testData)

	expectedResult := []string{"\"Jack\"", "\"Via Rocco Chinnici 4d\"", "3423 ba", "0313-111475", "22", "05/04/1984",
		"\"Charlie\"", "\"Via Aldo Moro, 7\"", "3209 DD", "30-34563332", "343.8", "04/10/1954"}
	var i int
	for _, row := range result {
		if len(row) == 0 {
			t.Fatal("Empty rows in result")
		}
		for _, cell := range row {
			if len(cell) == 0 {
				t.Fatal("Empty cells in result")
			} else if string(cell) != expectedResult[i] {
				fmt.Println(cell, []byte(expectedResult[i]))
				t.Errorf("Wrong parsing: expected %s, got %s", expectedResult[i], cell)
			}
			i++
		}
	}
}
