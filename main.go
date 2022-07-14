package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
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

func findFilesByHashExt(rootPath string, ext string) map[int64]map[string][]string {
	filesByHash := make(map[int64]map[string][]string)
	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if len(ext) != 0 && filepath.Ext(path) != ext {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			log.Fatal(err)
		}
		if _, ok := filesByHash[info.Size()]; !ok {
			filesByHash[info.Size()] = make(map[string][]string)
		}
		hashStr := hex.EncodeToString(h.Sum(nil))
		filesByHash[info.Size()][hashStr] = append(filesByHash[info.Size()][hashStr], path)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return filesByHash
}

func printFilesBySize(files map[int64][]string, isAscending bool) {
	var sizes []int64
	for size := range files {
		sizes = append(sizes, size)
	}
	sort.Slice(sizes, func(i, j int) bool {
		if isAscending {
			return sizes[i] < sizes[j]
		} else {
			return sizes[i] > sizes[j]
		}
	})

	for _, size := range sizes {
		fmt.Println(size, "bytes")
		for _, filePath := range files[size] {
			fmt.Println(filePath)
		}
		fmt.Println()
	}
}

func printFilesByHash(files map[int64]map[string][]string, isAscending bool) {
	var sizes []int64
	for size := range files {
		sizes = append(sizes, size)
	}
	sort.Slice(sizes, func(i, j int) bool {
		if isAscending {
			return sizes[i] < sizes[j]
		} else {
			return sizes[i] > sizes[j]
		}
	})

	counter := 1
	for _, size := range sizes {
		fmt.Println(size, "bytes")
		for h := range files[size] {
			if len(files[size][h]) == 1 {
				continue
			}
			fmt.Println("Hash:", h)
			for _, filePath := range files[size][h] {
				fmt.Printf("%d. %s\n", counter, filePath)
				counter++
			}
			fmt.Println()
		}

		fmt.Println()
	}
}

func askForSortOpt() bool {
	opt := 3
	for opt != 1 && opt != 2 {
		fmt.Println("Enter a sorting option:")
		if _, err := fmt.Scan(&opt); err != nil {
			log.Fatal(err)
		}
		fmt.Println()

		switch opt {
		case 1:
			return false
		case 2:
			return true
		default:
			fmt.Println("Wrong option")
			fmt.Println()
		}
	}
	return true
}

func askForHashOpt() bool {
	opt := ""
	for opt != "yes" && opt != "no" {
		fmt.Println("Check for duplicates?")
		if _, err := fmt.Scan(&opt); err != nil {
			log.Fatal(err)
		}
		fmt.Println()

		switch opt {
		case "yes":
			return true
		case "no":
			return false
		default:
			fmt.Println("Wrong option")
			fmt.Println()
		}
	}
	return true
}

func askForExt() string {
	var ext string
	fmt.Println("Enter file format:")
	fmt.Scanln(&ext)
	fmt.Println()
	return ext
}

func printSortingMenu() {
	fmt.Println("Size sorting options:")
	fmt.Println("1. Descending")
	fmt.Println("2. Ascending")
	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	}
	rootPath := os.Args[1]

	ext := askForExt()
	printSortingMenu()
	isAscending := askForSortOpt()
	filesOfSize := findFilesBySizeExt(rootPath, ext)
	printFilesBySize(filesOfSize, isAscending)
	if askForHashOpt() {
		filesOfHash := findFilesByHashExt(rootPath, ext)
		printFilesByHash(filesOfHash, isAscending)
	}

}
