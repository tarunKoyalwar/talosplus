package stringutils

import (
	"regexp"
	"strings"
)

var variablereg *regexp.Regexp = regexp.MustCompile(`(([@]{1}[0-9a-zA-Z]+[{][a-z!]+[}])|([@]{1}[0-9a-zA-Z]+))`)

// SplitAtSpace : Similar to strings.Feilds but only considers ' '
func SplitAtSpace(s string) []string {

	return Split(s, ' ')

}

// Split : Similar to Strings.Feilds with Custom separator
func Split(s string, delim rune) []string {

	// Must trim the string first
	s = strings.TrimSpace(s)

	arr := []string{}

	var sb strings.Builder

	for _, v := range s {
		if v != delim {
			sb.WriteRune(v)
		} else {
			if sb.Len() != 0 {
				arr = append(arr, sb.String())
				sb.Reset()
			}
		}
	}

	if sb.Len() != 0 {
		arr = append(arr, sb.String())
		sb.Reset()
	}

	return arr
}

// UniqueArray : Get Array with unique items
func UniqueArray(s ...string) []string {

	z := map[string]bool{}

	for _, v := range s {
		//split string
		tarr := strings.Fields(v)
		for _, elem := range tarr {
			z[elem] = true
		}
	}

	keys := make([]string, 0, len(z))
	for k := range z {
		keys = append(keys, k)
	}

	return keys
}

// UniqueElements : Get Unique Elements
func UniqueElements(s ...string) string {
	u := map[string]bool{}
	for _, v := range s {
		for _, b := range Split(v, '\n') {
			u[b] = true
		}
	}

	tmp := ""
	for k := range u {
		if tmp == "" {
			tmp = k
		} else {
			tmp += "\n" + k
		}
	}

	return tmp
}

// ExtractVar :  Extract var name using regex
func ExtractVar(raw string) string {

	/*
		Handle edge cases like printf "@somevar\n"
	*/

	return variablereg.FindString(raw)
}
