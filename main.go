package main

import (
	"fmt"
	"image/color"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kpechenenko/rword"
	"golang.org/x/image/font/basicfont"
)

type Game struct{
	bg *ebiten.Image
	start time.Time

	error_message string

	backspacePrev bool
	enterKeyPrev bool
	game_won int // 1 = win; -1 = loss

	//constants
	rows int
	wordLength int
	width float32
	height float32

	unknown_word string
	word_letters []string

	guessed_word string
	word_before string

	tried_words []string
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
	}
	g.backspacePrev = back_space_pressed

	enter_key_pressed := ebiten.IsKeyPressed(ebiten.KeyEnter)
	if enter_key_pressed && !g.enterKeyPrev {
		if !slices.Contains(g.tried_words, g.guessed_word) && len(g.guessed_word) == g.wordLength {
			g.tried_words = append(g.tried_words, g.guessed_word)
			g.word_before = g.guessed_word
			g.guessed_word = ""
			if g.error_message != "" {
				g.error_message = ""
			}
		} else {
			g.raise_error()
		}
	} 
	g.enterKeyPrev = enter_key_pressed

	// TODO: maybe method for logic regarding win/loss window 
	if g.word_before == g.unknown_word || len(g.tried_words) == g.rows { 
		if g.start.IsZero() {
			g.start = time.Now()
			if g.word_before == g.unknown_word { //Game result logic
				g.game_won = 1
			} else {
				g.game_won = -1
			}
		}
		if time.Since(g.start) > 3 * time.Second {
			return ebiten.Termination
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{242, 241, 144, 255})
	if g.bg != nil {
		screen.DrawImage(g.bg, nil)
	}

	ebitenutil.DebugPrint(screen, "herr sergeant! - swedish soldier")	
	text.Draw(screen, "Wordle!", basicfont.Face7x13, 50, 70, color.Black)
	text.Draw(screen, g.guessed_word, basicfont.Face7x13, 100, 90, color.Black)
	text.Draw(screen, g.error_message, basicfont.Face7x13, 200, 90, color.RGBA{255, 0, 0, 255})

	for i := 0; i < len(g.tried_words); i++ {
		if g.tried_words[i] != "" {
			for j := 0; j < g.wordLength; j++ {
				c := g.letter_exist(string(g.tried_words[i][j]), j)
				ebitenutil.DrawRect(screen, float64(100 + 80 * j), float64(100 + 80 * i), 70, 70, color.RGBA{c[0], c[1], c[2], c[3]})
				text.Draw(screen, string(g.tried_words[i][j]), basicfont.Face7x13, 130 + 80 * j, 130 + 80 * i, color.Black)
			}
		}
	}

	if g.game_won != 0 {
		ebitenutil.DrawRect(screen, 51, 51, 1000-102, 650-102, color.Black)
		ebitenutil.DrawRect(screen, 50, 50, 1000-100, 650-100, color.White)
		switch g.game_won {
		case 1:
			text.Draw(screen, "You won!", basicfont.Face7x13, 1000/2, 700/2, color.Black)
		case -1:
			text.Draw(screen, "You lost!", basicfont.Face7x13, 1000/2, 700/2, color.Black)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) initBG(xPos, yPos float32) {
	g.bg = ebiten.NewImage(1000, 700)
	g.bg.Fill(color.RGBA{242, 241, 144, 255}) // or draw shapes onto bg

	for i := 0; i < g.rows; i++ {
		var y_pos_step float32 = yPos + (g.height + 10) * float32(i)
		for j := 0; j < g.wordLength; j++ {
			var x_pos_step float32 = xPos + (g.width + 10) * float32(j)
			vector.FillRect(g.bg, x_pos_step - 1, y_pos_step - 1, g.width + 2, g.width + 2, color.Black, false)
			vector.FillRect(g.bg, x_pos_step, y_pos_step, g.width, 70, color.White, false)
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

func (g *Game) raise_error() {
	g.error_message = "Enter a word that you haven't in format ABCDE"
}

func (g *Game) letter_exist(char string, index int) [4]uint8 {
	if g.word_letters[index] == char {
		color_code := [4]uint8{0, 255, 0, 255}
		return color_code
	} else if slices.Contains(g.word_letters, char) {
		color_code := [4]uint8{204, 102, 0, 255}
		return color_code
	} else {
		color_code := [4]uint8{128, 128, 128, 255}
		return color_code
	}
}

func main() {
	ebiten.SetWindowSize(1000, 700)
	ebiten.SetWindowTitle("wordle")	
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := &Game{wordLength: 5, rows: 6, width: 70, height: 70}
	g.initBG(100, 100)

	g.unknown_word = random_word()
	g.word_letters = strings.Split(g.unknown_word, "")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
