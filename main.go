package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	csvFile, errOpen1 := os.Open("testdata.csv")
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

	WriteHtml(htmlFile, data, ".csv")
	//WriteHtml(htmlFile, data, ".prn")
}

func WriteHtml(file *os.File, data []byte, doctype string) {

	var re *regexp.Regexp
	switch doctype {
	case ".csv":
		re, _ = regexp.Compile(`".+?",|.+?,`)
	case ".prn":
		re, _ = regexp.Compile(`.+? +`)
	}

	file.WriteString(
		"<!DOCTYPE html>\n<html" +
			" lang=\"en\">\n<head>\n" +
			"    <meta charset=\"UTF-8\">\n" +
			"    <title>Title</title>\n</head>\n",
	)
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

	for i := 1; i < len(rows)-1; i++ {
		//tds := strings.Split(rows[i], ",")
		r := []byte(rows[i])
		r = append(r, ',')
		tds := re.FindAll(r, -1)
		if len(tds) != len(ths) {
			panic("bad data")
		}
		file.WriteString("\t<tr>")
		for k := 0; k < len(tds); k++ {

			file.WriteString("<td>")
			file.Write(tds[k][:len(tds[k])-1])
			file.WriteString("</td>")
		}
		file.WriteString("</tr>\n")
	}

	file.WriteString("</table>\n")
	file.WriteString("</body>\n")
	file.WriteString("</html>\n")
}
