package game

func (b *board) canMove(c colorKind) bool {
	for r, rank := range b.Ranks() {
		for f, piecePtr := range rank {
			var (
				piece = *piecePtr
				sq = square{int32(r), int32(f)}
			)

			if piece == nil || piece.color != c {
				continue
			}

			if len(piece.GetMoveOptions(b, sq)) != 0 {
				return true
			}
		}
	}

	return false
}

func (b *board) GameStatus() (gameStatus, colorKind) {
	var (
		whiteCanMove, blackCanMove = b.canMove(colorKindWhite), b.canMove(colorKindBlack)
		whiteInCheck, blackInCheck = b.inCheck(colorKindWhite), b.inCheck(colorKindBlack)
	)

	// the second result is the winner: the side that delivered mate
	if whiteInCheck && !whiteCanMove {
		return gameStatusCheckmate, colorKindBlack
	}

	if blackInCheck && !blackCanMove {
		return gameStatusCheckmate, colorKindWhite
	}

	var (
		whiteInStalemate = turn == colorKindWhite && !whiteCanMove
		blackInStalemate = turn == colorKindBlack && !blackCanMove
	)
	if whiteInStalemate || blackInStalemate {
		return gameStatusStalemate, 2
	}

	return gameStatusPlayable, 2
}
