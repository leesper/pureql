package token

import "sort"

// Position describes a source code position including the file name, line,
// column and offset.
type Position struct {
	Filename string
	Line     int
	Column   int
	Offset   int
}

// Pos is a compact encoding of a source position within a file set.
// It can be converted into a Position for a more convenient, but much
// larger, representation.
//
// The Pos value for a given file is a number in the range [base, base+size],
// where base and size are specified when adding the file to the file set.
type Pos int

// The zero value for Pos is NoPos.
const (
	NoPos Pos = 0
)

// FileSet represents a set of *.graphql files.
type FileSet struct {
	base  int // base offset of the next file
	files []*File
}

// NewFileSet returns *FileSet.
func NewFileSet() *FileSet {
	return &FileSet{
		base: 1,
	}
}

// AddFile adds a new file with the given name and file size to the file set s
// and returns the file. The size must not be negative.
func (s *FileSet) AddFile(name string, size int) *File {
	if size < 0 {
		panic("invalid base or size")
	}

	file := &File{
		set:   s,
		name:  name,
		base:  s.base,
		size:  size,
		lines: []int{0},
	}

	// calculating base offset for the next file, +1 for EOF
	s.base += size + 1
	if s.base < 0 {
		panic("source files too large")
	}
	s.files = append(s.files, file)

	return file
}

// File returns file that contains the position p. If cannot find returns nil.
func (s *FileSet) File(p Pos) *File {
	if p == NoPos {
		return nil
	}
	i := sort.Search(len(s.files), func(i int) bool {
		return s.files[i].base <= int(p) && int(p) <= s.files[i].base+s.files[i].size
	})
	return s.files[i]
}

// Position converts Pos into a Position value.
func (s *FileSet) Position(p Pos) Position {
	f := s.File(p)
	if f == nil {
		return Position{}
	}
	return f.position(p)
}

// File is a handle for *.graphql files in FileSet.
type File struct {
	set   *FileSet
	name  string
	base  int
	size  int
	lines []int
}

// AddLine adds the line offset for a new line. The offset must be greater than
// the previous one and less than file size, otherwise it will be ignored.
func (f *File) AddLine(offset int) {
	if i := len(f.lines); (i == 0 || offset > f.lines[i-1]) && offset < f.size {
		f.lines = append(f.lines, offset)
	}
}

// Pos returns a Pos based on offset, offset must not be greater than file size.
func (f *File) Pos(offset int) Pos {
	if offset > f.size {
		panic("offset out of range")
	}
	return Pos(f.base + offset)
}

// Name returns file name.
func (f *File) Name() string {
	return f.name
}

// Size returns file size in bytes.
func (f *File) Size() int {
	return f.size
}

// LineCount returns file line count.
func (f *File) LineCount() int {
	return len(f.lines)
}

func (f *File) position(p Pos) Position {
	offset := f.offset(p)
	i := sort.Search(len(f.lines), func(i int) bool {
		return offset < f.lines[i]
	}) - 1

	return Position{
		Filename: f.name,
		Line:     i + 1,
		Column:   offset - f.lines[i] + 1,
		Offset:   offset,
	}
}

func (f *File) offset(p Pos) int {
	if int(p) < f.base || int(p) > f.base+f.size {
		panic("illegal Pos value")
	}
	return int(p) - f.base
}
