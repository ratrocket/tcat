Forked from rsc.io/tcat in order to add option to split input on commas.

Tcat is a table cat.

Usage:

    tcat [-c] [file ...]

Tcat reads the named input files, splits each line into space-(or
comma-)separate fields, and then reprints the input aligning columns of
fields.

See it in action:

    go build
    ./tcat testdata/spaces
    ./tcat -c testdata/commas
    ./tcat -p testdata/pipes
