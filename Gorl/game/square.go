package game

import (
	"fmt"
	"iter"
)

// square is a board position, {rank, file}, both 0-based: {ONE, A} is a1.
// It's also used for offsets between squares (move directions).
type square [2]int32

func (s square) Rank() int32 { return s[0] }
func (s square) File() int32 { return s[1] }

// Add offsets the square by d, a {ranks, files} offset.
func (s square) Add(d square) square {
	return square{s[0] + d[0], s[1] + d[1]}
}

func (s square) OnBoard() bool {
	return s[0] >= 0 && s[0] < 8 && s[1] >= 0 && s[1] < 8
}

// String returns the algebraic name, e.g. "e4".
func (s square) String() string {
	if !s.OnBoard() {
		return fmt.Sprintf("square{%d, %d}", s[0], s[1])
	}
	return fmt.Sprintf("%c%d", 'a'+s.File(), s.Rank()+1)
}

// allSquares yields every square on the board, rank by rank from a1 to h8.
func allSquares() iter.Seq[square] {
	return func(yield func(square) bool) {
		for r := range int32(8) {
			for f := range int32(8) {
				if !yield(square{r, f}) {
					return
				}
			}
		}
	}
}
