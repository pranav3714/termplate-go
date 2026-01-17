package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/blacksilver/ever-so-powerful/internal/config"
)

// Formatter handles formatting output in different formats
type Formatter struct {
	config config.OutputConfig
	writer io.Writer
}

// NewFormatter creates a new output formatter
func NewFormatter(cfg config.OutputConfig) *Formatter {
	return &Formatter{
		config: cfg,
		writer: os.Stdout,
	}
}

// NewFormatterWithWriter creates a formatter with a custom writer
func NewFormatterWithWriter(cfg config.OutputConfig, w io.Writer) *Formatter {
	return &Formatter{
		config: cfg,
		writer: w,
	}
}

// Print formats and prints data based on the configured output format
func (f *Formatter) Print(data interface{}) error {
	switch f.config.Format {
	case "json":
		return f.printJSON(data)
	case "yaml":
		return f.printYAML(data)
	case "table":
		return f.printTable(data)
	case "csv":
		return f.printCSV(data)
	default:
		return f.printText(data)
	}
}

// printJSON outputs data as JSON
func (f *Formatter) printJSON(data interface{}) error {
	var output []byte
	var err error

	if f.config.Pretty {
		output, err = json.MarshalIndent(data, "", "  ")
	} else {
		output, err = json.Marshal(data)
	}

	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if _, err = fmt.Fprintln(f.writer, string(output)); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}

// printYAML outputs data as YAML
func (f *Formatter) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(f.writer)
	if f.config.Pretty {
		encoder.SetIndent(2)
	}
	defer encoder.Close()

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encoding YAML: %w", err)
	}
	return nil
}

// printTable outputs data as a table
func (f *Formatter) printTable(data interface{}) error {
	// Convert data to table format
	table, err := f.toTable(data)
	if err != nil {
		return err
	}

	// Print table based on style
	switch f.config.TableStyle {
	case "unicode":
		f.printUnicodeTable(table)
	case "markdown":
		f.printMarkdownTable(table)
	default:
		f.printASCIITable(table)
	}

	return nil
}

// printCSV outputs data as CSV
func (f *Formatter) printCSV(data interface{}) error {
	table, err := f.toTable(data)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(f.writer)
	defer writer.Flush()

	for _, row := range table {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("writing CSV: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}
	return nil
}

// printText outputs data as plain text
func (f *Formatter) printText(data interface{}) error {
	if _, err := fmt.Fprintln(f.writer, data); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}

// toTable converts various data types to table format
func (f *Formatter) toTable(data interface{}) ([][]string, error) {
	switch v := data.(type) {
	case [][]string:
		return v, nil
	case []map[string]string:
		return f.mapSliceToTable(v), nil
	case map[string]string:
		return f.mapToTable(v), nil
	default:
		return nil, fmt.Errorf("unsupported data type for table output")
	}
}

// mapSliceToTable converts a slice of maps to table format
func (f *Formatter) mapSliceToTable(data []map[string]string) [][]string {
	if len(data) == 0 {
		return [][]string{}
	}

	// Get headers from first map
	headers := make([]string, 0, len(data[0]))
	for k := range data[0] {
		headers = append(headers, k)
	}

	// Build table
	table := [][]string{headers}
	for _, row := range data {
		var rowData []string
		for _, h := range headers {
			rowData = append(rowData, row[h])
		}
		table = append(table, rowData)
	}

	return table
}

// mapToTable converts a single map to table format
func (f *Formatter) mapToTable(data map[string]string) [][]string {
	table := [][]string{{"Key", "Value"}}
	for k, v := range data {
		table = append(table, []string{k, v})
	}
	return table
}

// printASCIITable prints a table using ASCII characters
func (f *Formatter) printASCIITable(table [][]string) {
	if len(table) == 0 {
		return
	}

	// Calculate column widths
	widths := f.calculateColumnWidths(table)

	// Print header
	f.printASCIIRow(table[0], widths, true)

	// Print separator
	f.printASCIISeparator(widths)

	// Print rows
	for _, row := range table[1:] {
		f.printASCIIRow(row, widths, false)
	}
}

// printUnicodeTable prints a table using Unicode box drawing characters
func (f *Formatter) printUnicodeTable(table [][]string) {
	if len(table) == 0 {
		return
	}

	widths := f.calculateColumnWidths(table)

	// Print top border
	f.printUnicodeBorder(widths, "┌", "┬", "┐")

	// Print header
	f.printUnicodeRow(table[0], widths)

	// Print header separator
	f.printUnicodeBorder(widths, "├", "┼", "┤")

	// Print rows
	for _, row := range table[1:] {
		f.printUnicodeRow(row, widths)
	}

	// Print bottom border
	f.printUnicodeBorder(widths, "└", "┴", "┘")
}

// printMarkdownTable prints a table in Markdown format
func (f *Formatter) printMarkdownTable(table [][]string) {
	if len(table) == 0 {
		return
	}

	widths := f.calculateColumnWidths(table)

	// Print header
	f.printMarkdownRow(table[0], widths)

	// Print separator
	fmt.Fprint(f.writer, "|")
	for _, w := range widths {
		fmt.Fprint(f.writer, strings.Repeat("-", w+2), "|")
	}
	fmt.Fprintln(f.writer)

	// Print rows
	for _, row := range table[1:] {
		f.printMarkdownRow(row, widths)
	}
}

// calculateColumnWidths calculates the width of each column
func (f *Formatter) calculateColumnWidths(table [][]string) []int {
	if len(table) == 0 {
		return nil
	}

	widths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	return widths
}

// Helper functions for ASCII table
func (f *Formatter) printASCIIRow(row []string, widths []int, _ bool) {
	fmt.Fprint(f.writer, "| ")
	for i, cell := range row {
		fmt.Fprintf(f.writer, "%-*s | ", widths[i], cell)
	}
	fmt.Fprintln(f.writer)
}

func (f *Formatter) printASCIISeparator(widths []int) {
	fmt.Fprint(f.writer, "|")
	for _, w := range widths {
		fmt.Fprint(f.writer, strings.Repeat("-", w+2), "|")
	}
	fmt.Fprintln(f.writer)
}

// Helper functions for Unicode table
func (f *Formatter) printUnicodeRow(row []string, widths []int) {
	fmt.Fprint(f.writer, "│ ")
	for i, cell := range row {
		fmt.Fprintf(f.writer, "%-*s │ ", widths[i], cell)
	}
	fmt.Fprintln(f.writer)
}

func (f *Formatter) printUnicodeBorder(widths []int, left, mid, right string) {
	fmt.Fprint(f.writer, left)
	for i, w := range widths {
		fmt.Fprint(f.writer, strings.Repeat("─", w+2))
		if i < len(widths)-1 {
			fmt.Fprint(f.writer, mid)
		}
	}
	fmt.Fprintln(f.writer, right)
}

// Helper functions for Markdown table
func (f *Formatter) printMarkdownRow(row []string, widths []int) {
	fmt.Fprint(f.writer, "| ")
	for i, cell := range row {
		fmt.Fprintf(f.writer, "%-*s | ", widths[i], cell)
	}
	fmt.Fprintln(f.writer)
}
