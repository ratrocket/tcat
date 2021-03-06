Forked from rsc.io/tcat in order to add option to split input on commas.

Tcat is a table cat.

Usage:

    tcat [-cpt] [-is <string>] [-os <string>] [file...]
      -c        use comma for input separator
      -p        use pipe ('|') for input separator
      -t        use tab for input separator
      -is       specify input separator
      -os       specify output separator
      -sections (bool) allow sections or not
      -secsep   prefix of section separator (default "----")
    ("-c" takes precedence over "-p" which takes precedence which takes
    precedence over -t which takes precedence over "-is")

Tcat reads the named input files, splits each line into space-(or
comma-)separate fields, and then reprints the input aligning columns of
fields.

"Sections" can be accommodated by passing `-sections=true` and then any
line starting with "secep" (default "----") will be passed through tcat
untouched.

See it in action:

    go build
    ./tcat testdata/spaces
    ./tcat -c testdata/commas
    ./tcat -p testdata/pipes
    ./tcat -t testdata/tabs
