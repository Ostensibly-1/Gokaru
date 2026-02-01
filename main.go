package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
)

func main() {
	App := app.New()
	LoadGame(&App, false)
	fmt.Println("All systems were operational")
}
