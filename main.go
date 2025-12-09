package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dir := flag.String("dir", ".", "Path to folder with your photos")
	flag.Parse()

	fmt.Println("Scanning folder:", *dir)

	var totalImages int
	var jpgCount, heicCount int
	var totalBytes int64

	err := filepath.WalkDir(*dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		isImage := false

		switch ext {
		case ".jpg", ".JPG", ".jpeg", ".JPEG":
			jpgCount++
			isImage = true
		case ".heic", ".HEIC":
			heicCount++
			isImage = true
		}

		if !isImage {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			return nil
		}

		size := info.Size()
		totalBytes += size
		totalImages++

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("===== SUMMARY =====")
	fmt.Printf("Total images: %d\n", totalImages)
	fmt.Printf("  JPG / JPEG: %d\n", jpgCount)
	fmt.Printf("  HEIC:       %d\n", heicCount)
	fmt.Printf("Total size: %d bytes\n", totalBytes)
}
