package main

import (
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 640
	screenHeight = 640
)

var mainFont font.Face

func reload() {
	fileContent, err := os.ReadFile("gameTest.sc")
	if err != nil {
		log.Fatalln(err)
	}
	tokens := Tokenize(string(fileContent))

	ast := ParseProgram(&tokens)

	ast.Interpret(nil)

	start, ok := functions["start"]

	if !ok {
		log.Fatalln("ERROR: \"Start\" is not defined!")
	}

	args := []Expression{
		ParseLiteral(
			Token{Type: TOK_NUM, Lexme: strconv.Itoa(screenWidth)},
		),
		ParseLiteral(
			Token{Type: TOK_NUM, Lexme: strconv.Itoa(screenHeight)},
		),
	}

	start.Call(args, nil)
}

type Game struct{}

func (g *Game) Update() error {
	update, _ := functions["update"]

	//if !ok {
	//	log.Fatalln("ERROR: Game Loop not defined!")
	//}

	jumpKeyPressed := 0
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		jumpKeyPressed = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		reload()
	}

	args := []Expression{
		ParseLiteral(
			Token{Type: TOK_NUM, Lexme: strconv.Itoa(jumpKeyPressed)},
		),
	}

	update.Call(args, nil)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{77, 240, 255, 255})

	playerX := float32(mainScope.Variables["playerX"].Value)
	playerY := float32(mainScope.Variables["playerY"].Value)
	playerW := float32(mainScope.Variables["playerWidth"].Value)
	playerH := float32(mainScope.Variables["playerHeight"].Value)

	vector.DrawFilledRect(screen, playerX, playerY, playerW, playerH, color.RGBA{255, 63, 46, 255}, true)

	obstacleX := float32(mainScope.Variables["obstacleX"].Value)
	obstacleY := float32(mainScope.Variables["obstacleY"].Value)
	obstacleW := float32(mainScope.Variables["obstacleWidth"].Value)
	obstacleH := float32(mainScope.Variables["obstacleHeight"].Value)

	vector.DrawFilledRect(screen, obstacleX, obstacleY, obstacleW, obstacleH, color.RGBA{255, 0, 0, 255}, true)

	score := mainScope.Variables["score"].Value

	text.Draw(screen, strconv.Itoa(score), mainFont, screenWidth/2, screenHeight/8, color.RGBA{255, 204, 64, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Game test
func main() {
	reload()

	// Load Font
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mainFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingVertical,
	})

	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game Test")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

/*
Old regular code for SEGC ( just comented for this test lmfao )
func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalln("ERROR: No input file being passed!")
	}
	fileContent, err := os.ReadFile(args[1])
	if err != nil {
		log.Fatalln(err)
	}

	tokens := Tokenize(string(fileContent))

	ast := ParseProgram(&tokens)

	ast.Interpret(nil)

	main, ok := functions["main"]

	if !ok {
		log.Fatalln("ERROR: Entry point not defined!")
	}

	var exprs []Expression

	if len(args)-2 != len(main.Args) {
		log.Println("ERROR: Number of Command-line arguments not matching!")
		log.Fatalln("Expected: ", len(main.Args))
		log.Fatalln("Recived: ", len(args)-2)
	}

	for i := 2; i < len(args); i++ {
		exprs = append(exprs, ExpressionLiteral{Type: num_literal, Tok: Token{Type: TOK_NUM, Lexme: args[i]}})
	}

	errCode := main.Call(exprs, nil)
	os.Exit(errCode)
}
*/
