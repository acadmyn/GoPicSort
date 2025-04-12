# GoPicSort

GoPicSort is a Go application that organizes your photos into folders based on when they were taken. It reads EXIF metadata from photos and sorts them into a year/month folder structure.

## Features

- Sorts photos into folders based on year and month (e.g., `2023/05/` for photos taken in May 2023)
- Supports common image formats (JPG, JPEG, PNG, TIFF, HEIC, RAW, etc.)
- Option to copy or move files
- Filter by specific file formats
- Skips files that already exist in the destination

## Installation

```bash
# Clone the repository
git clone https://github.com/acadmyn/GoPicSort.git
cd gopicsort

# Build the application
go build
```

## Usage

```bash
# Copy photos to a sorted structure
./gopicsort -source /path/to/photos -dest /path/to/sorted/photos

# Move photos instead of copying them
./gopicsort -source /path/to/photos -dest /path/to/sorted/photos -move

# Process only JPG and PNG files
./gopicsort -source /path/to/photos -dest /path/to/sorted/photos -format "jpg,png"
```

### Command-line Options

- `-source`: Source directory containing photos (required)
- `-dest`: Destination directory for sorted photos (required)
- `-move`: Move files instead of copying them (optional, default is to copy)
- `-format`: Specific file format(s) to process, comma-separated (e.g., "jpg,png,heic"). Leave empty to process all supported formats.

## How It Works

1. The application walks through all files in the source directory
2. For each image file (filtered by format if specified), it extracts the date taken from EXIF metadata
3. It creates a directory structure based on year and month (YYYY/MM)
4. It copies or moves the file to the appropriate directory

## Requirements

- Go 1.18 or higher
- The `github.com/rwcarlsen/goexif` package for EXIF metadata extraction 
