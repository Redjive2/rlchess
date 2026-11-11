package game

import "slices"

func (p piece) GetMoveOptions(b *board, pos square) []square {
	var options []square

	switch p.kind {
	case pieceKindPawn:
		options = p.pawnOptions(b, pos)
	case pieceKindRook:
		options = p.rookOptions(b, pos)
	case pieceKindKnight:
		options = p.knightOptions(b, pos)
	case pieceKindBishop:
		options = p.bishopOptions(b, pos)
	case pieceKindQueen:
		options = p.queenOptions(b, pos)
	case pieceKindKing:
		options = p.kingOptions(b, pos)
	}

	return p.filterCheck(b, pos, clipToBoard(options))
}

// clipToBoard removes options that fall off the board, so the xOptions
// functions can propose squares without checking the edges themselves.
func clipToBoard(options []square) []square {
	return slices.DeleteFunc(options, func(s square) bool {
		return !s.OnBoard()
	})
}

// inCheck reports whether color's king is attacked. Boards without that king
// (e.g. in tests) are never in check.
func (b *board) inCheck(color colorKind) bool {
	pos, ok := b.find(pieceKindKing, color)
	return ok && b.attacked(pos, color.Opponent())
}

func (p piece) filterCheck(b *board, pos square, options []square) []square {
	var filtered []square

	for _, opt := range options {
		if p.isCastle(pos, opt) {
			passed := square{pos.Rank(), (pos.File() + opt.File()) / 2}
			if b.inCheck(p.color) || b.attacked(passed, p.color.Opponent()) {
				continue
			}
		}

		if b.afterMove(pos, opt).inCheck(p.color) {
			continue
		}

		filtered = append(filtered, opt)
	}

	return filtered
}

func (b *board) afterMove(from, to square) *board {
	next := *b
	p := next.pieceAt(from)

	// en passant: a pawn moving diagonally onto an empty square captures the
	// pawn beside its starting square
	if p.kind == pieceKindPawn && from.File() != to.File() && !next.Occupied(to) {
		*next.At(square{from.Rank(), to.File()}) = nil
	}

	// castling: the rook jumps to the square the king passed over
	if p.isCastle(from, to) {
		rookFile := int32(H)
		if to.File() < from.File() {
			rookFile = A
		}
		rookFrom := square{from.Rank(), rookFile}
		*next.At(square{from.Rank(), (from.File() + to.File()) / 2}) = next.pieceAt(rookFrom)
		*next.At(rookFrom) = nil
	}

	*next.At(to) = p
	*next.At(from) = nil
	return &next
}

func (b *board) find(kind pieceKind, color colorKind) (square, bool) {
	for s := range allSquares() {
		if b.is(s, kind, color) {
			return s, true
		}
	}

	return square{}, false
}

var knightOffsets = []square{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}

// attacked reports whether any piece of color by attacks pos. It works from
// pos outward instead of calling GetMoveOptions, so it can't recurse.
func (b *board) attacked(pos square, by colorKind) bool {

	// pawns attack diagonally forward, so look one rank behind pos (from by's side)
	behind := int32(-1)
	if by == colorKindBlack {
		behind = 1
	}
	if b.is(pos.Add(square{behind, -1}), pieceKindPawn, by) || b.is(pos.Add(square{behind, 1}), pieceKindPawn, by) {
		return true
	}

	for _, d := range knightOffsets {
		if b.is(pos.Add(d), pieceKindKnight, by) {
			return true
		}
	}

	for _, d := range allDirections {
		if b.is(pos.Add(d), pieceKindKing, by) {
			return true
		}
	}

	// sliders: the first piece along each line, if it's by's and moves that way
	lines := []struct {
		directions []square
		kind       pieceKind
	}{{straightDirections, pieceKindRook}, {diagonalDirections, pieceKindBishop}}

	for _, line := range lines {
		for _, d := range line.directions {
			s := pos.Add(d)
			for s.OnBoard() && !b.Occupied(s) {
				s = s.Add(d)
			}

			if b.is(s, line.kind, by) || b.is(s, pieceKindQueen, by) {
				return true
			}
		}
	}

	return false
}

