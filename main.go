package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
)

type FilesSummary struct {
	Files       int
	Dirs        int
	FileSizeSum int64
}

func main() {
	flag.Parse()
	paths := flag.Args()

	if len(paths) <= 0 {
		fmt.Println("usage: rr file ...")
		os.Exit(1)
	}

	if err := execute(paths); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}
}

func execute(paths []string) error {
	summary, err := calcFilesSummary(paths)
	if err != nil {
		return err
	}

	fmt.Printf("%d directories, %d files (%s)\n",
		summary.Dirs, summary.Files, humanize.Bytes(uint64(summary.FileSizeSum)))
	removeOk, err := confirm(fmt.Sprintf("remove %s?", strings.Join(paths, ", ")))
	if err != nil {
		return err
	}
	if removeOk {
		if err := removeAll(paths); err != nil {
			return err
		}
	}
	return nil
}

func removeAll(paths []string) error {
	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func confirm(message string) (bool, error) {
	var input string

	fmt.Printf("%s (y/N): ", message)
	if _, err := fmt.Scan(&input); err != nil {
		return false, err
	}
	input = strings.TrimSpace(input)
	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true, nil
	} else {
		return false, nil
	}
}

func calcFilesSummary(paths []string) (*FilesSummary, error) {
	summary := &FilesSummary{}

	for _, path := range paths {
		err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				summary.Dirs++
			} else {
				summary.Files++
				summary.FileSizeSum += info.Size()
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return summary, nil
}
