package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"unicode"
)

func main() {
	dir := flag.String("dir", "", "Cache source directory.")

	flag.Parse()

	if *dir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := run(*dir); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func run(dir string) error {
	outFile, err := os.Create("output.csv")
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := csv.NewWriter(outFile)
	defer w.Flush()

	if err := w.Write([]string{"file", "content"}); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := path.Join(dir, entry.Name())

		valid, err := isJSON(filePath)
		if err != nil {
			return err
		}

		if !valid {
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		if err := w.Write([]string{entry.Name(), string(content)}); err != nil {
			return err
		}
	}

	return nil
}

func isJSON(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	// Skip UTF-8 BOM: EF BB BF
	b, err := r.Peek(3)
	if err == nil && b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF {
		_, _ = r.Discard(3)
	}

	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}

		if !unicode.IsSpace(ch) {
			return ch == '{', nil
		}
	}
}
