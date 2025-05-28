# tcat

> The authoritative version of `tcat` (of my fork, that is) lives on
> [sourcehut](https://git.sr.ht/~md0/tcat).  Anything else is a mirror
> and may lag behind.
>
> The import path is md0.org/tcat.

## Installation

```
go install md0.org/tcat@latest
```

## md0.org fork

Forked from rsc.io/tcat in order to add option to split input on commas
(-c option).

Later extended to split on pipes ("|"; `-p`), tabs (`-t`), and also
anything else via `--is` (see usage).

## modified original README

Tcat is a tabular cat.

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

## Examples

All of these are based on the testscript/txtar tests in
testdata/script/, named, e.g., "spaces.txtar".

### spaces

```
> cat spaces.txt
a b c
cat dogs   birds
space   tab     other
```

```
> tcat spaces.txt  # separating on spaces is the default
a      b     c
cat    dogs  birds
space  tab   other
```

### commas

```
> cat commas.txt
label0,label1,label2,label3
Alice,45,6,cats|coffee
Bob,32,8,dogs|walking|water
```

```
> tcat -c  # or --comma
label0  label1  label2  label3
Alice   45      6       cats|coffee
Bob     32      8       dogs|walking|water
```

### pipes

```
> cat pipes.txt
sushi|toast|crackers
cat|dog|badger
vim|emacs|ed
```

```
> tcat -p  # or --pipe
sushi  toast  crackers
cat    dog    badger
vim    emacs  ed
```

### tabs

```
> cat tabs.txt
col1    col2    longcol 3 with a lot of different kind of space in it   col4
cat     dog     person  warthog
this is a column        this is also a column   this is a column        and... so is this
```

```
> tcat -t  # or --tab
col1              col2                   longcol 3 with a lot of different kind of space in it  col4
cat               dog                    person                                                 warthog
this is a column  this is also a column  this is a column                                       and... so is this
```

## License

BSD 3-Clause "New" or "Revised" License

(This repository retains the license from the original repository it was
forked from, rsc.io/tcat.  That repository's license is copied verbatim.
See LICENSE for details.)
