package lspci

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type LSPCI struct {
	Data       map[string]map[string]string
	flagNumber bool
}

func New(vendorInNumber bool) *LSPCI {
	return &LSPCI{
		Data:       make(map[string]map[string]string),
		flagNumber: vendorInNumber,
	}
}

func (l *LSPCI) Run() error {
	bin, findErr := FindBin("lspci")
	if findErr != nil {
		return findErr
	}
	args := []string{"-vmm", "-D"}
	if l.flagNumber {
		args = append(args, "-n")
	}
	cmd := exec.Command(bin, args...)

	out := &bytes.Buffer{}
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		return err
	}
	l.Data, err = parseLSPCI(out)
	return err
}

var sep = []byte{'\n', '\n'}

func scanDoubleNewLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, sep); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func parseLSPCI(r io.Reader) (map[string]map[string]string, error) {
	ret := make(map[string]map[string]string)
	scanner := bufio.NewScanner(r)
	scanner.Split(scanDoubleNewLine)
	for scanner.Scan() {
		// Per sector
		section := make(map[string]string)
		subScanner := bufio.NewScanner(bytes.NewBuffer(scanner.Bytes()))
		for subScanner.Scan() {
			data := strings.SplitN(subScanner.Text(), ":\t", 2)
			section[data[0]] = data[1]
		}
		if err := subScanner.Err(); err != nil {
			return nil, err
		}
		ret[section["Slot"]] = section
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func FindBin(binary string) (string, error) {
	locations := []string{"/sbin", "/usr/sbin", "/usr/local/sbin"}

	for _, path := range locations {
		lookup := path + "/" + binary
		fileInfo, err := os.Stat(path + "/" + binary)

		if err != nil {
			continue
		}

		if !fileInfo.IsDir() {
			return lookup, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Unable to find the '%v' binary", binary))
}
