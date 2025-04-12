package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	// Parse command-line arguments
	sourceDir := flag.String("source", "", "Source directory containing photos")
	destDir := flag.String("dest", "", "Destination directory for sorted photos")
	moveFiles := flag.Bool("move", false, "Move files instead of copying them")
	fileFormat := flag.String("format", "", "Specific file format to process (e.g., 'jpg,png'). Leave empty for all supported formats")
	flag.Parse()

	// Validate command-line arguments
	if *sourceDir == "" || *destDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Ensure the source directory exists
	sourceStat, err := os.Stat(*sourceDir)
	if err != nil || !sourceStat.IsDir() {
		log.Fatalf("Source directory does not exist or is not a directory: %v", *sourceDir)
	}

	// Ensure the destination directory exists, create if not
	if err := os.MkdirAll(*destDir, 0755); err != nil {
		log.Fatalf("Failed to create destination directory: %v", err)
	}

	// Process the file format parameter
	var formats []string
	if *fileFormat != "" {
		// Split the format string by comma and trim spaces
		for _, f := range strings.Split(*fileFormat, ",") {
			format := strings.TrimSpace(f)
			if format != "" {
				// Add dot prefix if not present
				if !strings.HasPrefix(format, ".") {
					format = "." + format
				}
				formats = append(formats, strings.ToLower(format))
			}
		}
	}

	// Walk through the source directory
	err = filepath.Walk(*sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		
		// Check if the file is an image and matches the format filter (if any)
		if !isValidFileFormat(ext, formats) {
			return nil
		}

		// Get date from EXIF data
		date, err := getPhotoDate(path)
		if err != nil {
			log.Printf("Warning: Could not get date for %s: %v", path, err)
			return nil
		}

		// Create destination directory structure: yyyy/mm/
		yearMonth := filepath.Join(*destDir, fmt.Sprintf("%04d", date.Year()), fmt.Sprintf("%02d", date.Month()))
		if err := os.MkdirAll(yearMonth, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", yearMonth, err)
		}

		// Destination file path
		destPath := filepath.Join(yearMonth, filepath.Base(path))

		// Copy or move the file
		if *moveFiles {
			if err := moveFile(path, destPath); err != nil {
				return fmt.Errorf("failed to move %s to %s: %v", path, destPath, err)
			}
			log.Printf("Moved %s to %s", path, destPath)
		} else {
			if err := copyFile(path, destPath); err != nil {
				return fmt.Errorf("failed to copy %s to %s: %v", path, destPath, err)
			}
			log.Printf("Copied %s to %s", path, destPath)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error processing files: %v", err)
	}

	log.Println("Photo sorting completed successfully!")
}

// isValidFileFormat checks if the file extension is valid based on the format filter
func isValidFileFormat(ext string, formats []string) bool {
	// If no specific formats are specified, check against all supported formats
	if len(formats) == 0 {
		return isImageFile(ext)
	}
	
	// Otherwise, check if the extension is in the list of specified formats
	for _, format := range formats {
		if ext == format {
			return true
		}
	}
	
	return false
}

// isImageFile returns true if the file extension corresponds to a common image format
func isImageFile(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".heic", ".heif", ".raw", ".cr2", ".nef":
		return true
	default:
		return false
	}
}

// getPhotoDate extracts the date when the photo was taken from EXIF metadata
func getPhotoDate(filepath string) (time.Time, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return time.Time{}, err
	}
	defer file.Close()

	// Decode EXIF data
	x, err := exif.Decode(file)
	if err != nil {
		return time.Time{}, err
	}

	// Try to get the date the photo was taken
	datetime, err := x.DateTime()
	if err != nil {
		return time.Time{}, err
	}

	return datetime, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Check if destination file already exists
	if _, err := os.Stat(dst); err == nil {
		// File exists, don't overwrite
		log.Printf("Skipping %s: file already exists at destination", dst)
		return nil
	}

	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write to destination
	return os.WriteFile(dst, data, 0644)
}

// moveFile moves a file from src to dst
func moveFile(src, dst string) error {
	// Check if destination file already exists
	if _, err := os.Stat(dst); err == nil {
		// File exists, don't overwrite
		log.Printf("Skipping %s: file already exists at destination", dst)
		return nil
	}

	// Use os.Rename to move the file
	return os.Rename(src, dst)
} 