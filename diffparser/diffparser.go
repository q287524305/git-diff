package diffparser

import (
	"regexp"
	"strconv"
	"strings"

	"errors"
)

// FileMode represents the file status in a diff
// FileMode表示差异中的文件状态
type FileMode int

const (
	// DELETED if the file is deleted
	// 如果文件已删除，则删除
	DELETED FileMode = iota
	// MODIFIED if the file is modified
	// 已修改（如果文件被修改）
	MODIFIED
	// NEW if the file is created and there is no diff
	//如果文件已创建且没有差异，则为NEW
	NEW
)

// DiffRange contains the DiffLine's
//DiffRange包含多个DiffLine
type DiffRange struct {

	//起始行号
	Start int

	// 行数
	Length int

	// 代码块
	Lines []*DiffLine
}

// DiffLineMode tells the line if added, removed or unchanged
// DiffLineMode告诉该行是添加，删除还是保持不变
type DiffLineMode rune

const (
	// ADDED if the line is added (shown green in diff)
	////添加（如果添加了该行）（差异显示为绿色）
	ADDED DiffLineMode = iota
	// REMOVED if the line is deleted (shown red in diff)
	//已删除（如果删除了该行）（差异显示为红色）
	REMOVED
	// UNCHANGED if the line is unchanged (not colored in diff)
	//如果该行未更改（未在diff中着色），则保持不变
	UNCHANGED
)

// DiffLine is the least part of an actual diff
// DiffLine是实际差异的最小部分
type DiffLine struct {
	Mode     DiffLineMode // DiffLineMode告诉该行是添加，删除还是保持不变
	Number   int
	Content  string
	Position int // the line in the diff //差异中的行

}

// DiffHunk is a group of difflines
// DiffHunk是一组difflines
type DiffHunk struct {
	HunkHeader string
	OrigRange  DiffRange
	NewRange   DiffRange
	WholeRange DiffRange
}

// DiffFile is the sum of diffhunks and holds the changes of the file features
//DiffFile 是diffhunks的总和，用于保存文件功能的更改
type DiffFile struct {
	DiffHeader string
	Mode       FileMode
	OrigName   string
	NewName    string
	Hunks      []*DiffHunk
}

// Diff is the collection of DiffFiles
// Diff是DiffFiles的集合
type Diff struct {
	Files []*DiffFile
	Raw   string `sql:"type:text"`

	PullID uint `sql:"index"`
}

func (d *Diff) addFile(file *DiffFile) {
	d.Files = append(d.Files, file)
}

// Changed returns a map of filename to lines changed in that file. Deleted
// files are ignored.

// Changed返回文件名映射到该文件中更改的行。已删除
//文件将被忽略。
func (d *Diff) Changed() map[string][]int {
	dFiles := make(map[string][]int)

	for _, f := range d.Files {
		if f.Mode == DELETED {
			continue
		}

		for _, h := range f.Hunks {
			for _, dl := range h.NewRange.Lines {
				if dl.Mode == ADDED { // TODO(waigani) return removed
					dFiles[f.NewName] = append(dFiles[f.NewName], dl.Number)
				}
			}
		}
	}

	return dFiles
}

func regFind(s string, reg string, group int) string {
	re := regexp.MustCompile(reg)
	return re.FindStringSubmatch(s)[group]
}

func lineMode(line string) (*DiffLineMode, error) {
	var m DiffLineMode
	switch line[:1] {
	case " ":
		m = UNCHANGED
	case "+":
		m = ADDED
	case "-":
		m = REMOVED
	default:
		return nil, errors.New("could not parse line mode for line: \"" + line + "\"")
	}
	return &m, nil
}

// Parse takes a diff, such as produced by "git diff", and parses it into a
// Diff struct.

