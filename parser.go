package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type filedata struct {
	binding string
	action  string
}

func parse(file string, data *[]filedata) (err error) {
	if data == nil {
		return errors.New("empty value was passed to parse function")
	}

	configFile, err := os.Open(file)

	if err != nil {
		return
	}

	defer configFile.Close()

	reader := bufio.NewReader(configFile)

	lineNumber := 0
	wasKeybinding := false
	wasPrefix := false
	type datumType struct {
		binding strings.Builder
		action  strings.Builder
	}
	datum := []datumType{}
	index := 0

	for {
		lineNumber++

		var (
			line     []byte
			isPrefix bool
		)

		line, isPrefix, err = reader.ReadLine()

		if err != nil {
			break
		}

		if index+1 != len(datum) {
			datum = append(datum, datumType{})
		}

		lineStr := string(line)

		// skip the shebang
		if lineNumber == 1 && strings.HasPrefix(lineStr, "#!") {
			continue
		}

		// skip an empty line
		if lineStr == "" {
			continue
		}

		// ignore comments (##+)
		if strings.HasPrefix(lineStr, "##") {
			continue
		}

		// decide whether the line is a keybinding or not
		if strings.HasPrefix(lineStr, "#") {
			if isPrefix {
				log.Fatalf("a keybinding can't be that long, line %d, file %s", lineNumber, file)
				os.Exit(1)
			}
			// erase spaces for key validation
			lineStr = strings.ReplaceAll(lineStr, " ", "")
			validator := regexp.MustCompile(`(?m)^#\w+((\+\w+)+)?$`)
			if validator.MatchString(lineStr) {
				if datum[index].action.Len() != 0 {
					index++
					datum = append(datum, datumType{})
				}
				lineStr = strings.ReplaceAll(strings.TrimPrefix(lineStr, "#"), "+", "-")
				// overwrite previous prefix
				if wasKeybinding {
					if datum[index].binding.Len() != 0 {
						datum[index].binding.Reset()
						fmt.Println(fmt.Sprintf("overwriting %d", lineNumber))
						fmt.Println(fmt.Sprintf("previous - \"%s\"", datum[index].binding.String()))
						fmt.Println(fmt.Sprintf("new - \"%s\"", lineStr))
					}
				}
				datum[index].binding.Write([]byte(lineStr))
				wasKeybinding = true
			}
		} else {
			wasKeybinding = false
			if isPrefix {
				if wasPrefix {
					datum[index].action.Write(line)
				} else {
					if datum[index].action.Len() != 0 {
						datum[index].action.Write([]byte("\n"))
					}
					datum[index].action.Write(line)
					wasPrefix = true
				}
				continue
			}

			if wasPrefix {
				datum[index].action.Write(line)
				wasPrefix = false
			} else {
				if datum[index].action.Len() != 0 {
					datum[index].action.Write([]byte("\n"))
				}
				datum[index].action.Write(line)
			}
		}
	}

	if err == io.EOF {
		err = nil
	} else {
		return
	}

	*data = nil
	for _, d := range datum {
		*data = append(*data, filedata{binding: d.binding.String(), action: d.action.String()})
	}

	return
}
