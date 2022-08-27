package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (dfh) findFilesBySizeExt(rootPath string, ext string) map[int64][]string {
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

func (dfh) findFilesByHashExt(rootPath string, ext string) map[int64]map[string][]string {
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
		defer func() {
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
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
