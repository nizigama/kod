package editor

import (
	"bufio"
	"io"
	"unicode/utf8"
)

type Pos struct {
	X, Y int
}

type Line struct {
	data []byte
}

type LineArray struct {
	lines []*Line
}

// converts a position of a rune(on the screen) to it's position in the byte array
// graciously stolen from micro
func runeToByteIndex(txt []byte, n int) int {
	if n == 0 {
		return 0
	}

	count := 0
	i := 0
	for len(txt) > 0 {
		_, size := utf8.DecodeRune(txt)

		txt = txt[size:]
		count += size
		i++

		if i == n {
			break
		}
	}
	return count
}

func NewLineArray(in io.Reader) *LineArray {
	// TODO: Warn about large files
	la := LineArray{}

	// alloc 1000 lines by default
	la.lines = make([]*Line, 0, 1000)

	br := bufio.NewReader(in)

	numLines := 0

	for {
		data, err := br.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				la.AppendLine()
				numLines++
				la.lines[numLines-1].data = data[:len(data)]
			}
			break
		} else {
			la.AppendLine()
			numLines++
			la.lines[numLines-1].data = data[:len(data)-1]
		}
	}

	return &la
}

func (la *LineArray) AppendLine() {
	la.lines = append(la.lines, &Line{data: []byte{}})
}

func (la *LineArray) InsertLineAfter(pos Pos) {
	// insert after current line
	pos.Y++
	la.AppendLine()
	// shift everything to the right
	copy(la.lines[pos.Y+1:], la.lines[pos.Y:])
	la.lines[pos.Y] = &Line{}
}

func (la *LineArray) RemoveLine(pos int) {
	la.lines = append(la.lines[:pos], la.lines[pos+1:]...)
}

func (la *LineArray) Split(pos Pos) {
	la.InsertLineAfter(pos)
	src := la.lines[pos.Y]
	la.Insert(src.data[pos.X:], Pos{0, pos.Y + 1})
	la.DeleteToEnd(pos)
}

func (la *LineArray) JoinLines(a int, b int) {
	line := la.lines[a]
	la.Insert(la.lines[b].data, Pos{len(line.data), a})
	la.RemoveLine(b)
}

func (la *LineArray) Insert(value []byte, pos Pos) {
	x, y := runeToByteIndex(la.lines[pos.Y].data, pos.X), pos.Y
	for i := 0; i < len(value); i++ {
		if value[i] == '\n' {
			la.Split(Pos{x, y})
			x = 0
			y++
			continue
		}
		la.insertByte(value[i], Pos{x, y})
		x++
	}
}

func (la *LineArray) insertByte(value byte, pos Pos) {
	line := la.lines[pos.Y]
	line.data = append(line.data, 0)
	// shift everything to the right
	copy(line.data[pos.X+1:], line.data[pos.X:])
	line.data[pos.X] = value
}

func (la *LineArray) DeleteToEnd(pos Pos) {
	la.lines[pos.Y].data = la.lines[pos.Y].data[:pos.X]
}

func (la *LineArray) RemoveRune(pos Pos) {
	line := la.lines[pos.Y]
	startX := runeToByteIndex(line.data, pos.X)
	endX := runeToByteIndex(line.data, pos.X+1)
	line.data = append(line.data[:startX], line.data[endX:]...)
}