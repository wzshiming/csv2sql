package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wzshiming/csv2sql"
)

const defaultTableName = "csv2sql"

var (
	tableName  string
	inputFile  string
	outputFile string
)

func init() {
	flag.StringVar(&tableName, "table", defaultTableName, "Name of the table to create")
	flag.StringVar(&inputFile, "input", "", "Input CSV file (default: stdin)")
	flag.StringVar(&outputFile, "output", "", "Output SQL file (default: stdout)")
	flag.Parse()
}

func main() {
	if tableName == "" {
		fmt.Println("Error: table name cannot be empty")
		flag.Usage()
		os.Exit(1)
	}

	// Set up input
	var input io.Reader = os.Stdin
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	}

	// Set up output
	var output io.Writer = os.Stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		output = file
	}

	err := csv2sql.Convert(tableName, input, output)
	if err != nil {
		fmt.Printf("Error converting CSV: %v\n", err)
		os.Exit(1)
	}
}
