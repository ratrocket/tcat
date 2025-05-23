# tcat

Forked from rsc.io/tcat in order to add option to split input on commas
(-c option).

Later extended to split on pipes ("|"; `-p`), tabs (`-t`), and also
anything else via `--is` (see usage).

Tcat is a table cat.

Usage:

    tcat [-cpt] [--is <string>] [--os <string>] [file...]
      --comma/-c use comma for input separator
      --pipe/-p  use pipe ('|') for input separator
      --tab/-t   use tab for input separator
      --is       specify input separator
      --os       specify output separator
      --sections (bool) allow sections or not
      --secsep   prefix of section separator (default "----")
    ("-t" takes precedence over "-c" which takes precedence over "-p"
    which takes precedence over "-is")

Tcat reads the named input files, splits each line into space-(or
comma-)separate fields, and then reprints the input aligning columns of
fields.

"Sections" can be accommodated by passing `--sections=true` and then any
line starting with "secep" (default "----") will be passed through tcat
untouched.

See it in action:

    go build
    ./tcat testdata/spaces
    ./tcat -c testdata/commas
    ./tcat -p testdata/pipes
    ./tcat -t testdata/tabs

## License

BSD 3-Clause "New" or "Revised" License

(This repository retains the license from the original repository it was
forked from, rsc.io/tcat.  That repository's license is copied verbatim.
See LICENSE for details.)
