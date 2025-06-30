package csv2sql

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
)

// Convert reads CSV data from the provided reader, creates a SQL table definition,
// and generates corresponding INSERT statements, writing the output to the given writer.
// The table name is used for both the CREATE TABLE and INSERT statements.
// Returns an error if any I/O operation fails or if the CSV data is malformed.
func Convert(tableName string, reader io.Reader, writer io.Writer) error {
	r := csv.NewReader(reader)
	r.ReuseRecord = true

	headers, err := r.Read()
	if err != nil {
		return err
	}

	values := make([]string, 0, len(headers))
	for i, header := range headers {
		headers[i] = formatField(header)
	}

	for _, header := range headers {
		values = append(values, "    "+header+" TEXT")
	}

	_, err = writer.Write([]byte("CREATE TABLE " + formatField(tableName) + " (\n" + strings.Join(values, ",\n") + "\n);\n"))
	if err != nil {
		return err
	}

	firstRecord := true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Write INSERT statement
		if firstRecord {
			_, err = writer.Write([]byte("INSERT INTO " + formatField(tableName) + " VALUES\n"))
			if err != nil {
				return err
			}
			firstRecord = false
		} else {
			_, err = writer.Write([]byte(",\n"))
			if err != nil {
				return err
			}
		}

		if len(record) > len(headers) {
			record = record[:len(headers)]
		} else {
			for len(record) < len(headers) {
				record = append(record, "")
			}
		}
		for i, v := range record {
			values[i] = strconv.Quote(v)
		}
		_, err = writer.Write([]byte("(" + strings.Join(values, ", ") + ")"))
		if err != nil {
			return err
		}
	}

	_, err = writer.Write([]byte(";\n"))
	if err != nil {
		return err
	}

	return nil
}

func formatField(value string) string {
	return "`" + value + "`"
}
