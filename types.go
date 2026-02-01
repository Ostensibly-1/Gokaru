package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// TEAMS
type TeamType int

const (
	White TeamType = iota
	Black
)

// PIECES
type PieceType int

const (
	Pawn PieceType = iota
	Knight
	Rook
	Bishop
	Queen
	King
)

type Piece struct {
	_type             PieceType
	team              TeamType
	available_squares *[]*Square
}

// SQUARES
type Square struct {
	widget.BaseWidget
	x        int
	y        int
	selected bool
	piece    *Piece
	content  fyne.CanvasObject
	on_tap   func(x, y int)
}

func (self *Square) Tapped(*fyne.PointEvent) {
	self.on_tap(self.x, self.y)
}

func (self *Square) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(self.content)
}

func (self *Square) MouseIn(*desktop.MouseEvent) {
	if result, ok := self.content.(*canvas.Rectangle); ok {
		if (self.piece != nil) && (self.selected != true) {
			result.FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
			self.Refresh()
		}
	} else {
		log.Fatal("Canvas object is not a rectangle")
	}
}

func (self *Square) MouseOut() {
	if result, ok := self.content.(*canvas.Rectangle); ok {
		if self.selected != true && self.piece != nil {
			if (self.x+self.y)%2 != 0 {
				result.FillColor = color.Black
			} else {
				result.FillColor = color.White
			}
			self.Refresh()
		}
	}
}

func (self *Square) MouseMoved(*desktop.MouseEvent) {
}

func CreateChessSquare(x, y int, content fyne.CanvasObject, on_tap func(x, y int)) *Square {
	NewSquare := &Square{
		x:       x,
		y:       y,
		content: content,
		on_tap:  on_tap,
	}
	NewSquare.ExtendBaseWidget(NewSquare)
	return NewSquare
}

// CHESSBOARD
type Board struct {
	squares         [8][8]*Square
	turn            TeamType
	selected_square *[2]int // x,y
	ready_to_move   bool
	// available_squares *[]*Square
}

func (self *Board) MovePiece(square1, square2 *Square) {
	square2.piece = square1.piece
	square1.piece = nil
}

// misc
type ChessRenderer func()
