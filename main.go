package main

import (
	"io/ioutil"
	"os"
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

	data, errRead := ioutil.ReadAll(csvFile)
	if errRead != nil {
		panic(errRead)
	}

	htmlFile, errOpen2 := os.OpenFile("output.html", os.O_CREATE|os.O_RDWR, 0777)
	defer func() {
		htmlFile.Close()
	}()
	if errOpen2 != nil {
		panic(errOpen2)
	}
	WriteHeader(htmlFile)
	WriteBody(htmlFile, data)
}

func WriteHeader(file *os.File) {
	file.WriteString(
		"<!DOCTYPE html>\n<html" +
			" lang=\"en\">\n<head>\n" +
			"    <meta charset=\"UTF-8\">\n" +
			"    <title>Title</title>\n</head>\n",
	)
}

func WriteBody(file *os.File, data []byte) {
	rows := strings.Split(string(data), "\n")

	file.WriteString("<body>\n")
	file.WriteString("<table border=\"1\">\n")

	ths := strings.Split(rows[0], ",")

	file.WriteString("\t<tr>")
	for k := range ths {
		file.WriteString("<th>")
		file.WriteString(ths[k])
		file.WriteString("</th>")
	}
	file.WriteString("\t</tr>\n")

	columnsNum := len(ths)
	/*
		Переменная columnsNum нужна далее для определения валидности склейки строк,
		предполагаем, что названия колонок написаны без ошибок и спецсимволов
	*/
	for i := 1; i < len(rows)-1; i++ {
		tds := strings.Split(rows[i], ",")

		file.WriteString("\t<tr>")
		for k := 0; k < len(tds); k++ {
			td := tds[k]
			file.WriteString("<td>")
			// разделение по запятым могло разделить данные, заключенные в кавычки
			t := 0
			if tds[k][0] == '"' && lastSymbol(tds[k]) != '"' {
				// склеиваем строки, пока не найдем ту, которая заканчивается кавычками
				for t < len(tds) && tds[t][0] != '"' && lastSymbol(tds[t]) == '"' && columnsNum > len(tds)-t {
					td += tds[t]
					t++
				}
			}
			if t == len(tds) {
				file.WriteString(tds[k])
			} else {
				file.WriteString(td)
			}
			file.WriteString("</td>")
		}
		file.WriteString("</tr>\n")
	}

	file.WriteString("</table>\n")
	file.WriteString("</body>\n")
	file.WriteString("</html>\n")
}

func lastSymbol(s string) byte {
	return s[len(s)-1]
}
