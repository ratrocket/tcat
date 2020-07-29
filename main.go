// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tcat is a tabular cat.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	rows [][]string
	c    = flag.Bool("c", false, "use comma for input separator")
	p    = flag.Bool("p", false, "use pipe ('|') for input separator")
	isep = flag.String("is", "", "specify `input separator`")
	osep = flag.String("os", "", "specify `output separator`")

	sections = flag.Bool("sections", false, "allow section separator")
	secsep   = flag.String("secsep", "----", "prefix of section separator")
)

func usage() {
	// Must note that "-c" takes precedence over "-p" which takes
	// precedence over "-is".
	fmt.Fprintf(os.Stderr, "usage: tcat [-cp] [-is <string>] [-os <string>] [file...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("tcat: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if *p {
		*isep = "|"
	}
	if *c {
		*isep = ","
	}
	if flag.NArg() == 0 {
		read(os.Stdin)
	} else {
		for _, arg := range flag.Args() {
			f, err := os.Open(arg)
			if err != nil {
				log.Print(err)
				continue
			}
			read(f)
			f.Close()
		}
	}
	printTable(os.Stdout, rows)
}

func read(r io.Reader) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		log.Print(err)
	}
	for _, line := range strings.SplitAfter(string(data), "\n") {
		if line == "" {
			continue
		}
		row := split(line)
		rows = append(rows, row)
	}
}

func split(s string) []string {
	if *sections && strings.HasPrefix(s, *secsep) {
		return []string{s}
	}
	if *isep != "" {
		return strings.Split(strings.TrimSpace(s), *isep)
	}
	return strings.Fields(s)
}

func printTable(w io.Writer, rows [][]string) {
	var max []int
	for _, row := range rows {
		if *sections && strings.HasPrefix(row[0], *secsep) {
			continue
		}
		for i, c := range row {
			n := utf8.RuneCountInString(c)
			if i >= len(max) {
				max = append(max, n)
			} else if max[i] < n {
				max[i] = n
			}
		}
	}

	b := bufio.NewWriter(w)
	for _, row := range rows {
		if *sections && strings.HasPrefix(row[0], *secsep) {
			b.WriteString(row[0])
			continue
		}
		for len(row) > 0 && row[len(row)-1] == "" {
			row = row[:len(row)-1]
		}
		for i, c := range row {
			b.WriteString(c)
			if i+1 < len(row) {
				for j := utf8.RuneCountInString(c); j < max[i]+2; j++ {
					b.WriteRune(' ')
				}
			}
			if *osep != "" && i < len(row)-1 {
				b.WriteString(*osep)
			}
		}
		b.WriteRune('\n')
	}
	b.Flush()
}
