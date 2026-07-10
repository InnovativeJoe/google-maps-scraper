package web

import (
	"bytes"
	"context"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/gosom/scrapemate"
)

func TestJobDataUnmarshalDefaultsFieldSelection(t *testing.T) {
	t.Parallel()

	var data JobData
	raw := []byte(`{"keywords":["coffee"],"lang":"en","depth":1,"max_time":60000000000,"email":true}`)

	if err := data.UnmarshalJSON(raw); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if !data.InputID || !data.PopularTimes || !data.Cid || !data.Status {
		t.Fatalf("expected optional field selection to default to true: %#v", data.FieldSelection)
	}

	if !data.Link || !data.Title || !data.Category || !data.Address {
		t.Fatalf("expected common result fields to default to true: %#v", data.FieldSelection)
	}

	if !data.OpenHours || !data.Website || !data.Phone || !data.PlusCode {
		t.Fatalf("expected contact fields to default to true: %#v", data.FieldSelection)
	}

	if !data.ReviewCount || !data.ReviewRating {
		t.Fatalf("expected review fields to default to true: %#v", data.FieldSelection)
	}

	if !data.Owner || !data.CompleteAddress {
		t.Fatalf("expected remaining optional fields to default to true: %#v", data.FieldSelection)
	}

	if !data.Email {
		t.Fatalf("expected explicit fields to remain intact")
	}
}

func TestFilteredCsvWriterFiltersColumns(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	writer := NewFilteredCsvWriter(csv.NewWriter(&buf), FieldSelection{
		Link:            true,
		InputID:         false,
		PopularTimes:    true,
		Cid:             false,
		Status:          true,
		Descriptions:    false,
		ReviewsLink:     true,
		Thumbnail:       false,
		DataID:          false,
		StreetViewURL:   false,
		PlaceID:         false,
		Images:          false,
		Reservations:    false,
		OrderOnline:     false,
		Menu:            false,
		Owner:           false,
		CompleteAddress: false,
	})

	input := make(chan scrapemate.Result, 1)
	input <- scrapemate.Result{
		Data: testCsvCapable{
			headers: []string{"input_id", "link", "popular_times", "cid", "status", "reviews_link"},
			row:     []string{"1", "https://example.com", "busy", "abc", "open", "reviews"},
		},
	}
	close(input)

	if err := writer.Run(context.Background(), input); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := strings.TrimSpace(buf.String())
	want := "link,popular_times,status,reviews_link\nhttps://example.com,busy,open,reviews"
	if got != want {
		t.Fatalf("unexpected csv output\nwant: %q\ngot:  %q", want, got)
	}
}

type testCsvCapable struct {
	headers []string
	row     []string
}

func (t testCsvCapable) CsvHeaders() []string {
	return t.headers
}

func (t testCsvCapable) CsvRow() []string {
	return t.row
}
