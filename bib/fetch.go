package bib

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
	"reflect"
)

const zoteroBase string = "https://api.zotero.org"

func SearchZotero(query string, userId string) ([]BibEntry, error) {
	var entries []BibEntry

	// Create the request URL
	u, _ := url.Parse(zoteroBase)
	u.Path = path.Join(u.Path, "users", userId, "items")

	// Add the query. Fetch bibliography entry as CSL JSON
	q := u.Query()
	q.Add("v", "3")
	q.Add("q", query)
	q.Add("format", "csljson")
	u.RawQuery = q.Encode()

	// Fetch the response
	resp, err := http.Get(u.String())
	if err != nil {
		return entries, err
	}
	defer resp.Body.Close()

	// Read into []byte
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entries, err
	}

	// Unmarshal "items" field into a list of raw JSON to parse individually
	var results struct {
		Items []json.RawMessage
	}
	json.Unmarshal(body, &results)

	// Do parsing of each result
	for _, entry := range results.Items {
		var e BibEntry

		// Unmarshal known fields
		json.Unmarshal(entry, &e)

		// Process dates
		var dateStruct struct {
			Accessed struct {
				DateParts [][]int	`json:"date-parts"`
			} `json:"accessed"`
			Issued struct {
				DateParts [][]int	`json:"date-parts"`
			} `json:"issued"`
			Submitted struct {
				DateParts [][]int	`json:"date-parts"`
			} `json:"accessed"`
		}

		// This is extremely verbose, but it works
		if err := json.Unmarshal(entry, &dateStruct); err == nil {
			process := func(d [][]int) time.Time {
				if len(d[0]) > 2 {
					return time.Date(d[0][0], time.Month(d[0][1]), d[0][2], 0, 0, 0, 0, time.UTC)
				} else if len(d[0]) == 2 {
					return time.Date(d[0][0], time.Month(d[0][1]), 1, 0, 0, 0, 0, time.UTC)
				} else {
					return time.Date(d[0][0], time.Month(1), 1, 0, 0, 0, 0, time.UTC)
				}
			}
			if d := dateStruct.Accessed; !reflect.ValueOf(d).IsZero() {
				e.Accessed = process(d.DateParts)
			}
			if d := dateStruct.Issued; !reflect.ValueOf(d).IsZero() {
				// NOTE: Assuming "date-parts". It seems Zotero parses dates themselves, and does not use "raw"
				e.Issued = process(d.DateParts)
			}
			if d := dateStruct.Submitted; !reflect.ValueOf(d).IsZero() {
				e.Submitted = process(d.DateParts)
			}
		}

		entries = append(entries, e)
	}

	return entries, nil
}
