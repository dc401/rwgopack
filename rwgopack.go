package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const xorKey = 0x42

func xorCipher(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ xorKey
	}
	return result
}

func packbin(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(data)
	w.Close()

	ciphered := xorCipher(compressed.Bytes())
	return ciphered, nil
}

func createSelfExtractingScript(cipheredData []byte, outputFilename string) error {
	script := fmt.Sprintf(`package main

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
)

const xorKey = 0x42

func main() {
	cipheredData, _ := hex.DecodeString("%x")
	deciphered := xorCipher(cipheredData)

	r, _ := zlib.NewReader(bytes.NewReader(deciphered))
	original, _ := ioutil.ReadAll(r)
	r.Close()

	tempFile, _ := ioutil.TempFile("", "packed_*.bin")
	tempFile.Write(original)
	tempFile.Close()

	os.Chmod(tempFile.Name(), 0755)
	cmd := exec.Command(tempFile.Name(), os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	os.Remove(tempFile.Name())
}

func xorCipher(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ xorKey
	}
	return result
}
`, cipheredData)

	err := ioutil.WriteFile(outputFilename+".go", []byte(script), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-o", outputFilename, outputFilename+".go")
	err = cmd.Run()
	if err != nil {
		return err
	}

	os.Remove(outputFilename + ".go")
	fmt.Printf("Self-extracting binary created: %s\n", outputFilename)
	fmt.Printf("Run it with: ./%s\n", outputFilename)

	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run rwgopack.go <input_file> <output_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	cipheredData, err := packbin(inputFile)
	if err != nil {
		fmt.Printf("Error packing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ciphered data size: %d bytes\n", len(cipheredData))

	err = createSelfExtractingScript(cipheredData, outputFile)
	if err != nil {
		fmt.Printf("Error creating self-extracting script: %v\n", err)
		os.Exit(1)
	}
}
