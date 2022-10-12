package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	csvFile, errOpen1 := os.Open("data.csv")
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

	outCsv, errOpen2 := os.OpenFile("outCsv.html", os.O_CREATE|os.O_RDWR, 0777)
	defer func() {
		outCsv.Close()
	}()
	if errOpen2 != nil {
		panic(errOpen2)
	}

	WriteHtml(outCsv, csvData, ".csv")

	prnFile, errOpen3 := os.Open("data.prn")
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

	outPrn, errOpen4 := os.OpenFile("outPrn.html", os.O_CREATE|os.O_RDWR, 0777)
	defer func() {
		outPrn.Close()
	}()
	if errOpen4 != nil {
		panic(errOpen4)
	}
	WriteHtml(outPrn, prnData, ".prn")
}

func WriteHtml(file *os.File, data []byte, docType string) {

	var re *regexp.Regexp
	var sepSymbol byte
	var border string
	var html string
	switch docType {
	case ".csv":
		re, _ = regexp.Compile(`".+?",|.+?,`)
		sepSymbol = ','
		border = "1"
	case ".prn":
		re, _ = regexp.Compile(`.+ `)
		sepSymbol = ' '
		border = "0"
	default:
		panic("unknown format")
	}
	html = "<!DOCTYPE html>\n" +
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
		"</style>\n" +
		"<body>\n" +
		fmt.Sprintf("<table border=\"%s\">\n", border)

	rows := strings.Split(string(data), "\n")

	for i := 0; i < len(rows)-1; i++ {
		r := []byte(rows[i])

		r = append(r, sepSymbol)
		tds := re.FindAll(r, -1)

		html += "\t<tr>"
		for k := range tds {
			html += "<td>" + string(tds[k][:len(tds[k])-1]) + "</td>"
		}
		html += "</tr>\n"
	}
	html += "</table>\n" +
		"</body>\n" +
		"</html>\n"

	file.WriteString(html)
}
