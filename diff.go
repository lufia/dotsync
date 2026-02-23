package main

import (
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func runDiff(r *Repository, args []string, w io.Writer) error {
	f := NewFlagSet("diff", "[path ...]")

	if err := f.Parse(args); err != nil {
		return err
	}
	where := make(map[string]struct{})
	for _, arg := range f.Args() {
		dest, err := filepath.Abs(arg)
		if err != nil {
			return err
		}
		where[dest] = struct{}{}
	}
	dir := filepath.Join(r.StateDir, "store")
	return filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		state, err := r.ReadState(p)
		if err != nil {
			return err
		}
		if _, ok := where[state.Target]; !ok && len(where) > 0 {
			return nil
		}
		h, err := ReadHash(state.Target)
		if err != nil {
			return err
		}
		ok, _, err := isModeEqual(state.Target, state.Mode)
		if err != nil {
			return err
		}
		if ok && h == state.Hash {
			return nil
		}
		return diffFiles(w, state.Source, state.Target)
	})
}

var (
	fileColor   = color.New(color.FgHiWhite, color.Bold)
	metaColor   = color.New(color.FgCyan)
	addColor    = color.New(color.FgGreen)
	deleteColor = color.New(color.FgRed)
)

func diffFiles(w io.Writer, f1, f2 string) error {
	t1, err := os.ReadFile(f1)
	if err != nil {
		return err
	}
	t2, err := os.ReadFile(f2)
	if err != nil {
		return err
	}
	fileColor.Fprintf(w, "--- %s\n", f1)
	fileColor.Fprintf(w, "+++ %s\n", f2)
	d := diffmatchpatch.New()
	a, b, c := d.DiffLinesToChars(string(t1), string(t2))
	diffs := d.DiffCharsToLines(d.DiffMain(a, b, false), c)
	for h := range reduceDiffs(diffs) {
		h.WriteTo(w)
	}
	return nil
}

type Hunk struct {
	// line number at the top of this hunk
	oldLine, oldLen int
	newLine, newLen int
	diffs           []*Diff
}

func (h *Hunk) WriteTo(w io.Writer) (int64, error) {
	metaColor.Fprintf(w, "@@ -%d,%d +%d,%d @@\n", h.oldLine, h.oldLen, h.newLine, h.newLen)
	var written int
	for _, d := range h.diffs {
		switch d.kind {
		case -1:
			n, _ := deleteColor.Fprintf(w, "-%s\n", d.text)
			written += n
		case 0:
			n, _ := fmt.Fprintf(w, " %s\n", d.text)
			written += n
		case 1:
			n, _ := addColor.Fprintf(w, "+%s\n", d.text)
			written += n
		}
	}
	return int64(written), nil
}

type Diff struct {
	kind int // -1: delete, 1: add, 0: equal
	text string
}

type pusher struct {
	prev  *Hunk
	yield func(*Hunk) bool
}

func (p *pusher) Push(h *Hunk) bool {
	if len(h.diffs) == 0 {
		return true
	}
	rv := true
	if p.prev != nil && hasModified(p.prev.diffs) {
		rv = p.yield(p.prev)
	}
	p.prev = h
	return rv
}

func (p *pusher) Flush() bool {
	rv := true
	if p.prev != nil {
		rv = p.yield(p.prev)
	}
	p.prev = nil
	return rv
}

func hasModified(diffs []*Diff) bool {
	return slices.ContainsFunc(diffs, func(d *Diff) bool {
		return d.kind != 0
	})
}

// reduceDiffs snips the middle of diff's text that has type=equal,
// and it returns sequence of a pair of current and next [Hunk]s.
func reduceDiffs(diffs []diffmatchpatch.Diff) iter.Seq[*Hunk] {
	return func(yield func(*Hunk) bool) {
		push := pusher{yield: yield}
		var (
			hunk    = &Hunk{oldLine: 1, newLine: 1}
			oldNext = 1
			newNext = 1
		)
		for _, d := range diffs {
			text := d.Text
			if text[len(text)-1] == '\n' {
				text = text[:len(text)-1]
			}
			a := strings.Split(text, "\n")
			if len(a) == 0 {
				continue
			}
			switch d.Type {
			case diffmatchpatch.DiffEqual:
				before, after := splitLines(a, 3)
				hunk.diffs = append(hunk.diffs, makeDiffs(before, 0)...)
				hunk.oldLen += len(before)
				hunk.newLen += len(before)
				if !push.Push(hunk) {
					return
				}
				hunk = &Hunk{
					oldLine: oldNext + len(a) - len(after),
					oldLen:  len(after),
					newLine: newNext + len(a) - len(after),
					newLen:  len(after),
					diffs:   makeDiffs(after, 0),
				}
				oldNext += len(a)
				newNext += len(a)
			case diffmatchpatch.DiffInsert:
				hunk.diffs = append(hunk.diffs, makeDiffs(a, 1)...)
				hunk.newLen += len(a)
				newNext += len(a)
			case diffmatchpatch.DiffDelete:
				hunk.diffs = append(hunk.diffs, makeDiffs(a, -1)...)
				hunk.oldLen += len(a)
				oldNext += len(a)
			default:
				panic("unknown type: " + d.Type.String())
			}
		}
		push.Push(hunk)
		push.Flush()
	}
}

func splitLines(a []string, n int) (before, after []string) {
	before = a
	if len(before) > n {
		before = before[:n]
	}
	after = a[len(before):]
	if len(after) > n {
		after = after[len(after)-n:]
	}
	return
}

func makeDiffs(a []string, kind int) []*Diff {
	diffs := make([]*Diff, len(a))
	for i, s := range a {
		diffs[i] = &Diff{kind: kind, text: s}
	}
	return diffs
}
