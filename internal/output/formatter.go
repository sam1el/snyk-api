// Package output provides output formatting for CLI commands.
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/itchyny/gojq"
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
	format       Format
	writer       io.Writer
	templateText string
	jqQuery      string
}

var (
	defaultTemplate string
	defaultJQ       string
)

// New creates a new formatter with the specified format.
func New(format string) *Formatter {
	return &Formatter{
		format:       Format(format),
		writer:       os.Stdout,
		templateText: defaultTemplate,
		jqQuery:      defaultJQ,
	}
}

// WithWriter sets a custom writer (useful for testing).
func (f *Formatter) WithWriter(w io.Writer) *Formatter {
	f.writer = w
	return f
}

// WithTemplate sets a text/template string to render output.
func (f *Formatter) WithTemplate(tmpl string) *Formatter {
	f.templateText = tmpl
	return f
}

// WithJQ sets a jq-style query to pre-filter data.
func (f *Formatter) WithJQ(query string) *Formatter {
	f.jqQuery = query
	return f
}

// SetDefaultTemplate sets a process-wide default template.
func SetDefaultTemplate(tmpl string) {
	defaultTemplate = tmpl
}

// SetDefaultJQ sets a process-wide default jq query.
func SetDefaultJQ(query string) {
	defaultJQ = query
}

// Print outputs data in the configured format.
func (f *Formatter) Print(data interface{}) error {
	// Apply jq filter first if provided.
	if f.jqQuery != "" {
		filtered, err := f.applyJQ(data)
		if err != nil {
			return err
		}
		data = filtered
	}

	// Template output takes precedence over format selection.
	if f.templateText != "" {
		return f.printTemplate(data)
	}

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

func (f *Formatter) printTemplate(data interface{}) error {
	tmpl, err := template.New("out").Option("missingkey=zero").Parse(f.templateText)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	_, err = f.writer.Write(buf.Bytes())
	return err
}

func (f *Formatter) applyJQ(data interface{}) (interface{}, error) {
	query, err := gojq.Parse(f.jqQuery)
	if err != nil {
		return nil, fmt.Errorf("parse jq query: %w", err)
	}

	normalized, err := normalizeForJQ(data)
	if err != nil {
		return nil, err
	}

	iter := query.Run(normalized)
	v, ok := iter.Next()
	if !ok {
		return nil, nil
	}
	if err, ok := v.(error); ok {
		return nil, fmt.Errorf("apply jq query: %w", err)
	}
	return v, nil
}

func normalizeForJQ(data interface{}) (interface{}, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal for jq: %w", err)
	}
	var out interface{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("unmarshal for jq: %w", err)
	}
	return out, nil
}
