// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tcat is a tabular cat.
package main

import (
	"os"

	"md0.org/tcat/internal/tabularcat"
)

func main() {
	os.Exit(tabularcat.Main())
}
