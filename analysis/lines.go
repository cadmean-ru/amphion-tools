package analysis

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type LinesInfo struct {
	Total, NotEmpty, Comments, Code int
}

func CountLines(path string) LinesInfo {
	var total, notEmpty, comments, code int

	var excluded = make([]string, 0, 10)
	if len(os.Args) > 4 {
		for i := 4; i < len(os.Args); i++ {
			excluded = append(excluded, filepath.Clean(os.Args[i]))
		}
	}

	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, ex := range excluded {
			if ex == path {
				fmt.Printf("Ignored: %s\n", path)
				return nil
			}
			dirPath, fileName := filepath.Split(path)
			cleanDirPath := filepath.Clean(dirPath)
			if strings.HasPrefix(cleanDirPath, ex) || fileName == ex {
				fmt.Printf("Ignored: %s\n", path)
				return nil
			}
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil
		}

		str := string(data)

		lines := strings.Split(str, "\n")
		for _, line := range lines {
			total++
			if lineIsEmpty(line) {
				continue
			}

			notEmpty++

			trimmed := strings.Trim(line, " ")

			if strings.HasPrefix(trimmed, "//") {
				comments++
			} else {
				code++
			}
		}

		return nil
	})

	return LinesInfo{
		Total:    total,
		NotEmpty: notEmpty,
		Comments: comments,
		Code:     code,
	}
}

func lineIsEmpty(line string) bool {
	if len(line) == 0 {
		return true
	}

	for _, c := range line {
		if c != ' ' {
			return false
		}
	}

	return true
}