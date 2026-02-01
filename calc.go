package main

import (
	"fyne.io/fyne/v2"
)

// dictionary
var ROOK_DIRS = [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
var BISHOP_DIRS = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var QUEEN_DIRS = [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var KNIGHT_DIRS = [][2]int{{1, 2}, {1, -2}, {-1, 2}, {-1, -2}, {2, 1}, {2, -1}, {-2, 1}, {-2, -1}}
var KING_DIRS = [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

// actually calculating moves (wtf)
func CalculatePieces(board *Board, render_board *ChessRenderer, app *fyne.App) {
	for x := range 8 {
		for y := range 8 {
			if board.squares[x][y].piece != nil {
				piece := board.squares[x][y].piece

				raw := GetPseudoMoves(board, piece, x, y)

				// if the move causes the current turn's king to be checked, then remove that
				movable_sqrs := []*Square{}
				for _, target := range raw {
					if SimulateMoveAndCheckSafety(board, board.squares[x][y], target, piece.team) {
						movable_sqrs = append(movable_sqrs, target)
					}
				}

				piece.available_squares = &movable_sqrs
			}
		}
	}
}

// for king safety
func SimulateMoveAndCheckSafety(board *Board, prev_sqr, next_sqr *Square, team TeamType) bool {
	captured := next_sqr.piece
	tomove := prev_sqr.piece

	next_sqr.piece = tomove
	prev_sqr.piece = nil

	safe := !IsKingInCheck(board, team)

	prev_sqr.piece = tomove
	next_sqr.piece = captured

	return safe
}

// coroutine to check if king is checked
func IsKingInCheck(board *Board, team TeamType) bool {
	var king_x, king_y int
	found := false

	for x := range 8 {
		for y := range 8 {
			p := board.squares[x][y].piece
			if p != nil && p._type == King && p.team == team {
				king_x = x
				king_y = y
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	// see if enemies could attack
	for x := range 8 {
		for y := range 8 {
			p := board.squares[x][y].piece
			if p != nil && p.team != team {
				enemy_sqrs := GetPseudoMoves(board, p, x, y)
				for _, mv := range enemy_sqrs {
					if mv.x == king_x && mv.y == king_y {
						return true
					}
				}
			}
		}
	}

	return false
}

// like CalculatePieces but limited to chess movement rules only
func GetPseudoMoves(board *Board, piece *Piece, x, y int) []*Square {
	moves := []*Square{}

	switch piece._type {
	case Pawn:
		direction := 1
		start := BLACK_START_LANE
		if piece.team == White {
			direction = -1
			start = WHITE_START_LANE
		}

		next := x + direction
		if IsValid(next, y) {
			if board.squares[next][y].piece == nil {
				moves = append(moves, board.squares[next][y])
				// if at start
				forward_twice := x + (direction * 2)
				if x == start && IsValid(forward_twice, y) {
					if board.squares[forward_twice][y].piece == nil {
						moves = append(moves, board.squares[forward_twice][y])
					}
				}
			}
		}

		// if enemy at front next corners
		for _, dy := range []int{-1, 1} {
			if IsValid(next, y+dy) {
				target := board.squares[next][y+dy]
				if target.piece != nil && target.piece.team != piece.team {
					moves = append(moves, target)
				}
			}
		}

	case Rook:
		moves = append(moves, GetMovement1(board, x, y, piece.team, ROOK_DIRS)...)
	case Bishop:
		moves = append(moves, GetMovement1(board, x, y, piece.team, BISHOP_DIRS)...)
	case Queen:
		moves = append(moves, GetMovement1(board, x, y, piece.team, QUEEN_DIRS)...)
	case Knight:
		moves = append(moves, GetMovement2(board, x, y, piece.team, KNIGHT_DIRS)...)
	case King:
		moves = append(moves, GetMovement2(board, x, y, piece.team, KING_DIRS)...)
	}

	return moves
}

// for pieces that cant skip a piece every move
func GetMovement1(board *Board, x, y int, team TeamType, dirs [][2]int) []*Square {
	res := []*Square{}
	for _, d := range dirs {
		dx, dy := d[0], d[1]
		cx, cy := x, y
		for {
			cx += dx
			cy += dy
			if !IsValid(cx, cy) {
				break
			}

			target := board.squares[cx][cy]
			if target.piece == nil {
				res = append(res, target)
			} else {
				if target.piece.team != team {
					res = append(res, target) // Capture
				}
				break // Blocked
			}
		}
	}
	return res
}

// for pieeces that can skip a piece when moving or capture it
func GetMovement2(board *Board, x, y int, team TeamType, offsets [][2]int) []*Square {
	res := []*Square{}
	for _, off := range offsets {
		nx, ny := x+off[0], y+off[1]
		if IsValid(nx, ny) {
			target := board.squares[nx][ny]
			if target.piece == nil || target.piece.team != team {
				res = append(res, target)
			}
		}
	}
	return res
}

// boundary check
func IsValid(x, y int) bool {
	return x >= 0 && x < 8 && y >= 0 && y < 8
}
