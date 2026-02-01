package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LoadGame(App *fyne.App, game_over bool) {
	MainWindow := (*App).NewWindow("Gokaru by Ostensibly")
	MainWindow.SetFixedSize(true)

	if game_over == true {
		GameOverLabel := widget.NewLabel("Game Over")
		RetryBtn := widget.NewButton("Retry", func() {
			(*App).Quit()
			main()
		})

		GameOverScreen := container.NewVBox(GameOverLabel, RetryBtn)

		MainWindow.SetContent(GameOverScreen)
		MainWindow.ShowAndRun()
	}

	Grid := container.NewAdaptiveGrid(8)

	// setup Chessboard
	var RenderChessboard ChessRenderer
	BoardStruct := Board{turn: White}

	for x := range 8 {
		for y := range 8 {
			// generate Squares
			Rect := canvas.NewRectangle(color.White)
			Rect.SetMinSize(fyne.NewSize(50, 50))

			if (x+y)%2 != 0 {
				Rect.FillColor = color.Black
			}

			var NewSquare *Square
			NewSquare = CreateChessSquare(x, y, Rect, func(x, y int) {
				if BoardStruct.ready_to_move == true {
					// chess piece move logic
					AvailSqrs := BoardStruct.squares[BoardStruct.selected_square[0]][BoardStruct.selected_square[1]].piece.available_squares
					for idx := range len(*AvailSqrs) {
						AvlSqr := (*AvailSqrs)[idx]

						if x == AvlSqr.x && y == AvlSqr.y {
							if BoardStruct.selected_square != nil {
								OldX := BoardStruct.selected_square[0]
								OldY := BoardStruct.selected_square[1]
								OldSqr := BoardStruct.squares[OldX][OldY]

								BoardStruct.MovePiece(OldSqr, NewSquare)

								// turn switching
								if BoardStruct.turn == Black {
									BoardStruct.turn = White
								} else {
									BoardStruct.turn = Black
								}
							}
						}
					}

					RefreshPieces(&BoardStruct, &RenderChessboard)
					RenderChessboard()

					BoardStruct.ready_to_move = false
				} else {
					HandlePieces(x, y, &BoardStruct, &RenderChessboard)
				}
			})

			// generate pieces
			if x == WHITE_START_LANE {
				NewSquare.piece = &Piece{
					_type: Pawn,
					team:  White,
				}
			} else if x == BLACK_START_LANE {
				NewSquare.piece = &Piece{
					_type: Pawn,
					team:  Black,
				}
			}

			if x == WHITE_END_LANE {
				if y == 0 || y == 7 {
					NewSquare.piece = &Piece{
						_type: Rook,
						team:  White,
					}
				} else if y == 1 || y == 6 {
					NewSquare.piece = &Piece{
						_type: Knight,
						team:  White,
					}
				} else if y == 2 || y == 5 {
					NewSquare.piece = &Piece{
						_type: Bishop,
						team:  White,
					}
				} else if y == 3 {
					NewSquare.piece = &Piece{
						_type: Queen,
						team:  White,
					}
				} else {
					NewSquare.piece = &Piece{
						_type: King,
						team:  White,
					}
				}
			} else if x == BLACK_END_LANE {
				if y == 0 || y == 7 {
					NewSquare.piece = &Piece{
						_type: Rook,
						team:  Black,
					}
				} else if y == 1 || y == 6 {
					NewSquare.piece = &Piece{
						_type: Knight,
						team:  Black,
					}
				} else if y == 2 || y == 5 {
					NewSquare.piece = &Piece{
						_type: Bishop,
						team:  Black,
					}
				} else if y == 3 {
					NewSquare.piece = &Piece{
						_type: Queen,
						team:  Black,
					}
				} else {
					NewSquare.piece = &Piece{
						_type: King,
						team:  Black,
					}
				}
			}

			BoardStruct.squares[x][y] = NewSquare
		}
	}

	RenderChessboard = func() {
		Grid.Objects = nil
		CalculatePieces(&BoardStruct, &RenderChessboard, App)
		for x := range 8 {
			for y := range 8 {
				// render Squares
				RawSquare := BoardStruct.squares[x][y]

				// render Pieces
				if RawSquare.piece != nil {
					PieceType := RawSquare.piece._type
					PieceTeam := RawSquare.piece.team
					var PieceRender fyne.CanvasObject

					if PieceType == Pawn {
						PieceRender = canvas.NewCircle(color.White)

						if PieceTeam == Black {
							if Circ, ok := PieceRender.(*canvas.Circle); ok {
								Circ.FillColor = color.Black
							}
						}

						if Circ, ok := PieceRender.(*canvas.Circle); ok {
							Circ.StrokeColor = color.NRGBA{R: 255, G: 10, B: 10, A: 255}
							Circ.StrokeWidth = 2
						}

						NewSquare := container.NewStack(RawSquare, PieceRender)
						Grid.Add(NewSquare)

						continue
					} else if PieceType == Rook {
						PieceRender = canvas.NewPolygon(4, color.White)
					} else if PieceType == Knight {
						PieceRender = canvas.NewPolygon(3, color.White)
					} else if PieceType == Bishop {
						PieceRender = canvas.NewPolygon(5, color.White)
					} else if PieceType == King {
						PieceRender = canvas.NewPolygon(6, color.White)
					} else if PieceType == Queen {
						PieceRender = canvas.NewPolygon(12, color.White)
					}

					if PieceTeam == Black {
						if Poly, ok := PieceRender.(*canvas.Polygon); ok {
							Poly.FillColor = color.Black
						}
					}

					if Poly, ok := PieceRender.(*canvas.Polygon); ok {
						Poly.StrokeColor = color.NRGBA{R: 255, G: 10, B: 10, A: 255}
						Poly.StrokeWidth = 2
					}

					NewSquare := container.NewStack(RawSquare, PieceRender)
					Grid.Add(NewSquare)

					continue
				}

				// render selected
				if BoardStruct.selected_square != nil {
					Selected := BoardStruct.selected_square

					if BoardStruct.squares[Selected[0]][Selected[1]].piece != nil {
						Content := BoardStruct.squares[Selected[0]][Selected[1]].content
						if Sqr, ok := Content.(*canvas.Rectangle); ok {
							Sqr.FillColor = color.NRGBA{R: 255, G: 255, B: 0, A: 255}
						}
					}
				}

				Grid.Add(RawSquare)
			}
		}
	}

	// render board + credits
	RenderChessboard()
	Credits := widget.NewLabel("Gokaru Version 0.1.0 By Ostensibly")
	Container := container.NewVBox(Grid, Credits)

	// start
	MainWindow.SetContent(Container)
	MainWindow.ShowAndRun()
}
