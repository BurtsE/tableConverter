package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

//Константа для парсинга prn файла
var columnLen = [6]int{16, 22, 9, 15, 12, 9}

func main() {
	csvFile, errOpen1 := os.Open("data/data.csv")
	defer func() {
		csvFile.Close()
	}()
	if errOpen1 != nil {
		panic(errOpen1)
	}

	csvData, errRead := ioutil.ReadAll(csvFile)
	if errRead != nil {
		panic(errRead)
	}

	outCsv, errOpen2 := os.OpenFile("output/outCsv.html", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
	defer func() {
		outCsv.Close()
	}()
	if errOpen2 != nil {
		panic(errOpen2)
	}

	WriteHtml(outCsv, csvData, ".csv")

	prnFile, errOpen3 := os.Open("data/data.prn")
	defer func() {
		prnFile.Close()
	}()
	if errOpen3 != nil {
		panic(errOpen3)
	}

	prnData, errRead2 := ioutil.ReadAll(prnFile)
	if errRead2 != nil {
		panic(errRead2)
	}

	outPrn, errOpen4 := os.OpenFile("output/outPrn.html", os.O_CREATE|os.O_RDWR, 0777)
	defer func() {
		outPrn.Close()
	}()
	if errOpen4 != nil {
		panic(errOpen4)
	}
	WriteHtml(outPrn, prnData, ".prn")
}

// WriteHtml выполняет запись данных  в указанный файл
func WriteHtml(file *os.File, data []byte, docType string) {
	var border string
	var html string
	var table [][][]byte
	switch docType {
	case ".csv":
		table = ParseCsv(data)
		border = "1"
	case ".prn":
		table = ParsePrn(data)
		border = "0"
	default:
		panic("unknown format")
	}
	html += WriteHeader()
	html += "<body>\n"
	html += CreateTable(table, border)
	html += "</body>\n" + "</html>\n"

	_, err := file.WriteString(html)
	if err != nil {
		panic(err)
	}
}

// WriteHeader генерирует заголовок html файла
func WriteHeader() string {
	return "<!DOCTYPE html>\n" +
		"<html lang=\"en\">\n" +
		"<head>\n" +
		"    <meta charset=\"UTF-8\">\n" +
		"    <title>Title</title>\n" +
		"</head>\n" +
		"<style>\n" +
		"body {\n" +
		"\twhite-space: pre;\n" +
		"\tfont-family: Consolas;\n" +
		"}\n" +
		"</style>\n"
}

// CreateTable генерирует html код таблицы
func CreateTable(data [][][]byte, border string) (html string) {
	html += fmt.Sprintf("<table border=\"%s\">\n", border)
	for i := 0; i < len(data); i++ {
		tds := data[i]
		html += "\t<tr>"
		for k := range tds {
			html += CreateTableCell(tds[k])
		}
		html += "</tr>\n"
	}
	html += "</table>\n"
	return
}

// CreateTableCell создает клетки таблицы в формате html
func CreateTableCell(data []byte) string {
	return "<td>" + string(data) + "</td>"
}

// ParseCsv представляет данные, полученные из csv файла
// в виде двумерного массива слов (срезов байт)
func ParseCsv(data []byte) (parsedData [][][]byte) {
	var sepSymbol byte = ','
	re, _ := regexp.Compile(`".+?",|.+?,`)
	rows := strings.Split(string(data), "\n")

	for i := range rows {
		row := []byte(rows[i])
		row = append(row, sepSymbol)
		tds := re.FindAll(row, -1)
		for k := range tds {
			tds[k] = tds[k][:len(tds[k])-1]
		}
		parsedData = append(parsedData, tds)
	}
	return
}

// ParsePrnSimple представляет данные, полученные из prn файла
// в виде двумерного массива слов (срезов байт).
// Словом будет являться вся строка, в табличном представлении будет единственный столбец
func ParsePrnSimple(data []byte) (parsedData [][][]byte) {
	var sepSymbol byte = ' '
	re, _ := regexp.Compile(`.+ `)
	rows := strings.Split(string(data), "\n")

	for i := 0; i < len(rows)-1; i++ {
		row := []byte(rows[i])
		row = append(row, sepSymbol)
		tds := re.FindAll(row, -1)
		parsedData = append(parsedData, tds)
	}
	return
}

// ParsePrn представляет данные, полученные из prn файла
// в виде двумерного массива слов (срезов байт).
// Предполагается, что ширина колонок с данными (как и их количество) известна
// В противном случае логическое деление столбцов не представляется возможным
// Более универсальный метод parsePrnSimple дает такое же графическое представление, но данные объединены
// в единый столбец таблицы
func ParsePrn(data []byte) (parsedData [][][]byte) {
	//columnLen := []int{16, 22, 9, 15, 12, 9}
	var sepSymbol byte = ' '
	re, _ := regexp.Compile(`.+ `)
	rows := strings.Split(string(data), "\n")

	for i := 0; i < len(rows)-1; i++ {
		var tRow [][]byte
		var startPos, endPos int
		row := string(append([]byte(rows[i]), sepSymbol))
		for k := 0; k < len(columnLen); k++ {
			endPos += columnLen[k]
			tRow = append(tRow, []byte(re.FindString(SepTd(row, startPos, endPos))))
			startPos += columnLen[k]
		}
		parsedData = append(parsedData, tRow)
	}
	return
}

//	SepTd нужна для правильного деления данных в случае, если данные
//	представлены не в ASCII (range итерируется по рунам), т.к. функция ParsePrn
//  использует ширину строки в символах
func SepTd(str string, start, end int) (cellData string) {
	var i int
	for _, symbol := range str {
		if i >= start && i < end {
			cellData += string(symbol)
		}
		i++
	}
	return
}
