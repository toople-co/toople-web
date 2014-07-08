package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/simleb/errors"
)

// strings contains the strings for all templates in all languages.
var strings = LoadStrings()

// A View is a map with string keys used to populate templates.
type View map[string]interface{}

// NewView creates a view containing the strings of a template in a given language.
func NewView(tmpl, lang string) View {
	m := make(View)
	for k, v := range strings[tmpl][lang] {
		m[k] = v
	}
	return m
}

// LoadStrings reads template strings from all "templates/<template>.<lang>.strings" files
// prepending the strings from "templates/layout.<lang>.strings".
// This function will panic in case of errors since it is meant to run during initialization only.
func LoadStrings() map[string]map[string]map[string]string {
	files, err := filepath.Glob(filepath.Join("templates", "*.*.strings"))
	if err != nil {
		panic(errors.Stack(err, "strings: can't read strings files in directory 'templates'"))
	}
	strings := make(map[string]map[string]map[string]string)
	re := regexp.MustCompile(`(\w+)\.(\w+)\.strings`)
	// First pass: load all strings
	for _, filename := range files {
		m := re.FindStringSubmatch(filepath.Base(filename))
		strings[m[1]] = make(map[string]map[string]string)
		strings[m[1]][m[2]] = LoadStringsFile(filename)
	}
	// Second pass: complete with layout strings
	for tmpl := range strings {
		for lang := range strings[tmpl] {
			for k, v := range strings["layout"][lang] {
				if strings[tmpl][lang][k] == "" {
					strings[tmpl][lang][k] = v
				}
			}
		}
	}
	return strings
}

// LoadStringsFile reads template strings from a file.
// This function will panic in case of errors since it is meant to run during initialization only.
func LoadStringsFile(filename string) map[string]string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(errors.Stackf(err, "load strings: can't read file %q", filename))
	}
	re := regexp.MustCompile(`(?m)^(\w+)\n(?:#.*\n)*"((?:[^"\\]*|\\["\\bfnrt\/]|\\u[0-9a-f]{4})*)"`)
	sm := re.FindAllStringSubmatch(string(content), -1)
	strings := make(map[string]string)
	for _, m := range sm {
		strings[m[1]] = m[2]
	}
	return strings
}
