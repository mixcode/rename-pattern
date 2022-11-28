package main

import (
	"testing"
)

/*
func TestSearchPattern(t *testing.T) {

	filenames := []string{
		"abcde_fghij_klmno",
		"pqrst_fghij_klmno",
		"abcde_pqrst_klmno",
		"abcde_fghij_pqrst",
	}

	type testp struct {
		Pattern1, Pattern2 string
		Result             []string
	}
	patterns := []testp{
		{"abcde*", "xyz*", []string{
			"xyz_fghij_klmno",
			"pqrst_fghij_klmno",
			"xyz_pqrst_klmno",
			"xyz_fghij_pqrst",
		}},
		{"*fghij*", "*xyz*", []string{
			"abcde_xyz_klmno",
			"pqrst_xyz_klmno",
			"abcde_pqrst_klmno",
			"abcde_xyz_pqrst",
		}},
		{"*klmno", "*xyz", []string{
			"abcde_fghij_xyz",
			"pqrst_fghij_xyz",
			"abcde_pqrst_xyz",
			"abcde_fghij_pqrst",
		}},
		{"?bcde*", "*xyz*", []string{
			"axyz_fghij_klmno",
			"pqrst_fghij_klmno",
			"axyz_pqrst_klmno",
			"axyz_fghij_pqrst",
		}},
		{"*?ghij*", "*?xyz*", []string{
			"abcde_fxyz_klmno",
			"pqrst_fxyz_klmno",
			"abcde_pqrst_klmno",
			"abcde_fxyz_pqrst",
		}},
		{"*fghi?*", "*xyz?*", []string{
			"abcde_xyzj_klmno",
			"pqrst_xyzj_klmno",
			"abcde_pqrst_klmno",
			"abcde_xyzj_pqrst",
		}},
		{"*klmn?", "*xyz?", []string{
			"abcde_fghij_xyzo",
			"pqrst_fghij_xyzo",
			"abcde_pqrst_xyzo",
			"abcde_fghij_pqrst",
		}},
		{"*klmn?", "*xyz?", []string{
			"abcde_fghij_xyzo",
			"pqrst_fghij_xyzo",
			"abcde_pqrst_xyzo",
			"abcde_fghij_pqrst",
		}},
		{"*bcd*_*", "*xyz", []string{ // test missing patterns
			"axyze_fghij_klmno",
			"pqrst_fghij_klmno",
			"axyze_pqrst_klmno",
			"axyze_fghij_pqrst",
		}},
		{"*fghij**", "*xyz**", []string{ // test muliple asterisk ignoring
			"abcde_xyz_klmno",
			"pqrst_xyz_klmno",
			"abcde_pqrst_klmno",
			"abcde_xyz_pqrst",
		}},
	}

	for _, p := range patterns {
		for i, s := range filenames {
			f, e := Rename_OLD(s, p.Pattern1, p.Pattern2)
			if e != nil {
				t.Error(e)
				continue
			}
			if f != p.Result[i] {
				t.Errorf("pattern not match: expected %s, actual %s", p.Result[i], s)
			}
		}

	}
}
*/

