package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	doIt           = false // actually rename the files
	useRegExpMatch = false // use regexp match pattern
	useStdIn       = false
	quiet          = false
	verbose        = false
	wholeMatch     = false // the pattern must match the whole filename

	patternRegexp *regexp.Regexp
	patternTo     []Pattern
)

func renameOneFile(filename string) error {
	d, f := filepath.Split(filename)
	newf, e := ReplaceName(f, patternRegexp, patternTo)
	if e != nil {
		return e
	}
	changed := (newf != f)
	if changed {
		if !quiet {
			fmt.Printf("%s -> %s\n", f, newf)
		}
		if doIt {
			e = os.Rename(filename, filepath.Join(d, newf))
			if e != nil {
				return e
			}
		}
	} else {
		if verbose {
			fmt.Printf("[NOT CHANGE] %s\n", f)
		}
	}
	return nil
}

func run() (err error) {
	a := flag.Args()
	if len(a) < 2 {
		return fmt.Errorf("rename pattern must given")
	}

	extendPattern := func(s string) string { // append all-matching marks at the end
		if wholeMatch {
			return s
		}
		if len(s) == 0 {
			return "*"
		}
		if s[0] != '*' {
			s = "*" + s
		}
		if s[len(s)-1] != '*' {
			s += "*"
		}
		return s
	}

	if useRegExpMatch {
		s := a[0]
		if !wholeMatch {
			if len(s) > 0 && s[0] != '^' {
				s = "^(.*?)" + s
			}
			if s[len(s)-1] != '$' {
				s = s + "(.*?)$"
			}
		}
		patternRegexp, err = regexp.Compile(s)
	} else {
		patternRegexp, err = compilePattern(extendPattern(a[0]))
	}
	if err != nil {
		return
	}
	patternTo, err = parseReplacePattern(extendPattern(a[1]))
	if err != nil {
		return
	}
	//log.Printf("match from: [%v]", patternRegexp)
	//log.Printf("match to: [%v]", patternTo)

	if useStdIn {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			err = renameOneFile(sc.Text())
			if err != nil {
				return
			}
		}

	} else {
		a = a[2:]
		for _, f := range a {
			err = renameOneFile(f)
			if err != nil {
				return
			}
		}

	}
	return
}

func main() {

	flag.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Rename file matching certain patterns.\n")
		fmt.Fprintf(o, "\n")
		fmt.Fprintf(o, "Usage:\n")
		fmt.Fprintf(o, "  %s [flags] matchPattern replacePattern [filename ...]\n", os.Args[0])
		fmt.Fprintf(o, "  %s [flags] -s matchPattern replacePattern\n", os.Args[0])
		fmt.Fprintf(o, "\n")
		fmt.Fprintf(o, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\n")
		fmt.Fprintf(o, "Match Patterns (without -r flag):\n")
		fmt.Fprintf(o, "  * : match multiple chars\n")
		fmt.Fprintf(o, "  ? : match one char\n")
		fmt.Fprintf(o, "  : : match numbers\n")
		fmt.Fprintf(o, "  | : match separator. the match string is separated\n")
		fmt.Fprintf(o, "  (other chars) : match as-is\n")
		fmt.Fprintf(o, "  Or, use a regexp string with submatches paired by parentheses.\n")
		fmt.Fprintf(o, "\n")
		fmt.Fprintf(o, "Replace Patterns:\n")
		fmt.Fprintf(o, "  * or ? : use the matched pattern as-is\n")
		fmt.Fprintf(o, "  $NUM or ${NUM}: use the NUMth pattern matched\n")
		fmt.Fprintf(o, "  %%...: use printf-like formatting. %%d and %%s are valid, and %%[POS] could be used\n")
	}

	flag.BoolVar(&doIt, "d", doIt, "Do rename the files")
	flag.BoolVar(&useRegExpMatch, "r", useRegExpMatch, "use regexp match pattern")
	flag.BoolVar(&useStdIn, "s", useStdIn, "read filenames from stdin")
	flag.BoolVar(&quiet, "q", quiet, "suppress output")
	flag.BoolVar(&verbose, "v", verbose, "print some more info")
	flag.BoolVar(&wholeMatch, "w", wholeMatch, "the pattern must match the whole filename")
	flag.Parse()

	err := run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
