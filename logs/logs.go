package logs

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

var LogsBuffer bytes.Buffer

func StartLogging() {
	// Créer un pipe pour rediriger stdout et stderr qui contiennent tous les outputs fmt et log
	// c'est du code chatGPT

	reader, writer, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}

	stdout := os.Stdout

	os.Stdout = writer
	os.Stderr = writer

	multiWriter := io.MultiWriter(&LogsBuffer, stdout)

	log.SetOutput(multiWriter)

	go func() {
		io.Copy(multiWriter, reader)
	}()
}

func ResetBuffer() {
	fmt.Println("Reset Buffer")
	for {
		if bytes.Count(LogsBuffer.Bytes(), []byte{'\n'}) > 70 {
			LogsBuffer.Reset()
		}
	}
}