func TestRenamePattern(t *testing.T) {

	filenames := []string{
		"abcde_fghij_klmno_12345",
		"pqrst_fghij_klmno_12345",
		"abcde_pqrst_klmno_12345",
		"abcde_fghij_pqrst_12345",
		"abcde_fghij_klmno_pqrst",
	}

	type testp struct {
		Pattern1, Pattern2 string
		Result             []string
	}
	patterns := []testp{
		{"abcde*", "xyz*", []string{
			"xyz_fghij_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"xyz_pqrst_klmno_12345",
			"xyz_fghij_pqrst_12345",
			"xyz_fghij_klmno_pqrst",
		}},
		{"ab???_*", "xy???_*", []string{
			"xycde_fghij_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"xycde_pqrst_klmno_12345",
			"xycde_fghij_pqrst_12345",
			"xycde_fghij_klmno_pqrst",
		}},
		{"*fghij*", "*xyz*", []string{
			"abcde_xyz_klmno_12345",
			"pqrst_xyz_klmno_12345",
			"abcde_pqrst_klmno_12345",
			"abcde_xyz_pqrst_12345",
			"abcde_xyz_klmno_pqrst",
		}},
		{"*12345", "*xyz", []string{
			"abcde_fghij_klmno_xyz",
			"pqrst_fghij_klmno_xyz",
			"abcde_pqrst_klmno_xyz",
			"abcde_fghij_pqrst_xyz",
			"abcde_fghij_klmno_pqrst",
		}},
		{"?bcde*", "*xyz*", []string{
			"axyz_fghij_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"axyz_pqrst_klmno_12345",
			"axyz_fghij_pqrst_12345",
			"axyz_fghij_klmno_pqrst",
		}},
		{"*?ghij*", "*?xyz*", []string{
			"abcde_fxyz_klmno_12345",
			"pqrst_fxyz_klmno_12345",
			"abcde_pqrst_klmno_12345",
			"abcde_fxyz_pqrst_12345",
			"abcde_fxyz_klmno_pqrst",
		}},
		{"*fghi?*", "*xyz?*", []string{
			"abcde_xyzj_klmno_12345",
			"pqrst_xyzj_klmno_12345",
			"abcde_pqrst_klmno_12345",
			"abcde_xyzj_pqrst_12345",
			"abcde_xyzj_klmno_pqrst",
		}},
		{"abcde|_*_*", "*_$1_*", []string{ // substitute test
			"abcde_abcde_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"abcde_abcde_klmno_12345",
			"abcde_abcde_pqrst_12345",
			"abcde_abcde_klmno_pqrst",
		}},
		{"abcde|_*_*", "*_*($1)_$5", []string{ // substitute test
			"abcde_fghij(abcde)_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"abcde_pqrst(abcde)_klmno_12345",
			"abcde_fghij(abcde)_pqrst_12345",
			"abcde_fghij(abcde)_klmno_pqrst",
		}},
		{"abcde|_*_*", "*_*($1)_*", []string{ // substitute test
			"abcde_fghij(abcde)_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"abcde_pqrst(abcde)_klmno_12345",
			"abcde_fghij(abcde)_pqrst_12345",
			"abcde_fghij(abcde)_klmno_pqrst",
		}},
		{"abcde|_*_*", "*_*(${1}23)_*", []string{ // substitute test
			"abcde_fghij(abcde23)_klmno_12345",
			"pqrst_fghij_klmno_12345",
			"abcde_pqrst(abcde23)_klmno_12345",
			"abcde_fghij(abcde23)_pqrst_12345",
			"abcde_fghij(abcde23)_klmno_pqrst",
		}},
		{"*_|12345", "*_%09d", []string{ // substitute test
			"abcde_fghij_klmno_000012345",
			"pqrst_fghij_klmno_000012345",
			"abcde_pqrst_klmno_000012345",
			"abcde_fghij_pqrst_000012345",
			"abcde_fghij_klmno_pqrst",
		}},
		{"*_:", "*_%09d", []string{ // substitute test
			"abcde_fghij_klmno_000012345",
			"pqrst_fghij_klmno_000012345",
			"abcde_pqrst_klmno_000012345",
			"abcde_fghij_pqrst_000012345",
			"abcde_fghij_klmno_pqrst",
		}},
	}

	for _, p := range patterns {
		for i, s := range filenames {
			p1, e := compilePattern(p.Pattern1)
			if e != nil {
				t.Error(e)
				continue
			}
			p2, e := parseReplacePattern(p.Pattern2)
			if e != nil {
				t.Error(e)
				continue
			}
			//log.Printf("p1: (%s)[%v]", p.Pattern1, p1)
			//log.Printf("p2: (%s)[%v]", p.Pattern2, p2)
			f, e := ReplaceName(s, p1, p2, 0) // zero offset
			if e != nil {
				t.Error(e)
				continue
			}
			if f != p.Result[i] {
				t.Errorf("pattern not match: expected %s, actual %s", p.Result[i], s)
			}
		}

	}
}
