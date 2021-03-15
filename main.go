package main

import (
	"bufio"
	"fmt"
	"git-diff/Const"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//为简洁起见，未进行错误处理
func main() {
	read()
}

type CodeBlock struct {
	StartLine int64
	EndLine   int64
}

type Resp struct {
	File string
	C    []CodeBlock
}

//写文件
func write() {
	content := []byte(Const.DiffFile)
	err := ioutil.WriteFile("test.txt", content, 0644)
	if err != nil {
		panic(err)
	}
}

//按照行读取文件  [3,6]
func read() {
	var CodeBlocks []CodeBlock
	var CodeBlock CodeBlock
	filepath := "/Users/chenbo/Documents/project/git-diff/test.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()

	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	//获取起始行数
	lineNum, err := StareLine()
	fmt.Println("代码块起始行数为：", lineNum)
	if err != nil {
		fmt.Println(err)
	}

	var k int64
	//循环的次数
	var y int64
	//上一行的状态 1 代表 空行 2 代表 + 3 代表-
	var status int64
	fmt.Println("status", status)
	// 代码块开始的行数
	var start int64
	//代码块结束的行数
	var end int64
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if IsStateToLine(line) == 2 {
			//判断start 是否有值
			if start == 0 {
				//+1表示 循环从0次开始 但是对应的代码行数是1
				start = start + k + 1
				end = start
				//判断下一行是否是连续的 + y+1 表示的是当前的循环次数（从0开始循环）y+2 表示下一行
				if getStatus(y+2) != 2 {
					//如果不是+则end =start
					CodeBlock.StartLine = start
					CodeBlock.EndLine = end
					CodeBlocks = append(CodeBlocks, CodeBlock)
					//此代码块结束 清零
					start = 0
					end = 0
				} else {
					//如果是连续的
					end = end + 1
				}

			} else {
				if getStatus(y+2) != 2 {
					//如果不是+则end =start
					CodeBlock.StartLine = start
					CodeBlock.EndLine = end
					CodeBlocks = append(CodeBlocks, CodeBlock)
					//此代码块结束 清零
					start = 0
					end = 0
				} else {
					//如果是连续的
					end = end + 1
				}

			}

		}
		if IsStateToLine(line) == 3 {
			k--
		}
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		k++
		y++
		status = IsStateToLine(line)
	}

	fmt.Println("=======", CodeBlocks)
}

//返回代码块增加的行数  +的行数
func AddLine() int64 {
	//解析正则表达式，如果成功返回解释器
	reg1 := regexp.MustCompile(`\n\+[^\n]+[\n$]|\+ +.+$`)
	if reg1 == nil {
		fmt.Println("regexp err")
		return 0
	}
	//根据规则提取关键信息
	result1 := reg1.FindAllStringSubmatch(Const.DiffFile, -1)
	return int64(len(result1))
}

//返回@@ -1,5 +1,7 @@ 中的 1
func StareLine() (stare int64, err error) {

	diff := strings.Split(Const.DiffFile, "@@")
	diff1 := strings.Split(diff[1], " ")
	diff2 := strings.Split(diff1[2], ",")
	str := strings.Replace(diff2[0], " ", "", -1)
	if stare, err = strconv.ParseInt(str, 10, 64); err != nil {
		fmt.Println("==========error=========", err)
	}
	return

	//for k, v := range diff1[2] {
	//	fmt.Println("k=============", k)
	//	fmt.Println("v=============", v)
	//}
	return
}

//判断每一行的代码的第一位是+ - 还是 ""
//1 代表 空行 2 代表 + 3 代表-
func IsStateToLine(line string) int64 {
	if len(line) == 0 {
		return 1
	}
	if string(line[0]) == "+" {
		return 2
	}
	if string(line[0]) == "-" {
		return 3
	}
	return 1
}

//获取指定代码块行数的状态  1 代表 空行 2 代表 + 3 代表-
func getStatus(y int64) int64 {
	var m map[int64]int64
	m = make(map[int64]int64)

	filepath := "/Users/chenbo/Documents/project/git-diff/test.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return 0
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()

	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)

	var k int64
	for {
		k++
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		m[k] = IsStateToLine(line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read file error!", err)
				return 0
			}
		}

	}
	return m[y]
}
