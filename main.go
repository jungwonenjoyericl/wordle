package main

import (
	"image/color"
	//"image/draw"
	"log"


	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{
	bg *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{242, 241, 144, 255})
	if g.bg != nil {
		screen.DrawImage(g.bg, nil)
	}
	const font_size = 32
	text.Draw(screen, "Wordle!", basicfont.Face7x13, 50, 70, color.Black)
	ebitenutil.DebugPrint(screen, "herr sergeant! - swedish soldier")

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) initBG(xPos, yPos float32) {
	g.bg = ebiten.NewImage(1000, 700)
	g.bg.Fill(color.RGBA{242, 241, 144, 255}) // or draw shapes onto bg

	const width, height float32 = 70, 70
	const rows, wordlength int = 6, 5
	for i := 0; i < rows; i++ {
		var y_pos_step float32 = yPos + (height + 10) * float32(i)
		for j := 0; j < wordlength; j++ {
			var x_pos_step float32 = xPos + (width + 10) * float32(j)
			vector.FillRect(g.bg, x_pos_step - 1, y_pos_step - 1, width + 2, width + 2, color.Black, false)
			vector.FillRect(g.bg, x_pos_step, y_pos_step, width, 70, color.White, false)
		}
	}
}

func main() {
	ebiten.SetWindowSize(1000, 700)
	ebiten.SetWindowTitle("wordle")	
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &Game{}
	g.initBG(100, 100)
	
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
