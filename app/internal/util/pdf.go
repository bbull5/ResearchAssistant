package util

import (
	"bytes"
	"io"
	"os"

	"github.com/ledongthuc/pdf"
)


func ExtractTextFromPDF(path string) (string, error) {
	f, rErr := os.Open(path)
	if rErr != nil {
		return "", rErr
	}
	defer f.Close()

	reader, err := pdf.NewReader(f, getSize(f))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(&buf, b)
	return buf.String(), err
}

func getSize(file *os.File) int64 {
	fi, _ := file.Stat()
	return fi.Size()
}