package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Metadata struct {
	Title          string   `json:"title"`
	PhotoTakenTime TimeInfo `json:"photoTakenTime"`
	GeoDataExif    GeoData  `json:"geoDataExif"`
}

type TimeInfo struct {
	Timestamp string `json:"timestamp"`
	Formatted string `json:"formatted"`
}

type GeoData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

func main() {
	dirPath := flag.String("dir", ".", "Directory containing photos and JSON metadata files")
	dryRun := flag.Bool("dry-run", false, "Preview changes without modifying files")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	log.Printf("Processing directory: %s", *dirPath)
	if *dryRun {
		log.Println("DRY RUN MODE - No files will be modified")
	}

	err := processDirectory(*dirPath, *dryRun, *verbose)

	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	log.Println("✓ Processing complete!")
}

func processDirectory(dir string, dryRun bool, verbose bool) error {
	fileCount := 0
	processedCount := 0
	errorCount := 0

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Warning: Error accessing path %s: %v", path, err)
			errorCount++
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(strings.ToLower(path), ".json") {
			fileCount++
			if err := processMetadataFile(path, dryRun, verbose); err != nil {
				log.Printf("Error processing %s: %v", filepath.Base(path), err)
				errorCount++
			} else {
				processedCount++
			}
		}

		return nil
	})

	log.Printf("Summary: %d JSON files found, %d processed, %d errors", fileCount, processedCount, errorCount)

	return err
}

func processMetadataFile(jsonPath string, dryRun bool, verbose bool) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return err
	}

	mediaPath := findMediaFile(jsonPath, metadata.Title)
	if mediaPath == "" {
		if verbose {
			log.Printf("⚠️  Media file not found for: %s", metadata.Title)
		}
		return nil
	}

	log.Printf("✓ Found: %s", filepath.Base(mediaPath))

	if verbose {
		log.Printf(" Date: %s", metadata.PhotoTakenTime.Formatted)
		log.Printf(" GPS: %.6f, %.6f", metadata.GeoDataExif.Latitude, metadata.GeoDataExif.Longitude)
	}

	if !dryRun {
		if err := writeMetadataToFile(mediaPath, metadata); err != nil {
			return err
		}
		if verbose {
			log.Printf("  ✓ Metadata written successfully")
		}
	}

	return nil
}

func findMediaFile(jsonPath, title string) string {
	baseDir := filepath.Dir(jsonPath)
	mediaPath := filepath.Join(baseDir, title)

	if _, err := os.Stat(mediaPath); err == nil {
		return mediaPath
	}

	return ""
}

func writeMetadataToFile(filePath string, metadata Metadata) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".jpg", ".jpeg":
		return writeJPEGMetadata(filePath, metadata)
	case ".heic":
		return writeHEICMetadata(filePath, metadata)
	case ".mov", ".mp4":
		return writeVideoMetadata(filePath, metadata)
	default:
		return nil
	}
}

func writeJPEGMetadata(filePath string, metadata Metadata) error {
	// TODO: Implement JPEG EXIF writing
	return nil
}

func writeHEICMetadata(filePath string, metadata Metadata) error {
	// TODO: Implement HEIC metadata writing
	return nil
}

func writeVideoMetadata(filePath string, metadata Metadata) error {
	// TODO: Implement video metadata writing
	return nil
}
