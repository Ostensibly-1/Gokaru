package main

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

// colors the movable zones every piece selection of the team that has the turn
func HandlePieces(x, y int, board *Board, render_board *ChessRenderer) {
	RefreshPieces(board, render_board)

	SelectedSquare := board.squares[x][y]

	if SelectedSquare.piece != nil && SelectedSquare.piece.team == board.turn {

		board.selected_square = &[2]int{x, y}
		board.squares[x][y].selected = true
		(*render_board)()

		ChessPiece := board.squares[x][y].piece
		AllowedSquares := ChessPiece.available_squares

		// color green if pawn can move/capture on that square
		for idx := range len(*AllowedSquares) {
			AllowedSqr := (*AllowedSquares)[idx]
			if res, ok := AllowedSqr.content.(*canvas.Rectangle); ok {
				res.FillColor = color.NRGBA{R: 0, G: 250, B: 0, A: 255}
			}
		}

		board.ready_to_move = true
	}
}

// refreshes everything.. DUH
func RefreshPieces(board *Board, render_board *ChessRenderer) {
	for x := range 8 {
		for y := range 8 {
			board.squares[x][y].selected = false

			if res, ok := (board.squares[x][y].content).(*canvas.Rectangle); ok {
				res.FillColor = color.White
			}

			if (x+y)%2 != 0 {
				if res, ok := (board.squares[x][y].content).(*canvas.Rectangle); ok {
					res.FillColor = color.Black
				}
			}
		}
	}

	board.selected_square = nil

	(*render_board)()
}
