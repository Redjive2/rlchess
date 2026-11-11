package game

type (
	board struct {
		squares    [8][8]*piece
		selected   square
		isSelected bool
	}

	piece struct {
		color         colorKind
		kind          pieceKind
		HasMoved      bool // moved at least once this game (for castling)
		JustMoved     bool // moved on the most recent turn (e.g. for en passant)
		DoubleStepped bool // pawn's last move was a two-square advance
	}

	pieceKind  byte
	colorKind  byte
	gameStatus byte
)

func (c colorKind) Opponent() colorKind {
	if c == colorKindWhite {
		return colorKindBlack
	}

	return colorKindWhite
}

const (
	pieceKindPawn pieceKind = iota
	pieceKindRook
	pieceKindKnight
	pieceKindBishop
	pieceKindKing
	pieceKindQueen
)

const (
	colorKindWhite colorKind = iota
	colorKindBlack
)

const (
	A = iota
	B
	C
	D
	E
	F
	G
	H
)

const (
	ONE = iota
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
)

const (
	gameStatusPlayable gameStatus = iota
	gameStatusStalemate
	gameStatusCheckmate
)

// Occupied reports whether there is a piece at pos, indexed like squares
// (squares[pos.Rank()][pos.File()]). Positions off the board are never occupied.
func (b *board) Occupied(pos square) bool {
	return b.pieceAt(pos) != nil
}

// OccupiedBy reports whether there is a piece of the given color at pos.
func (b *board) OccupiedBy(pos square, color colorKind) bool {
	p := b.pieceAt(pos)
	return p != nil && p.color == color
}

// JustOccupied reports whether there is a piece at pos that moved on the most
// recent turn.
func (b *board) JustOccupied(pos square) bool {
	p := b.pieceAt(pos)
	return p != nil && p.JustMoved
}

// Ranks returns every square in rank-major order: Ranks()[rank][file].
func (b *board) Ranks() [8][8]**piece {
	var ranks [8][8]**piece
	for s := range allSquares() {
		ranks[s.Rank()][s.File()] = b.At(s)
	}
	return ranks
}

// Files returns every square in file-major order: Files()[file][rank].
func (b *board) Files() [8][8]**piece {
	var files [8][8]**piece
	for s := range allSquares() {
		files[s.File()][s.Rank()] = b.At(s)
	}
	return files
}

// At returns the square at pos, which can be read or assigned through, or
// nil if pos is off the board.
func (b *board) At(pos square) **piece {
	if !pos.OnBoard() {
		return nil
	}
	return &b.squares[pos.Rank()][pos.File()]
}

// pieceAt returns the piece at pos, or nil if the square is empty or off the
// board.
func (b *board) pieceAt(pos square) *piece {
	if square := b.At(pos); square != nil {
		return *square
	}
	return nil
}
