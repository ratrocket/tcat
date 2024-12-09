// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tcat is a tabular cat.
package tabularcat

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"md0.org/mflag"
)

var (
	rows [][]string

	c    *bool   // flag.Bool("c", false, "use comma for input separator")
	p    *bool   // flag.Bool("p", false, "use pipe ('|') for input separator")
	t    *bool   // flag.Bool("t", false, "use tab for input separator")
	isep *string // flag.String("is", "", "specify `input separator`")
	osep *string // flag.String("os", "", "specify `output separator`")

	sections *bool   // flag.Bool("sections", false, "allow section separator")
	secsep   *string // flag.String("secsep", "----", "prefix of section separator")
)

func Main() int {
	log.SetPrefix("tcat: ")
	log.SetFlags(0)
	fs := mflag.NewFlagSet(os.Args[0], mflag.ExitOnError)
	var (
		help = fs.BoolP("help", "h", false, "print help and exit")
	)

	c = fs.BoolP("comma", "c", false, "use comma for input separator")
	p = fs.BoolP("pipe", "p", false, "use pipe ('|') for input separator")
	t = fs.BoolP("tab", "t", false, "use tab for input separator")

	isep = fs.String("is", "", "specify `input separator`")
	osep = fs.String("os", "", "specify `output separator`")

	sections = fs.Bool("sections", false, "allow section separator")
	secsep = fs.String("secsep", "----", "prefix of section separator")

	fs.Usage = func() {
		// Must note that "-c" takes precedence over "-p" which takes
		// precedence over "-is".
		fmt.Fprintf(os.Stderr, "usage: tcat [-cpt] [--is <string>] [--os <string>] [file...]\n")
		fmt.Fprintf(os.Stderr, "       (also): --sections={true,false} --secsep string\n")
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[1:])
	if *help {
		fs.Usage()
		return 0
	}
	if *p {
		*isep = "|"
	}
	if *c {
		*isep = ","
	}
	if *t {
		*isep = "\t"
	}
	if fs.NArg() == 0 {
		read(os.Stdin)
	} else {
		for _, arg := range fs.Args() {
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

	return 0
}

func read(r io.Reader) {
	data, err := io.ReadAll(r)
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
		parts := []string{}
		for _, col := range strings.Split(strings.TrimSpace(s), *isep) {
			parts = append(parts, strings.TrimSpace(col))
		}
		return parts

	}
	return strings.Fields(s)
}

func printTable(w io.Writer, rows [][]string) {
	var max []int

	// figure out what the widest column in each row is.
	for _, row := range rows {
		if *sections && strings.HasPrefix(row[0], *secsep) {
			continue
		}
		for i, c := range row {
			n := runewidth.StringWidth(c)
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

		// get rid of blank columns hanging off the end
		for len(row) > 0 && row[len(row)-1] == "" {
			row = row[:len(row)-1]
		}

		for i, col := range row {
			b.WriteString(col)

			// TODO I think you can do an early "continue"
			// here if i == len(row)-1.  Both subsequent
			// "ifs" will change their condition.

			if i+1 < len(row) {
				for j := runewidth.StringWidth(col); j < max[i]+2; j++ {
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
