package game

import (
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	// configuration
	squareSize int32 = 65
	margin     int32 = 4

	// starting turn color
	turn = colorKindWhite

	b *board

	inPromotion     bool
	promotionSquare square // where the promoting pawn landed
	status          = gameStatusPlayable
	winningColor    colorKind
)

func Load() {
	b = SetupBoard()
}

func Unload() {
	for _, tex := range textureCache {
		rl.UnloadTexture(tex)
	}
}

func Loop(frameTime int) {
	// snapshot which screen is showing, so a click that ends the game or starts
	// a promotion isn't also handled by the overlay it just opened
	var (
		gameOver  = status != gameStatusPlayable
		promoting = inPromotion
	)

	if !gameOver && !promoting && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		handleBoardClick()
	}

	moveOptions := selectedMoveOptions()

	rl.ClearBackground(rl.RayWhite)
	b.DrawBoard(moveOptions)
	b.DrawPieces()

	// overlays are drawn after the board so they sit on top of it
	switch {
	case gameOver:
		if drawEnd() {
			b = SetupBoard()
			status = gameStatusPlayable
			winningColor = 2
			turn = colorKindWhite
		}
	case promoting:
		if promotedPiece, ok := drawPromotion(); ok {
			promotedPiece.HasMoved = true // a promoted rook can't castle
			*b.At(promotionSquare) = &promotedPiece
			inPromotion = false
			endTurn()
		}
	}
}

// selectedMoveOptions returns the move options of the selected piece, if it's
// the side to move's.
func selectedMoveOptions() []square {
	selected, ok := b.GetSelection()
	if !ok {
		return nil
	}

	piece := b.pieceAt(selected)
	if piece == nil || piece.color != turn {
		return nil
	}

	return piece.GetMoveOptions(b, selected)
}

// handleBoardClick moves the selected piece if a highlighted square was
// clicked; any other click selects (or deselects) instead.
func handleBoardClick() {
	target, onBoard := mainBoard().findMouseTarget()
	selected, isSelected := b.GetSelection()

	if !isSelected || !onBoard || !slices.Contains(selectedMoveOptions(), target) {
		if onBoard {
			b.Select(target)
		} else {
			b.Deselect()
		}
		return
	}

	var (
		piece            = b.pieceAt(selected)
		isPawn           = piece.kind == pieceKindPawn
		isWhitePromotion = turn == colorKindWhite && target.Rank() == EIGHT
		isBlackPromotion = turn == colorKindBlack && target.Rank() == ONE
	)

	b.Move(selected, target)

	if isPawn && (isWhitePromotion || isBlackPromotion) {
		inPromotion = true
		promotionSquare = target
		return
	}

	endTurn()
}

// endTurn passes the turn to the other side and checks for the end of the game.
func endTurn() {
	turn = turn.Opponent()
	b.Deselect()
	status, winningColor = b.GameStatus()
}

func drawEnd() bool {
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.NewColor(48, 38, 38, 64))

	message := "Game over! Stalemate."
	if status == gameStatusCheckmate {
		winner := "White"
		if winningColor == colorKindBlack {
			winner = "Black"
		}
		
		message = "Game over! " + winner + " wins."
	}

	var (
		spacing           float32 = 1.25
		fontSize                  = squareSize / 2
		messageDimensions         = rl.MeasureTextEx(rl.GetFontDefault(), message, float32(fontSize), spacing)
		centerX                   = int32(rl.GetScreenWidth()) / 2
		centerY                   = int32(rl.GetScreenHeight()) / 2
		startX                    = centerX - int32(messageDimensions.X/2)
		startY                    = centerY - int32(messageDimensions.Y/2)
		marginX           int32   = 6
		marginY           int32   = 4
	)

	rl.DrawRectangle(startX-marginX, startY-marginY, int32(messageDimensions.X)+marginX*2, int32(messageDimensions.Y)+marginY*2, rl.NewColor(40, 24, 24, 240))
	rl.DrawTextEx(rl.GetFontDefault(), message, rl.NewVector2(float32(startX), float32(startY)), float32(fontSize), spacing, rl.White)


	return rl.IsMouseButtonPressed(rl.MouseLeftButton)
}

func drawPromotion() (piece, bool) {
	var (
		pieceList = [4]piece{
			{color: turn, kind: pieceKindRook},
			{color: turn, kind: pieceKindKnight},
			{color: turn, kind: pieceKindBishop},
			{color: turn, kind: pieceKindQueen},
		}
		promotionLayout     = promotionBoard()
		target, isTargeting = promotionLayout.findMouseTarget()
	)

	promotionLayout.drawBorder()

	for i, piece := range pieceList {
		var (
			sq     = square{0, int32(i)}
			vector = promotionLayout.vector(sq)
			color  = rl.NewColor(235, 235, 235, 255) // white by default
		)

		if i%2 == 0 { // dark when rank+file is even (A1 is dark)
			color = rl.NewColor(60, 55, 50, 255)
		}

		if isTargeting && sq == target {
			if color.R == 235 {
				color = rl.NewColor(250, 200, 200, 205)
			} else {
				color = rl.NewColor(100, 60, 60, 170)
			}
		}

		rl.DrawRectangle(int32(vector.X), int32(vector.Y), squareSize, squareSize, color)
		piece.Draw(vector, squareSize)
	}

	if isTargeting && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		return pieceList[target.File()], true
	}

	return piece{}, false
}
