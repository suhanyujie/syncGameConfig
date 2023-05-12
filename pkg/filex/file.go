package filex

import (
	"bytes"
	"io"
	"log"
	"os"
)

func ReadFile(file string) string {
	tmpFs, err := os.OpenFile(file, os.O_RDONLY, 0666)
	defer tmpFs.Close()
	if err != nil {
		log.Printf("[readFile] OpenFile err: %v", err)
		return ""
	}
	// 10M
	buf := make([]byte, 0, 1024*1024*10)
	writer := bytes.NewBuffer(buf)
	if _, err = io.Copy(writer, tmpFs); err != nil {
		log.Printf("[readFile] io.Copy err: %v", err)
		return ""
	}
	return writer.String()
}