//解析会获取一个差异（例如由“ git diff”产生的差异），并将其解析为
//区分结构。
func Parse(diffString string) (*Diff, error) {
	var diff Diff
	diff.Raw = diffString
	lines := strings.Split(diffString, "\n")

	var file *DiffFile
	var hunk *DiffHunk
	var ADDEDCount int
	var REMOVEDCount int
	var inHunk bool
	oldFilePrefix := "--- a/"
	newFilePrefix := "+++ b/"

	var diffPosCount int
	var firstHunkInFile bool
	// Parse each line of diff.
	for idx, l := range lines {
		diffPosCount++
		switch {
		case strings.HasPrefix(l, "diff "):
			inHunk = false

			// Start a new file.
			file = &DiffFile{}
			header := l
			if len(lines) > idx+3 {
				rein := regexp.MustCompile(`^index .+$`)
				remp := regexp.MustCompile(`^(-|\+){3} .+$`)
				index := lines[idx+1]
				if rein.MatchString(index) {
					header = header + "\n" + index
				}
				mp1 := lines[idx+2]
				mp2 := lines[idx+3]
				if remp.MatchString(mp1) && remp.MatchString(mp2) {
					header = header + "\n" + mp1 + "\n" + mp2
				}
			}
			file.DiffHeader = header
			diff.Files = append(diff.Files, file)
			firstHunkInFile = true

			// File mode.
			file.Mode = MODIFIED
		case l == "+++ /dev/null":
			file.Mode = DELETED
		case l == "--- /dev/null":
			file.Mode = NEW
		case strings.HasPrefix(l, oldFilePrefix):
			file.OrigName = strings.TrimPrefix(l, oldFilePrefix)
		case strings.HasPrefix(l, newFilePrefix):
			file.NewName = strings.TrimPrefix(l, newFilePrefix)
		case strings.HasPrefix(l, "@@ "):
			if firstHunkInFile {
				diffPosCount = 0
				firstHunkInFile = false
			}

			inHunk = true
			// Start new hunk.
			hunk = &DiffHunk{}
			file.Hunks = append(file.Hunks, hunk)

			// Parse hunk heading for ranges
			re := regexp.MustCompile(`@@ \-(\d+),?(\d+)? \+(\d+),?(\d+)? @@ ?(.+)?`)
			m := re.FindStringSubmatch(l)
			if len(m) < 5 {
				return nil, errors.New("Error parsing line: " + l)
			}
			a, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, err
			}
			b := a
			if len(m[2]) > 0 {
				b, err = strconv.Atoi(m[2])
				if err != nil {
					return nil, err
				}
			}
			c, err := strconv.Atoi(m[3])
			if err != nil {
				return nil, err
			}
			d := c
			if len(m[4]) > 0 {
				d, err = strconv.Atoi(m[4])
				if err != nil {
					return nil, err
				}
			}
			if len(m[5]) > 0 {
				hunk.HunkHeader = m[5]
			}

			// hunk orig range.
			hunk.OrigRange = DiffRange{
				Start:  a,
				Length: b,
			}

			// hunk new range.
			hunk.NewRange = DiffRange{
				Start:  c,
				Length: d,
			}

			// (re)set line counts
			ADDEDCount = hunk.NewRange.Start
			REMOVEDCount = hunk.OrigRange.Start
		case inHunk && isSourceLine(l):
			m, err := lineMode(l)
			if err != nil {
				return nil, err
			}
			line := DiffLine{
				Mode:     *m,
				Content:  l[1:],
				Position: diffPosCount,
			}
			newLine := line
			origLine := line

			// add lines to ranges
			switch *m {
			case ADDED:
				newLine.Number = ADDEDCount
				hunk.NewRange.Lines = append(hunk.NewRange.Lines, &newLine)
				hunk.WholeRange.Lines = append(hunk.WholeRange.Lines, &newLine)
				ADDEDCount++

			case REMOVED:
				origLine.Number = REMOVEDCount
				hunk.OrigRange.Lines = append(hunk.OrigRange.Lines, &origLine)
				hunk.WholeRange.Lines = append(hunk.WholeRange.Lines, &origLine)
				REMOVEDCount++

			case UNCHANGED:
				newLine.Number = ADDEDCount
				hunk.NewRange.Lines = append(hunk.NewRange.Lines, &newLine)
				hunk.WholeRange.Lines = append(hunk.WholeRange.Lines, &newLine)
				origLine.Number = REMOVEDCount
				hunk.OrigRange.Lines = append(hunk.OrigRange.Lines, &origLine)
				ADDEDCount++
				REMOVEDCount++
			}
		}
	}

	return &diff, nil
}

func isSourceLine(line string) bool {
	if line == `\ No newline at end of file` {
		return false
	}
	if l := len(line); l == 0 || (l >= 3 && (line[:3] == "---" || line[:3] == "+++")) {
		return false
	}
	return true
}

// Length returns the hunks line length
//长度返回粗线的长度

func (hunk *DiffHunk) Length() int {
	return len(hunk.WholeRange.Lines) + 1
}
