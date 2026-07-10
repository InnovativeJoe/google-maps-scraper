package web

import (
	"context"
	"encoding/csv"
	"fmt"
	"reflect"
	"sync"

	"github.com/gosom/scrapemate"
)

var _ scrapemate.ResultWriter = (*filteredCSVWriter)(nil)

type filteredCSVWriter struct {
	w         *csv.Writer
	once      sync.Once
	selection FieldSelection
}

// NewFilteredCsvWriter creates a CSV writer that omits unchecked optional
// fields from the output headers and rows.
func NewFilteredCsvWriter(w *csv.Writer, selection FieldSelection) scrapemate.ResultWriter {
	return &filteredCSVWriter{w: w, selection: selection}
}

func (c *filteredCSVWriter) Run(_ context.Context, in <-chan scrapemate.Result) error {
	for result := range in {
		elements, err := c.getCsvCapable(result.Data)
		if err != nil {
			return err
		}

		if len(elements) == 0 {
			continue
		}

		c.once.Do(func() {
			_ = c.w.Write(c.filterHeaders(elements[0].CsvHeaders()))
		})

		for _, element := range elements {
			if err := c.w.Write(c.filterRow(element.CsvHeaders(), element.CsvRow())); err != nil {
				return err
			}
		}

		c.w.Flush()
	}

	return c.w.Error()
}

func (c *filteredCSVWriter) getCsvCapable(data any) ([]scrapemate.CsvCapable, error) {
	var elements []scrapemate.CsvCapable

	if isSlice(data) {
		s := reflect.ValueOf(data)

		for i := 0; i < s.Len(); i++ {
			val := s.Index(i).Interface()
			element, ok := val.(scrapemate.CsvCapable)
			if !ok {
				return nil, fmt.Errorf("%w: unexpected data type: %T", scrapemate.ErrorNotCsvCapable, val)
			}

			elements = append(elements, element)
		}
		return elements, nil
	}

	element, ok := data.(scrapemate.CsvCapable)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected data type: %T", scrapemate.ErrorNotCsvCapable, data)
	}

	return append(elements, element), nil
}

func (c *filteredCSVWriter) filterHeaders(headers []string) []string {
	filtered := make([]string, 0, len(headers))

	for _, header := range headers {
		if c.selection.includesHeader(header) {
			filtered = append(filtered, header)
		}
	}

	return filtered
}

func (c *filteredCSVWriter) filterRow(headers, row []string) []string {
	filtered := make([]string, 0, len(row))

	for i, header := range headers {
		if c.selection.includesHeader(header) {
			filtered = append(filtered, row[i])
		}
	}

	return filtered
}

func isSlice(t any) bool {
	if t == nil {
		return false
	}

	//nolint:exhaustive // only need to check for slices
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		return true
	default:
		return false
	}
}
