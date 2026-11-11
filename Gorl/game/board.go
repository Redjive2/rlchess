package game

func (b *board) GetSelection() (square, bool) {
	return b.selected, b.isSelected
}

func (b *board) Select(s square) {
	b.selected, b.isSelected = s, true
}

func (b *board) Deselect() {
	b.isSelected = false
}

// Move moves the piece at from to to, assuming it's one of its move options.
// It handles captures, the en passant capture and the castling rook, and
// updates the moved-piece flags.
func (b *board) Move(from, to square) {
	p := b.pieceAt(from)

	// only the piece moving now counts as having just moved
	for s := range allSquares() {
		if other := b.pieceAt(s); other != nil {
			other.JustMoved = false
		}
	}

	b.squares = b.afterMove(from, to).squares

	p.JustMoved = true
	p.HasMoved = true
	p.DoubleStepped = p.isDoubleStep(from, to)

	// castling also moves the rook, onto the square the king passed over
	if p.isCastle(from, to) {
		rook := b.pieceAt(square{from.Rank(), (from.File() + to.File()) / 2})
		rook.HasMoved = true
	}
}

func SetupBoard() *board {
	b := &board{}

	backRank := [8]pieceKind{
		pieceKindRook, pieceKindKnight, pieceKindBishop, pieceKindQueen,
		pieceKindKing, pieceKindBishop, pieceKindKnight, pieceKindRook,
	}

	for file, kind := range backRank {
		f := int32(file)
		*b.At(square{ONE, f}) = &piece{color: colorKindWhite, kind: kind}
		*b.At(square{TWO, f}) = &piece{color: colorKindWhite, kind: pieceKindPawn}
		*b.At(square{SEVEN, f}) = &piece{color: colorKindBlack, kind: pieceKindPawn}
		*b.At(square{EIGHT, f}) = &piece{color: colorKindBlack, kind: kind}
	}

	return b
}
