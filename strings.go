package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

// Strings maps a template name and a language in the form <template>/<lang>
// to a map of strings corresponding to keys in the template.
type Strings map[string]map[string]string

// LoadStrings reads template strings from all "templates/<template>.<lang>.strings" files
// prepending the strings from "templates/layout.<lang>.strings".
// This function will panic in case of errors since it is meant to run during initialization only.
func LoadStrings() Strings {
	files, err := filepath.Glob("templates/*.*.strings")
	if err != nil {
		panic(fmt.Errorf("strings: can't read files in 'templates/'"))
	}
	re := regexp.MustCompile(`templates/(\w+)\.(\w+)\.strings`)
	strings := make(Strings)
	for _, filename := range files {
		key := re.ReplaceAllString(filename, "$1/$2")
		strings[key] = loadStrings(filename)
	}
	// for each template/<lang>, add strings from layout/<lang> if not present
	re = regexp.MustCompile(`(?:\w+)/(\w+)`)
	for key, val := range strings {
		key := re.ReplaceAllString(key, "layout/$1")
		for k, v := range strings[key] {
			if val[k] == "" {
				val[k] = v
			}
		}
	}
	return strings
}

func loadStrings(filename string) map[string]string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("strings: can't read file %q", filename))
	}
	re := regexp.MustCompile(`(?m)^(\w+)\n(?:#.*\n)*"((?:[^"\\]*|\\["\\bfnrt\/]|\\u[0-9a-f]{4})*)"`)
	sm := re.FindAllSubmatch(content, -1)
	m := make(map[string]string)
	for _, p := range sm {
		m[string(p[1])] = string(p[2])
	}
	return m
}

func NewView(src map[string]string) map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
