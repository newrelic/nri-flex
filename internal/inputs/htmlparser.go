/*
* Copyright 2019 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
 */

package inputs

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/newrelic/nri-flex/internal/formatter"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Table holds a simple table of headers and rows.
type Table struct {
	Attributes map[string]string
	Headers    []string
	Rows       [][]string
}

// ParseToJSON parses a html fragment or whole document looking for HTML
func ParseToJSON(s []byte, htmlAttributes map[string]string) (string, error) {

	tables, err := Parse(s, htmlAttributes)
	if err != nil {
		return "", err
	}
	jsonString := "["
	closeString := ""
	numberOfTables := len(tables)

	for i, t := range tables {
		if len(t.Rows) == 0 {
			continue
		}
		if i != numberOfTables-1 {
			closeString = ","
		} else {
			closeString = ""
		}
		tString := convertTable(t, i)
		jsonString = jsonString + tString + closeString
	}
	jsonString = jsonString + "]"
	return jsonString, nil
}

func convertTable(t *Table, i int) string {
	j := `{"table":[`

	header := t.Headers
	numberOfRows := len(t.Rows)
	for i, row := range t.Rows {

		j = j + "{"
		numberOfCells := len(row)
		for c := range row {
			r := ""
			if c != numberOfCells-1 {
				r = fmt.Sprintf(" \"%s\": \"%s\",", header[c], row[c])
			} else {
				r = fmt.Sprintf(" \"%s\": \"%s\"", header[c], row[c])
			}

			j = j + r
		}
		if i != numberOfRows-1 {
			j = j + "},"
		} else {
			j = j + "}]"
		}

	}
	for k, v := range t.Attributes {
		kv := fmt.Sprintf(", \"%s\": \"%s\"", k, v)
		j = j + kv

	}
	j = fmt.Sprintf("%s,\"Index\":%d }", j, i)
	return j
}

// Parse parses a html fragment or whole document looking for HTML
// tables. It converts all cells into text, stripping away any HTML content.
func Parse(s []byte, htmlAttributes map[string]string) ([]*Table, error) {
	node, err := html.Parse(bytes.NewReader(s))
	if err != nil {
		return nil, err
	}
	tables := []*Table{}
	var vThead = true
	parse(node, &tables, vThead, htmlAttributes)
	for kk, t := range tables {

		tables[kk] = addMissingColumns(t)
	}

	return tables, nil
}

func innerText(n *html.Node, parseAttribute bool, htmlAttributes map[string]string) string {
	if n.Type == html.TextNode {
		stripResult := stripChars(n.Data)
		return stripResult
	}
	var result string = ""
	if n.Type == html.ElementNode {
		if parseAttribute {
			result = parseAttributes(n.Attr, htmlAttributes)
		}
	}
	for x := n.FirstChild; x != nil; x = x.NextSibling {
		result += innerText(x, parseAttribute, htmlAttributes)
	}
	return result
}

func stripChars(input string) string {
	var result []string

	for _, i := range input {
		switch {
		// all these considered as space, including tab \t
		// '\t', '\n', '\v', '\f', '\r',' ', 0x85, 0xA0
		case unicode.IsSpace(i):
			result = append(result, " ") // replace tab with space
		case unicode.IsPunct(i):
			result = append(result, " ") // replace tab with space
		default:
			result = append(result, string(i))
		}
	}
	return strings.Join(result, "")
}

func containTable(n *html.Node) bool {
	if n.DataAtom == atom.Table {
		return true
	}
	result := false
	for x := n.FirstChild; x != nil && !result; x = x.NextSibling {
		result = containTable(x)
	}
	return result
}

func parse(n *html.Node, tables *[]*Table, vThead bool, htmlAttributes map[string]string) {
	strip := strings.TrimSpace
	switch n.DataAtom {
	case atom.Table:
		t := &Table{}
		for _, at := range n.Attr {
			if t.Attributes == nil {
				t.Attributes = map[string]string{}
			}
			t.Attributes[at.Key] = at.Val
		}
		*tables = append(*tables, t)
		vThead = true
	case atom.Th:
		if vThead {
			t := (*tables)[len(*tables)-1]
			t.Headers = append(t.Headers, strip(innerText(n, false, htmlAttributes)))
		} else {
			if !containTable(n) {
				t := (*tables)[len(*tables)-1]
				l := len(t.Rows) - 1
				t.Rows[l] = append(t.Rows[l], strip(innerText(n, true, htmlAttributes)))
				return
			}
			t := (*tables)[len(*tables)-1]
			l := len(t.Rows) - 1
			// If the <td> contains <table> element, set the <td> content to "TableElement"
			t.Rows[l] = append(t.Rows[l], "TableElement")
		}

	case atom.Tr:
		t := (*tables)[len(*tables)-1]
		t.Rows = append(t.Rows, []string{})
	case atom.Td:
		if !containTable(n) {
			t := (*tables)[len(*tables)-1]
			l := len(t.Rows) - 1
			t.Rows[l] = append(t.Rows[l], strip(innerText(n, true, htmlAttributes)))
			return
		}
		t := (*tables)[len(*tables)-1]
		l := len(t.Rows) - 1
		// If the <td> contains <table> element, set the <td> content to "TableElement"
		t.Rows[l] = append(t.Rows[l], "TableElement")

	case atom.Thead:
		vThead = true
	case atom.Tbody:
		vThead = false
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		parse(child, tables, vThead, htmlAttributes)
	}
}

func addMissingColumns(t *Table) *Table {
	cols := len(t.Headers)
	rows := make([][]string, 0, len(t.Rows))
	for _, row := range t.Rows {
		if len(row) > 0 {
			rows = append(rows, row)
		}
		if len(row) > cols {
			cols = len(row)
		}
	}
	for len(t.Headers) < cols {
		name := "Col " + strconv.Itoa(len(t.Headers)+1)
		t.Headers = append(t.Headers, name)
	}
	for i, row := range t.Headers {
		if len(row) == 0 {
			row = "Col " + strconv.Itoa(i)
		}
		t.Headers[i] = row
	}
	for kk := range rows {
		for len(rows[kk]) < cols {
			rows[kk] = append(rows[kk], "")
		}
	}
	t.Rows = rows
	return t
}

func parseAttributes(input []html.Attribute, parseAttributes map[string]string) string {
	var result []string
	for _, attr := range input {
		for key, val := range parseAttributes {
			if formatter.KvFinder("regex", attr.Key, key) {
				if formatter.KvFinder("regex", attr.Val, val) {
					result = append(result, attr.Key+":"+attr.Val+";")
				}
			}
		}

	}
	return strings.Join(result, "")
}
