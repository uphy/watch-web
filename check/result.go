package check

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type (
	Result struct {
		JobID    string
		Label    string
		Link     string
		Previous string
		Current  string
	}
	Diff struct {
		diff []diffmatchpatch.Diff
	}
)

func (r *Result) Diff() *Diff {
	d := diffmatchpatch.New()
	a, b, c := d.DiffLinesToChars(r.Previous, r.Current)
	diffs := d.DiffMain(a, b, false)
	diff := d.DiffCharsToLines(diffs, c)
	return &Diff{diff}
}

func (d *Diff) Changed() bool {
	if len(d.diff) == 0 {
		return false
	}
	for _, d := range d.diff {
		if d.Type == diffmatchpatch.DiffEqual {
			continue
		}
		return true
	}
	return false
}

func (d *Diff) String() string {
	w := new(bytes.Buffer)
	for _, diff := range d.diff {
		text := diff.Text
		text = strings.Trim(text, "\r\n")
		texts := strings.Split(text, "\n")
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			for _, t := range texts {
				fmt.Fprintf(w, "- %s\n", t)
			}
		case diffmatchpatch.DiffEqual:
			for _, t := range texts {
				fmt.Fprintf(w, "  %s\n", t)
			}
		case diffmatchpatch.DiffInsert:
			for _, t := range texts {
				fmt.Fprintf(w, "+ %s\n", t)
			}
		}
	}
	return w.String()
}
