package main

import (
	"image/color"
	"log"
	"fmt"
	"strings"
	
	"github.com/kpechenenko/rword"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct{
	bg *ebiten.Image

	backspacePrev bool

	rows int
	wordLength int

	unknown_word string
	word_letters []string

	guessed_word string
	guess_word_letters []string
}

func (g *Game) Update() error {
	for _,r := range ebiten.InputChars() {
		if r >= 'a' && r <= 'z' && len(g.guessed_word) < g.wordLength {
			g.guessed_word += string(r)
		}
	}
	
	back_space_pressed := ebiten.IsKeyPressed(ebiten.KeyBackspace) // remove 1 char from guess per press
	if back_space_pressed && !g.backspacePrev && len(g.guessed_word) > 0 {
		g.guessed_word = g.guessed_word[:len(g.guessed_word) - 1]
		fmt.Println(g.guessed_word)
	}
	g.backspacePrev = back_space_pressed

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.split_into_chars(g.guessed_word)
	}

	
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{242, 241, 144, 255})
	if g.bg != nil {
		screen.DrawImage(g.bg, nil)
	}

	ebitenutil.DebugPrint(screen, "herr sergeant! - swedish soldier")
	text.Draw(screen, g.guessed_word, basicfont.Face7x13, 100, 90, color.Black)
	text.Draw(screen, "Wordle!", basicfont.Face7x13, 50, 70, color.Black)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) initBG(xPos, yPos float32) {
	g.bg = ebiten.NewImage(1000, 700)
	g.bg.Fill(color.RGBA{242, 241, 144, 255}) // or draw shapes onto bg

	const width, height float32 = 70, 70
	for i := 0; i < g.rows; i++ {
		var y_pos_step float32 = yPos + (height + 10) * float32(i)
		for j := 0; j < g.wordLength; j++ {
			var x_pos_step float32 = xPos + (width + 10) * float32(j)
			vector.FillRect(g.bg, x_pos_step - 1, y_pos_step - 1, width + 2, width + 2, color.Black, false)
			vector.FillRect(g.bg, x_pos_step, y_pos_step, width, 70, color.White, false)
		}
	}
}

func random_word() string {
	var g rword.GenerateRandom
	var err error
	g, err = rword.New()
	if err != nil {
		panic(err)
	}
	var word string
	for {
		word = g.Word()
		if len(word) == 5{
			fmt.Println(word)
			return word
		}
	}
}

func (g *Game) split_into_chars(guess string) {
	g.guess_word_letters = strings.Split(guess, "")
}

func main() {
	ebiten.SetWindowSize(1000, 700)
	ebiten.SetWindowTitle("wordle")	
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &Game{wordLength: 5, rows: 6}
	g.initBG(100, 100)

	g.unknown_word = random_word()
	g.word_letters = strings.Split(g.unknown_word, "")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
