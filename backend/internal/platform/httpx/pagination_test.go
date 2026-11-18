package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParsePageQuery(t *testing.T) {
	valid := map[string]PageQuery{
		"":                         {Page: 1, PageSize: 20},
		"?page=3":                  {Page: 3, PageSize: 20},
		"?page=2&page_size=100":    {Page: 2, PageSize: 100},
		"?page_size=1&other=thing": {Page: 1, PageSize: 1},
	}
	for query, want := range valid {
		got, err := ParsePageQuery(httptest.NewRequest(http.MethodGet, "/items"+query, nil))
		if err != nil || got != want {
			t.Errorf("%q = %+v, %v; want %+v", query, got, err, want)
		}
	}
	invalid := map[string]string{
		"?page=0":                 "page",
		"?page=x":                 "page",
		"?page_size=0":            "page_size",
		"?page_size=101":          "page_size",
		"?page=-1&page_size=1000": "page,page_size",
	}
	for query, want := range invalid {
		_, err := ParsePageQuery(httptest.NewRequest(http.MethodGet, "/items"+query, nil))
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.Status != http.StatusUnprocessableEntity {
			t.Fatalf("%q err = %v, want a validation error", query, err)
		}
		var fields []string
		for _, f := range apiErr.Errors {
			fields = append(fields, f.Field)
		}
		if strings.Join(fields, ",") != want {
			t.Errorf("%q fields = %v, want %s", query, fields, want)
		}
	}
}

func TestPageQueryOffset(t *testing.T) {
	if got := (PageQuery{Page: 3, PageSize: 20}).Offset(); got != 40 {
		t.Errorf("Offset() = %d, want 40", got)
	}
}

func TestNewPageResponse(t *testing.T) {
	q := PageQuery{Page: 2, PageSize: 20}
	cases := map[int]int{0: 0, 1: 1, 20: 1, 21: 2, 42: 3}
	for total, pages := range cases {
		if got := NewPageResponse([]string{"a"}, q, total).TotalPages; got != pages {
			t.Errorf("total %d: TotalPages = %d, want %d", total, got, pages)
		}
	}
	body, err := json.Marshal(NewPageResponse[string](nil, q, 0))
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"items":[],"page":2,"page_size":20,"total_items":0,"total_pages":0}`; string(body) != want {
		t.Errorf("body %s, want %s", body, want)
	}
}
