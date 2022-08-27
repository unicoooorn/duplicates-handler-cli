package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

type dfh struct{}

func (dfh) printFilesBySize(files map[int64][]string, isAscending bool) {
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

func (dfh) printFilesByHash(files map[int64]map[string][]string, isAscending bool) {
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

func (dfh) askForSortOpt() bool {
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

func (dfh) askForHashOpt() bool {
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

func (dfh) askForExt() string {
	var ext string
	fmt.Println("Enter file format:")
	if bytesCount, err := fmt.Scanln(&ext); err != nil && bytesCount > 0 {
		log.Fatal(err)
	}
	fmt.Println()
	return ext
}

func (dfh) printSortingMenu() {
	fmt.Println("Size sorting options:")
	fmt.Println("1. Descending")
	fmt.Println("2. Ascending")
	fmt.Println()
}

func (cur dfh) Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Directory is not specified")
		return
	}
	rootPath := os.Args[1]

	ext := cur.askForExt()
	cur.printSortingMenu()
	isAscending := cur.askForSortOpt()
	filesOfSize := cur.findFilesBySizeExt(rootPath, ext)
	cur.printFilesBySize(filesOfSize, isAscending)
	if cur.askForHashOpt() {
		filesOfHash := cur.findFilesByHashExt(rootPath, ext)
		cur.printFilesByHash(filesOfHash, isAscending)
	}

}
