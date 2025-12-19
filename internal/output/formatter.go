// Package output provides output formatting for CLI commands.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format represents an output format.
type Format string

const (
	// FormatJSON outputs data as JSON.
	FormatJSON Format = "json"
	// FormatYAML outputs data as YAML.
	FormatYAML Format = "yaml"
	// FormatTable outputs data as an ASCII table.
	FormatTable Format = "table"
)

// Formatter handles output formatting.
type Formatter struct {
	format Format
	writer io.Writer
}

// New creates a new formatter with the specified format.
func New(format string) *Formatter {
	return &Formatter{
		format: Format(format),
		writer: os.Stdout,
	}
}

// WithWriter sets a custom writer (useful for testing).
func (f *Formatter) WithWriter(w io.Writer) *Formatter {
	f.writer = w
	return f
}

// Print outputs data in the configured format.
func (f *Formatter) Print(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.printJSON(data)
	case FormatYAML:
		return f.printYAML(data)
	case FormatTable:
		return f.printTable(data)
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

// printJSON outputs data as formatted JSON.
func (f *Formatter) printJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printYAML outputs data as YAML.
func (f *Formatter) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(f.writer)
	encoder.SetIndent(2)
	defer encoder.Close() //nolint:errcheck // Best effort cleanup
	return encoder.Encode(data)
}

// printTable outputs data as an ASCII table using tabwriter.
func (f *Formatter) printTable(data interface{}) error {
	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush() //nolint:errcheck // Best effort cleanup

	// Handle different data types
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return f.printSliceTable(w, v)
	case reflect.Struct:
		return f.printStructTable(w, v)
	case reflect.Map:
		return f.printMapTable(w, v)
	default:
		// Fallback to JSON for unsupported types
		return f.printJSON(data)
	}
}

// printSliceTable prints a slice as a table.
func (f *Formatter) printSliceTable(w *tabwriter.Writer, v reflect.Value) error {
	if v.Len() == 0 {
		_, _ = fmt.Fprintln(w, "(no results)") //nolint:errcheck // Output formatting
		return nil
	}

	// Get first element to determine structure
	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}

	if first.Kind() == reflect.Struct {
		// Extract headers from struct fields
		headers := []string{}
		typ := first.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.IsExported() {
				headers = append(headers, strings.ToUpper(field.Name))
			}
		}
		_, _ = fmt.Fprintln(w, strings.Join(headers, "\t")) //nolint:errcheck // Output formatting

		// Add rows
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			row := []string{}
			for j := 0; j < elem.NumField(); j++ {
				if typ.Field(j).IsExported() {
					row = append(row, fmt.Sprintf("%v", elem.Field(j).Interface()))
				}
			}
			_, _ = fmt.Fprintln(w, strings.Join(row, "\t")) //nolint:errcheck // Output formatting
		}
	}

	return nil
}

// printStructTable prints a struct as a key-value table.
func (f *Formatter) printStructTable(w *tabwriter.Writer, v reflect.Value) error {
	_, _ = fmt.Fprintln(w, "FIELD\tVALUE") //nolint:errcheck // Output formatting

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := typ.Field(i)
		if field.IsExported() {
			_, _ = fmt.Fprintf(w, "%s\t%v\n", field.Name, v.Field(i).Interface()) //nolint:errcheck // Output formatting
		}
	}

	return nil
}

// printMapTable prints a map as a key-value table.
func (f *Formatter) printMapTable(w *tabwriter.Writer, v reflect.Value) error {
	_, _ = fmt.Fprintln(w, "KEY\tVALUE") //nolint:errcheck // Output formatting

	for _, key := range v.MapKeys() {
		_, _ = fmt.Fprintf(w, "%v\t%v\n", key.Interface(), v.MapIndex(key).Interface()) //nolint:errcheck // Output formatting
	}

	return nil
}
