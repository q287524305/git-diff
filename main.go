package main

import (
	"fmt"
	"git-diff/diffparser"
	"io/ioutil"
)

// error handling left out for brevity
func main() {
	byt, _ := ioutil.ReadFile("example.diff")
	diff, _ := diffparser.Parse(string(byt))

	// You now have a slice of files from the diff,
	file := diff.Files[0]

	// diff hunks in the file,
	hunk := file.Hunks[0]

	// new and old ranges in the hunk
	newRange := hunk.NewRange

	// and lines in the ranges.
	line := newRange.Lines[0]
	fmt.Println("line=============", line)
}
