package gpio

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
)

func read(file string) (string, error) {
	data, err := ioutil.ReadFile(file)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), err
}

func write(file string, data string) error {
	f, err := os.OpenFile(file, os.O_WRONLY, 0644)
	defer f.Close()

	if err == nil {
		_, err = f.WriteString(data)
	}

	return err
}

func exportPin(baseDir string, pin int) error {
	return write(fmt.Sprintf("%s/export", baseDir),
		fmt.Sprintf("%d", pin))
}

func unexportPin(baseDir string, pin int) error {
	return write(fmt.Sprintf("%s/unexport", baseDir),
		fmt.Sprintf("%d", pin))
}

func setDirection(baseDir string, pin int, direction int) error {
	pinDirection := map[int]string{1: "out", 0: "in"} // should/could be a const?

	return write(fmt.Sprintf("%s/gpio%d/direction", baseDir, pin),
		pinDirection[direction])
}

func readValue(baseDir string, pin int) (string, error) {
	return read(fmt.Sprintf("%s/gpio%d/value", baseDir, pin))
}

func writeValue(baseDir string, pin int, value int) error {
	return write(fmt.Sprintf("%s/gpio%d/value", baseDir, pin),
		fmt.Sprintf("%d", value))
}
