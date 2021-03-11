package main

import (
	"fmt"
	"git-diff/diffparser"
	"io/ioutil"
)

//为简洁起见，未进行错误处理
func main() {
	byt, _ := ioutil.ReadFile("example.diff")
	//解析diff 文件
	diff, _ := diffparser.Parse(string(byt))
	fmt.Println("diff=======", diff.Changed())

	fmt.Println("===================================")
	// You now have a slice of files from the diff,
	//for k := range diff.Files {
	//	file := diff.Files[k]
	//	fmt.Println("file================", file)
	//	//fmt.Println("file.DiffHeader=====", file.DiffHeader)
	//	//fmt.Println("file.NewName=====", file.NewName)
	//	//fmt.Println("file.OrigName=====", file.OrigName)
	//	fmt.Println("file.Hunks=====", file.Hunks)
	//	fmt.Println("file.Mode=====", file.Mode)
	//}
	//file := diff.Files[0]
	//
	//// diff hunks in the file,
	//hunk := file.Hunks[0]
	//
	//// new and old ranges in the hunk
	//newRange := hunk.NewRange
	//
	//// and lines in the ranges.
	//line := newRange.Lines[0]
	//fmt.Println("line", line)

}
