package main

import (
	"fmt"
	"git-diff/diffparser"
	"io/ioutil"
)

//为简洁起见，未进行错误处理
func main() {

	byt, _ := ioutil.ReadFile("example.diff")
	diff, _ := diffparser.Parse(string(byt))

	//现在，您可以从差异中获取文件的一部分，
	for k := range diff.Files {
		//获取diff 头
		file := diff.Files[k]
		for k1 := range file.Hunks {
			hunk := file.Hunks[k1]
			// new and old ranges in the hunk
			newRange := hunk.NewRange

			// and lines in the ranges.
			line := newRange.Lines[0]
			fmt.Println(line)
		}
	}
}
