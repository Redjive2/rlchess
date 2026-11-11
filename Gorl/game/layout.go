package game

import rl "github.com/gen2brain/raylib-go/raylib"

// layout describes where a grid of squares sits on screen, so the main board
// and smaller faux boards (like the promotion picker) share the same
// mouse-to-square and square-to-screen logic. Rank 0 is the bottom row.
type layout struct {
	x, y         int32 // top-left corner of the squares (inside the border)
	ranks, files int32
	squareSize   int32
}

// mainBoard is the 8x8 playing board, centered in the window.
func mainBoard() layout {
	return centered(8, 8,
		(int32(rl.GetScreenWidth())-squareSize*8)/2,
		(int32(rl.GetScreenHeight())-squareSize*8)/2,
	)
}

// promotionBoard is the 1x4 promotion picker, centered horizontally in the
// gap above the main board.
func promotionBoard() layout {
	return centered(1, 4,
		(int32(rl.GetScreenWidth())-squareSize*4)/2,
		(mainBoard().y-margin-squareSize)/2,
	)
}

func centered(ranks, files, x, y int32) layout {
	return layout{x: x, y: y, ranks: ranks, files: files, squareSize: squareSize}
}

// width and height are the size of the squares, not counting the border.
func (l layout) width() int32  { return l.files * l.squareSize }
func (l layout) height() int32 { return l.ranks * l.squareSize }

// row is s's row on screen, counted from the top.
func (l layout) row(s square) int32 { return l.ranks - 1 - s.Rank() }

// column is s's column on screen, counted from the left.
func (l layout) column(s square) int32 { return s.File() }

// vector returns the screen position of s's top-left corner.
func (l layout) vector(s square) rl.Vector2 {
	return rl.NewVector2(
		float32(l.x+l.column(s)*l.squareSize),
		float32(l.y+l.row(s)*l.squareSize),
	)
}

// mousePos returns the mouse's screen position. It's the same for every
// layout; pass it to squareAt to find the square under it.
func (l layout) mousePos() rl.Vector2 {
	return rl.GetMousePosition()
}

// findMouseTarget returns the square under the mouse, and whether the mouse
// is over l at all.
func (l layout) findMouseTarget() (square, bool) {
	return l.squareAt(l.mousePos())
}

// squareAt returns the square under screen position pos, and whether pos is
// inside l.
func (l layout) squareAt(pos rl.Vector2) (square, bool) {
	var (
		mx, my = int32(pos.X), int32(pos.Y)

		outOfRangeX = mx < l.x || mx >= l.x+l.width()
		outOfRangeY = my < l.y || my >= l.y+l.height()
	)

	if outOfRangeX || outOfRangeY {
		return square{}, false
	}

	var (
		column = (mx - l.x) / l.squareSize
		row    = (my - l.y) / l.squareSize
	)

	return square{l.ranks - 1 - row, column}, true
}

// drawBorder draws the black border around l's squares.
func (l layout) drawBorder() {
	rl.DrawRectangle(l.x-margin, l.y-margin, l.width()+margin*2, l.height()+margin*2, rl.Black)
}