// is reports whether the piece at pos has the given kind and color.
func (b *board) is(pos square, kind pieceKind, color colorKind) bool {
	p := b.pieceAt(pos)
	return p != nil && p.kind == kind && p.color == color
}

func abs(n int32) int32 {
	if n < 0 {
		return -n
	}

	return n
}

func (p piece) pawnOptions(b *board, pos square) []square {
	var (
		r, f = pos.Rank(), pos.File()

		// white moves up the ranks, black down
		forward       int32 = 1
		startRank     int32 = TWO
		enPassantRank int32 = FIVE
		opponent            = p.color.Opponent()
	)

	if p.color == colorKindBlack {
		forward, startRank, enPassantRank = -1, SEVEN, FOUR
	}

	var options []square

	// one square ahead, or two from the starting rank; both must be empty
	ahead := square{r + forward, f}
	if !b.Occupied(ahead) {
		options = append(options, ahead)

		twoAhead := square{r + 2*forward, f}
		if r == startRank && !b.Occupied(twoAhead) {
			options = append(options, twoAhead)
		}
	}

	for _, side := range []int32{-1, 1} {
		diagonal := square{r + forward, f + side}

		// regular capture
		if b.OccupiedBy(diagonal, opponent) {
			options = append(options, diagonal)
		}

		// en passant: an enemy pawn beside us that just advanced two squares can be
		// captured as if it had moved one, landing on the square it passed over
		beside := square{r, f + side}
		if r == enPassantRank && b.JustOccupied(beside) && b.OccupiedBy(beside, opponent) {
			if victim := b.pieceAt(beside); victim.kind == pieceKindPawn && victim.DoubleStepped {
				options = append(options, diagonal)
			}
		}
	}

	return options
}

var (
	straightDirections = []square{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	diagonalDirections = []square{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	allDirections      = slices.Concat(straightDirections, diagonalDirections)
)

// isCastle reports whether moving p between from and to is castling.
func (p piece) isCastle(from, to square) bool {
	return p.kind == pieceKindKing && from.Rank() == to.Rank() && abs(to.File()-from.File()) == 2
}

// isDoubleStep reports whether moving p between from and to is a pawn's
// two-square advance.
func (p piece) isDoubleStep(from, to square) bool {
	return p.kind == pieceKindPawn && from.File() == to.File() && abs(to.Rank()-from.Rank()) == 2
}

func (p piece) slide(b *board, pos square, directions []square) []square {
	var options []square

	for _, d := range directions {
		for s := pos.Add(d); ; s = s.Add(d) {
			if !s.OnBoard() || b.OccupiedBy(s, p.color) {
				break
			}

			options = append(options, s)
			if b.Occupied(s) { // opposing piece: capture, then stop
				break
			}
		}
	}

	return options
}

func (p piece) rookOptions(b *board, pos square) []square {
	return p.slide(b, pos, straightDirections)
}

func (p piece) knightOptions(b *board, pos square) []square {
	var options []square

	for _, d := range knightOffsets {
		if s := pos.Add(d); !b.OccupiedBy(s, p.color) {
			options = append(options, s)
		}
	}

	return options
}

func (p piece) bishopOptions(b *board, pos square) []square {
	return p.slide(b, pos, diagonalDirections)
}

func (p piece) queenOptions(b *board, pos square) []square {
	return p.slide(b, pos, allDirections)
}

func (p piece) kingOptions(b *board, pos square) []square {
	var options []square

	for _, d := range allDirections {
		if s := pos.Add(d); !b.OccupiedBy(s, p.color) {
			options = append(options, s)
		}
	}

	if p.HasMoved {
		return options
	}

	r := pos.Rank()
	for _, side := range []struct{ rookFile, step int32 }{{H, 1}, {A, -1}} {
		rook := b.pieceAt(square{r, side.rookFile})
		if rook == nil || rook.kind != pieceKindRook || rook.color != p.color || rook.HasMoved {
			continue
		}

		clear := true
		for f := pos.File() + side.step; f != side.rookFile; f += side.step {
			if b.Occupied(square{r, f}) {
				clear = false
				break
			}
		}

		if clear {
			options = append(options, square{r, pos.File() + 2*side.step})
		}
	}

	return options
}
