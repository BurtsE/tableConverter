package main

import "os"

func main() {
	csvFile, errOpen1 := os.Open("data.csv")
	defer func() {
		csvFile.Close()
	}()
	if errOpen1 != nil {
		panic(errOpen1)
	}

	htmlFile, errOpen2 := os.OpenFile("output.html", os.O_CREATE|os.O_RDWR, 0777)
	defer func() {
		htmlFile.Close()
	}()
	if errOpen2 != nil {
		panic(errOpen2)
	}
	WriteHeader(htmlFile)
	WriteBody(htmlFile, []byte(""))
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
	file.WriteString("<body>\n")

	file.WriteString("</body>\n")
}
