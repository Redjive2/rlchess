package game

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// DrawBoard draws the border and squares, highlighting the square under the
// mouse and the selected square.
func (b *board) DrawBoard(moveOptions []square) {
	var (
		l                    = mainBoard()
		target, isTargeting  = l.findMouseTarget()
		selected, isSelected = b.GetSelection()
	)

	l.drawBorder()

	kingSquare, _ := b.find(pieceKindKing, turn)

	for sq := range allSquares() {
		color := rl.NewColor(235, 235, 235, 255) // white by default

		if sq.Rank()%2 == sq.File()%2 { // dark when rank+file is even (A1 is dark)
			color = rl.NewColor(60, 55, 50, 255)
		}

		switch {
		case isTargeting && sq == target:
			if color.R == 235 {
				color = rl.NewColor(250, 200, 200, 205)
			} else {
				color = rl.NewColor(100, 60, 60, 170)
			}
		case slices.Contains(moveOptions, sq):
			if color.R == 235 {
				color = rl.NewColor(215, 255, 220, 200)
			} else {
				color = rl.NewColor(60, 90, 60, 170)
			}
		case isSelected && sq == selected:
			if color.R == 235 {
				color = rl.NewColor(255, 200, 200, 215)
			} else {
				color = rl.NewColor(105, 60, 60, 180)
			}
		case b.inCheck(turn) && sq == kingSquare:
			if color.R == 235 {
				color = rl.NewColor(255, 118, 118, 255)
			} else {
				color = rl.NewColor(208, 108, 108, 215)
			}
		}

		// files run left to right, ranks run bottom (ONE) to top (EIGHT)
		corner := l.vector(sq)
		rl.DrawRectangle(int32(corner.X), int32(corner.Y), squareSize, squareSize, color)
	}
}

func (b *board) DrawPieces() {
	for s := range allSquares() {
		if p := b.pieceAt(s); p != nil {
			p.Draw(mainBoard().vector(s), squareSize)
		}
	}
}

func (p piece) Draw(pos rl.Vector2, squareSize int32) {
	texPath := "assets/textures/"

	switch p.color {
	case colorKindWhite:
		texPath += "white"
	case colorKindBlack:
		texPath += "black"
	}

	texPath += "/"

	switch p.kind {
	case pieceKindPawn:
		texPath += "pawn"
	case pieceKindRook:
		texPath += "rook"
	case pieceKindKnight:
		texPath += "knight"
	case pieceKindBishop:
		texPath += "bishop"
	case pieceKindQueen:
		texPath += "queen"
	case pieceKindKing:
		texPath += "king"
	}

	texPath += ".png"

	var (
		tex  = loadTexture(texPath)
		size = float32(squareSize)

		// fit the texture inside the square, keeping its aspect ratio
		scale = min(size/float32(tex.Width), size/float32(tex.Height))

		// center it in the square
		position = rl.NewVector2(
			pos.X+(size-float32(tex.Width)*scale)/2,
			pos.Y+(size-float32(tex.Height)*scale)/2,
		)
	)

	rl.DrawTextureEx(tex, position, 0, scale, rl.White)
}

var textureCache = map[string]rl.Texture2D{}

func loadTexture(path string) rl.Texture2D {
	if tex, found := textureCache[path]; found {
		return tex
	}

	var (
		image = rl.LoadImage(path)
		tex   = rl.LoadTextureFromImage(image)
	)

	rl.UnloadImage(image)

	textureCache[path] = tex
	return tex
}
