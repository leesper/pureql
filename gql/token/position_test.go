package token

import (
	"reflect"
	"testing"
)

func TestInvalidSize(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("should panic")
		}
	}()
	fs := NewFileSet()
	fs.AddFile("invalidSize", -1)
}

func TestFiles(t *testing.T) {
	tests := []struct {
		filename string
		source   string
		size     int
		lines    []int // lines contains the offset of the first character for each line (the first entry is always 0)
	}{
		{"a", "", 0, []int{0}},
		{"b", "01234", 5, []int{0}},
		{"c", "\n\n\n\n\n\n\n\n\n", 9, []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
		{"d", "package p\n\nimport \"fmt\"", 23, []int{0, 10, 11}},
		{"e", "package p\n\nimport \"fmt\"\n", 24, []int{0, 10, 11}},
		{"f", "package p\n\nimport \"fmt\"\n ", 25, []int{0, 10, 11, 24}},
	}

	fs := NewFileSet()
	for _, test := range tests {
		f := fs.AddFile(test.filename, test.size)
		if f.Name() != test.filename {
			t.Errorf("expecting %s, found %s", test.filename, f.Name())
		}

		if f.Size() != test.size {
			t.Errorf("expecting %d, found %d", test.size, f.Size())
		}

		if fs.File(f.Pos(0)) != f {
			t.Errorf("expecting %v, found %v", f, fs.File(f.Pos(0)))
		}

		for index, offset := range test.lines {
			f.AddLine(offset)
			if f.LineCount() != index+1 {
				t.Errorf("expecting %d, found %d", index+1, f.LineCount())
			}

			// add again should be ignored
			f.AddLine(offset)
			if f.LineCount() != index+1 {
				t.Errorf("expecting %d, found %d", index+1, f.LineCount())
			}

			verifyPosition(t, fs, f, test.lines[0:index+1])
		}

		if f.LineCount() != len(test.lines) {
			t.Errorf("expecting %d, found %d", len(test.lines), f.LineCount())
		}

		verifyPosition(t, fs, f, test.lines)
	}
}

func verifyPosition(t *testing.T, fs *FileSet, f *File, lines []int) {
	for offset := 0; offset < f.Size(); offset++ {
		p := f.Pos(offset)
		offset2 := f.offset(p)
		if offset2 != offset {
			t.Errorf("expecting %d, found %d", offset, offset2)
		}

		line, column := linecol(offset, lines)
		expect := Position{
			Filename: f.Name(),
			Line:     line,
			Column:   column,
			Offset:   offset,
		}
		pos := fs.Position(p)

		if !reflect.DeepEqual(expect, pos) {
			t.Errorf("expecting %v, found %v", expect, pos)
		}
	}
}

func linecol(offset int, lines []int) (int, int) {
	prevLineOff := 0
	for index, lineOff := range lines {
		if offset < lineOff {
			return index, offset - prevLineOff + 1
		}
		prevLineOff = lineOff
	}
	return len(lines), offset - prevLineOff + 1
}
