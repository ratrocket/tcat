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
			// row[0] is the entire row.
			b.WriteString(row[0])
			continue
		}

		// Get rid of blank columns hanging off the end, but
		// only if osep is blank.  If osep isn't blank, we still
		// want to represent the presence of a blank final
		// column.
		//
		// Note this only matters for non-default
		// (default=whitespace) input separator; it's impossible
		// to have a blank last column if your columns are
		// separated by whitespace.  You just have uneven
		// columns (which is fine for tcat).
		for *osep == "" && len(row) > 0 && row[len(row)-1] == "" {
			row = row[:len(row)-1]
		}

		for i, col := range row {

			// Write the contents...
			b.WriteString(col)

			// ... then, unless it's the last column (in
			// which case we break right here), add
			// whatever's necessary, taking into account the
			// width of the content relative to this
			// column's max width, and "osep".
			if i == len(row)-1 {
				break
			}

			// First write spaces up until max[i].
			for j := runewidth.StringWidth(col); j < max[i]; j++ {
				b.WriteRune(' ')
			}

			// Next, if osep is blank, add 2 spaces (so the
			// widest column is separated from its
			// neighbor-to-the-right by a *total* of two
			// spaces).
			//
			// If osep isn't blank, add a space to either
			// side of it and write that.  (With a little
			// nuance around the last column being blank...
			// keep reading the comments.)

			// This does the right thing whether osep is
			// blank or not.
			b.WriteString(" " + *osep)

			// Prevent writing a trailing space in the case
			// that we're on the 2nd to last column and the
			// last column is blank.  (Note that with input
			// separated by whitespace, it's impossible to
			// have a "blank" last column, and you can never
			// reach this code.)
			if i == len(row)-2 && row[len(row)-1] == "" {
				// nop
			} else {
				b.WriteRune(' ')
			}

			// The above is equivalent (by DeMorgan's law)
			// to the following (but I prefer the above; the
			// intent seems *much* clearer to me):
			// if i < len(row)-2 || row[len(row)-1] != "" {
			// 	b.WriteRune(' ')
			// }

		}

		b.WriteRune('\n')
	}
	b.Flush()
}
