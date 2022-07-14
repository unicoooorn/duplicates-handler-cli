package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func findFilesBySizeExt(rootPath string, ext string) map[int64][]string {
	filesOfSize := make(map[int64][]string)
	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if len(ext) != 0 && filepath.Ext(path) != ext {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		filesOfSize[info.Size()] = append(filesOfSize[info.Size()], path)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return filesOfSize
}

func printFilesBySize(files map[int64][]string, opt int) {
	var sizes []int64
	for size := range files {
		sizes = append(sizes, size)
	}
	sort.Slice(sizes, func(i, j int) bool {
		switch opt {
		case 1:
			return sizes[i] > sizes[j]

		case 2:
			return sizes[i] < sizes[j]
		}
		log.Fatal("There is no such option")
		return true
	})

	for _, size := range sizes {
		fmt.Println(size, "bytes")
		for _, filePath := range files[size] {
			fmt.Println(filePath)
		}
		fmt.Println()
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	}
	rootPath := os.Args[1]

	// reading extension to look for (ext is empty = looking for any extension)
	var ext string
	fmt.Println("Enter file format:")
	if _, err := fmt.Scanln(&ext); err != nil {
		log.Fatal(err)
	}
	fmt.Println()

	// traversing directory
	filesOfSize := findFilesBySizeExt(rootPath, ext)
	// sorting menu
	fmt.Println("Size sorting options:")
	fmt.Println("1. Descending")
	fmt.Println("2. Ascending")
	fmt.Println()
	// scanning for sorting options
	isCorrect := false
	var opt int
	for !isCorrect {
		fmt.Println("Enter a sorting option:")
		if _, err := fmt.Scan(&opt); err != nil {
			log.Fatal(err)
		}
		fmt.Println()

		switch opt {
		case 1:
			isCorrect = true
			break
		case 2:
			isCorrect = true
		default:
			fmt.Println("Wrong option")
			fmt.Println()
		}
	}
	printFilesBySize(filesOfSize, opt)
}
