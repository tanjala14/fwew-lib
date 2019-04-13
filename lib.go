//	This file is part of Fwew.
//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/

// Package main contains all the things. lib.go handles common functions.
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

// Contains returns true if anything in q is also in s
func Contains(s []string, q []string) bool {
	if len(q) == 0 || len(s) == 0 {
		return false
	}
	for _, x := range q {
		for _, y := range s {
			if y == x {
				return true
			}
		}
	}
	return false
}

// ContainsStr returns true if q is in s
func ContainsStr(s []string, q string) bool {
	if len(q) == 0 || len(s) == 0 {
		return false
	}
	for _, x := range s {
		if q == x {
			return true
		}
	}
	return false
}

// DeleteElement "deletes" all occurrences of q in s
// actually returns a new string slice containing the original minus all q
func DeleteElement(s []string, q string) []string {
	if len(s) == 0 {
		return s
	}
	var r []string
	for _, str := range s {
		if str != q {
			r = append(r, str)
		}
	}
	return r
}

// DeleteEmpty "deletes" all empty string entries in s
// actually returns a new string slice containing all non-empty strings in s
func DeleteEmpty(s []string) []string {
	return DeleteElement(s, "")
}

// Index return the index of q in s
func Index(s []string, q string) int {
	for i, v := range s {
		if v == q {
			return i
		}
	}
	return -1
}

// IndexStr return the index of q in s
func IndexStr(s string, q rune) int {
	for i, v := range s {
		if v == q {
			return i
		}
	}
	return -1
}

// IsLetter returns true if s is an alphabet character or apostrophe
func IsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || r == '\'' || r == '‘' {
			return true
		}
	}
	return false
}

// Reverse returns the reversed version of s
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes[n:])
}

// StripChars strips all the characters in chr out of str
func StripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

// Valid validates range of integers for input
func Valid(input int64, reverse bool) bool {
	const (
		maxIntDec int64 = 32767
		maxIntOct int64 = 77777
	)
	if reverse {
		if 0 <= input && input <= maxIntDec {
			return true
		}
		return false
	}
	if 0 <= input && input <= maxIntOct {
		return true
	}
	return false
}

// DownloadDict downloads the latest released version of the dictionary file
func DownloadDict() error {
	var (
		filepath = Text("dictionary")
		url      = Text("dictURL")
	)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = out.Close()
	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	fmt.Println(Text("dlSuccess"))
	return nil
}

// https://github.com/ryanuber/go-glob
// The character which is treated like a glob
const GLOB = "*"

// Glob will test a string pattern, potentially containing globs, against a
// subject string. The result is a simple true/false, determining whether or
// not the glob pattern matched the subject text.
func Glob(pattern, subj string) bool {
	// Empty pattern can only match empty subject
	if pattern == "" {
		return subj == pattern
	}

	// If the pattern _is_ a glob, it matches everything
	if pattern == GLOB {
		return true
	}

	parts := strings.Split(pattern, GLOB)

	if len(parts) == 1 {
		// No globs in pattern, so test for equality
		return subj == pattern
	}

	leadingGlob := strings.HasPrefix(pattern, GLOB)
	trailingGlob := strings.HasSuffix(pattern, GLOB)
	end := len(parts) - 1

	// Go over the leading parts and ensure they match.
	for i := 0; i < end; i++ {
		idx := strings.Index(subj, parts[i])

		switch i {
		case 0:
			// Check the first section. Requires special handling.
			if !leadingGlob && idx != 0 {
				return false
			}
		default:
			// Check that the middle parts match.
			if idx < 0 {
				return false
			}
		}

		// Trim evaluated text from subj as we loop over the pattern.
		subj = subj[idx+len(parts[i]):]
	}

	// Reached the last section. Requires special handling.
	return trailingGlob || strings.HasSuffix(subj, parts[end])
}

func SHA1Hash(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))[0:8]
}
