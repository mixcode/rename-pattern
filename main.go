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
	useStdIn       = false // use stdin for file list
	wholeMatch     = false // the pattern must match the whole filename
	quiet          = false // suppress output
	verbose        = false // more output

	// filename matching regexp pattern
	matchRegexpPattern *regexp.Regexp
	matchPatternOffset int // index offset of matching patters. increases to 1 if an automatic leftmost match pattern is added.

	// replace pattern
	replacePattern []Pattern
)

func renameOneFile(filename string) error {
	d, f := filepath.Split(filename)
	newf, e := ReplaceName(f, matchRegexpPattern, replacePattern, matchPatternOffset)
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

	// build match pattern
	if useRegExpMatch { // input pattern is a regexp pattern
		s := a[0]
		if !wholeMatch {
			// leftmost and rightmost pattern should be added
			if len(s) > 0 && s[0] != '^' {
				s = "^(.*?)" + s
				matchPatternOffset = 1
			}
			if s[len(s)-1] != '$' {
				s = s + "(.*?)$"
			}
		}
		matchRegexpPattern, err = regexp.Compile(s)
	} else {
		s := a[0]
		if !wholeMatch {
			// leftmost and rightmost pattern should be added
			if len(s) > 0 && s[0] != '*' {
				s = "*" + s
				matchPatternOffset = 1
			}
			if s[len(s)-1] != '*' {
				s = s + "*"
			}
		}
		matchRegexpPattern, err = compilePattern(s)
	}
	if err != nil {
		return
	}

	// build replace pattern
	s := a[1]
	if !wholeMatch {
		if matchPatternOffset > 0 && len(s) > 0 && s[0] != '*' {
			// Prepend automatic match only if the automatic leftmost match added at the search
			s = "*" + s
		}
		if s[len(s)-1] != '*' {
			s += "*"
		}
	}
	replacePattern, err = parseReplacePattern(s)
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
		fmt.Fprintf(o, "  | : match group separator\n")
		fmt.Fprintf(o, "  (other chars) : match as-is\n")
		fmt.Fprintf(o, "  Or, use a regexp string with submatches paired by parentheses\n")
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
