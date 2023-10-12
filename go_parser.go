package languages

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"

	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type GoParser struct {
	fileList   []*object.File
	annotation Annotation
}

func NewGoparser(fileList []*object.File, language *Language) *GoParser {
	return &GoParser{
		fileList:   fileList,
		annotation: language.Annotation,
	}
}

func (g *GoParser) Parse(fileList []*object.File) error {
	for _, file := range g.fileList {

		content, err := file.Contents()
		if err != nil {
			return err
		}

		lineCount, err := g.scan(content)
		if err != nil {
			return err
		}
		fmt.Printf("File: %s, Line count: %d\n", file.Name, lineCount)

		// Create a new CSV writer.
		f, err := os.Create("count.csv")
		defer f.Close()
		if err != nil {
			log.Fatalln("Failed to open file ", err)
		}
		w := csv.NewWriter(f)
		defer w.Flush()

		// Write the header row.
		if err := w.Write([]string{"File name", "Line count"}); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}

		// Write the data to the CSV file.
		for _, file2 := range g.fileList {

			content2, err := file2.Contents()
			if err != nil {
				return err
			}
			lineCount2, err := g.scan(content2)
			if err != nil {
				return err
			}

			if err := w.Write([]string{file2.Name, fmt.Sprintf("%d", lineCount2)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

	}

	return nil
}

func (g *GoParser) scan(content string) (int, error) {

	// Create a scanner to read the file line by line.
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Initialize variables to keep track of the start and end annotations
	startAnnotationFound := false
	endAnnotationFound := false

	// Initialize a variable to keep track of the line count
	lineCount := 0

	// Loop through each line in the file
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the start annotation
		if !startAnnotationFound && line == g.annotation.Start {
			startAnnotationFound = true
			continue
		}

		// Check if the line contains the end annotation
		if startAnnotationFound && line == g.annotation.End {
			endAnnotationFound = true
			break
		}

		// If we're between the start and end annotations, increment the line count
		if startAnnotationFound && !endAnnotationFound {
			lineCount++
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("error: %v", err)
		return 0, err
	}

	// Return the line count
	return lineCount, nil
}

// +copilot
// func (g *GoParser) Parse(fileList []*object.File) error {
// 	fmt.Println(g.fileList)
// 	for _, file := range g.fileList {

// 		content, err := file.Contents()
// 		// fmt.Println(file.Name + "#####" + content)
// 		if err != nil {
// 			return err
// 		}

// 		lineCount, err := g.scan(content)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Printf("File: %s, Line count: %d\n", file.Name, lineCount)

// 	}
// 	return nil
// }

// -copilot
