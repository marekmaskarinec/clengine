package cliw

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"image/png"

	"github.com/fatih/color"
	"github.com/nsf/termbox-go"
)

/* Tile is a part of a world. It holds all important info. Name can be used for storing data about it in files. */
type Tile struct {
	Name    string
	Tile    string
	Damage  int
	Color   string
	BgColor string
}

/* Vector2 but shorter, so you don't have to type :] */
type Ve2 struct {
	X int
	Y int
}

type Layer struct {
	World [][]Tile
	Pos   Ve2
}

type PixLayer struct {
	PixMap [][]string
	Pos    Ve2
}

/*Changes one specific tile in the world. Can append tiles, if the world is too small.*/
func EditTile(world [][]Tile, pos Ve2, t Tile) ([][]Tile, error) {
	w := DuplicateWorld(world)
	if pos.X < 0 || pos.Y < 0 {
		return w, errors.New("You entered value smaller than zero")
	} else {
		for len(w) <= pos.X {
			w = append(w, make([]Tile, 0))
		}
		for len(w[pos.X]) <= pos.Y {
			w[pos.X] = append(w[pos.X], Tile{})
		}
		w[pos.X][pos.Y] = t
		return w, nil
	}
}

/*func EditLine(world [][]Tile, from, to Ve2, t Tile) ([][]Tile, error){
	if from.X < 0 || from.Y < 0 || to.X < from.X || to.Y < from.Y {
		return nil, errors.New("Invalid number")
	} else {
		if to.X >= to.Y{

		}
	}
}*/

/*Gets size of terminal window*/
func GetSize() (Ve2) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("getting terminal size: %s\n", err)
	}

	outStr := strings.Split(out.String(), " ")
	wString := strings.Split(outStr[1], "\n")
	h, err := strconv.Atoi(outStr[0])
	w, err := strconv.Atoi(wString[0])
	if err != nil {
		fmt.Printf("getting terminal size: %s", err)
	}

	return Ve2{h, w}
}

/*Draws world on the center of the screen
can get additional blank margin for user input in ui*/
func DrawCentered(w [][]Tile, additionalRow bool) {
	s := GetSize()
	he, wi := s.X, s.Y
	var rows, colls, wwidth int
	var currentRow [][]Tile
	currentRow = append(currentRow, make([]Tile, 0))

	rows = (he - len(w)) / 2

	/*first row has to be the longest*/
	for i := range w[0] {
		wwidth += len(w[0][i].Tile)
	}
	colls = (wi - wwidth) / 2

	for i := 0; i < len(w); i++ {
		SetCursor(Ve2{rows+i, colls})
		currentRow[0] = w[i]
		DrawWorld(currentRow)
	}
	if additionalRow {
		fmt.Print(strings.Repeat(" ", colls))
	}
}

/*makes cut from the world*/
func CutWorld(world [][]Tile, from Ve2, to Ve2) ([][]Tile, error) {
	if from.X+to.X <= len(world) && from.Y+to.Y <= len(world[0]) {
		var toReturn [][]Tile
		var toAppend []Tile
		for i := 0; i < to.X; i++ {
			toReturn = append(toReturn, toAppend)
			for j := 0; j < to.Y; j++ {
				toReturn[i] = append(toReturn[i], world[i+from.X][j+from.Y])
			}
		}
		return toReturn, nil
	} else {
		return world, errors.New("cut: Out of boundaries")
	}
}

/*returns new Ve2*/
func V2(x, y int) Ve2 {
	return Ve2{X: x, Y: y}
}

/*Adds on Ve2 to another*/
func (v1 *Ve2) Add(v2 Ve2) {
	v1.X += v2.X
	v1.Y += v2.Y
}

/* Changes all tiles in a rectangular shape. Fills the inside of the rectangle. */
func EditWorld(world [][]Tile, from, to Ve2, tile Tile) ([][]Tile, error) {
	if from.X < 0 || from.Y < 0 || to.X < from.X || to.Y < from.Y {
		return world, errors.New("Invalid number")
	} else {
		for len(world) <= from.X+to.X {
			world = append(world, make([]Tile, 0))
		}
		for i := 0; i <= to.X; i++ {
			for len(world[from.X+i]) <= from.Y+to.Y-1 {
				//fmt.Println("x")
				world[from.X+i] = append(world[from.X+i], Tile{})
			}
		}

		for i := 0; i < to.X; i++ {
			for r := 0; r < to.Y; r++ {
				world[from.X+i][from.Y+r] = tile
			}
		}
		return world, nil
	}
}

/* Saves world to a text file. Don't use this! Doesn't support background color! */
func SaveWorld(world [][]Tile, path string) {
	var c Tile
	var toWrite string
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			c = world[i][j]
			toWrite += strconv.Itoa(i) + "\n" + strconv.Itoa(j) + "\n" + c.Name + "\n" + c.Tile + "\n" + strconv.Itoa(c.Damage) + "\n" + c.Color + "\n"
		}
	}
	ioutil.WriteFile(path, []byte(toWrite), 0644)
}

/* Saves world to JSON. Use this! */

func SaveWorldJSON(w [][]Tile, path string) {
	toSave, err := json.Marshal(w)
	if err == nil {
		ioutil.WriteFile(path, toSave, 0644)
	}
}

/* Duplicates a world */

func DuplicateWorld(w [][]Tile) [][]Tile {
	var nw [][]Tile

	for i := 0; i < len(w); i++ {
		nw = append(nw, make([]Tile, 0))
		for j := 0; j < len(w[i]); j++ {
			nw[i] = append(nw[i], w[i][j])
		}
	}
	return nw
}

/* Loads world from file. This shouldn't be used anymore, but if you have old worlds, you can use this to convert them. Doesn't support background color! */
func LoadWorld(path string) [][]Tile {
	var world [][]Tile
	text := []string{}
	var damage, x, y int

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	for i := 0; i < len(text)-5; i++ {
		damage, _ = strconv.Atoi(text[i+4])
		x, _ = strconv.Atoi(text[i])
		y, _ = strconv.Atoi(text[i+1])
		for len(world) <= x {
			world = append(world, make([]Tile, 0))
		}
		for len(world[x]) <= y {
			world[x] = append(world[x], Tile{})
		}
		world[x][y] = Tile{Name: text[i+2], Tile: text[i+3], Damage: damage, Color: text[i+5]}
		i += 5
	}
	return world
}

/* Loads world from JSON */

func LoadWorldJSON(path string) [][]Tile {
	var w [][]Tile

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &w)
	return w
}

/*prints out the world to terminal*/
func DrawWorldOld(world [][]Tile) {
	var c Tile
	palette := palette()
	bgPalette := bgPalette()
	col := color.New(color.FgWhite)
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[i]); j++ {
			c = world[i][j]
			if c.Color == "blnk" {
				fmt.Print(c.Tile)
			} else {
				if c.BgColor != "" {
					col = color.New(palette[c.Color], bgPalette[c.BgColor])
				} else {
					col = color.New(palette[c.Color])
				}
				col.Print(c.Tile)
			}
		}
		fmt.Println("")
	}
}

func DrawWorld(w [][]Tile) {
	var c Tile
	var tp string
	var prev Tile

	tp += GetColor(w[0][0].Color, false) + GetColor(w[0][0].BgColor, true)

	w[len(w)-1] = append(w[len(w)-1], Tile{Tile: "", Color: "ffffff", BgColor: "000000"})

	for i := range w {
		for j := range w[i] {
			c = w[i][j]
			if prev.Color != c.Color {
				tp += GetColor(c.Color, false)
			}
			if prev.BgColor != c.BgColor {
				tp += GetColor(c.BgColor, true)
			}

			tp += c.Tile
			prev = c
		}
		if i != len(w)-1 {
			tp += "\n"
		}
	}
	fmt.Println(tp)
}

/*returns color palette*/
func palette() map[string]color.Attribute {
	colors := make(map[string]color.Attribute)
	colors["green"] = color.FgGreen
	colors["yellow"] = color.FgYellow
	colors["blue"] = color.FgBlue
	colors["red"] = color.FgRed
	colors["cyan"] = color.FgCyan
	colors["black"] = color.FgBlack
	colors["white"] = color.FgWhite
	colors["magenta"] = color.FgMagenta
	return colors
}

func bgPalette() map[string]color.Attribute {
	colors := make(map[string]color.Attribute)
	colors["green"] = color.BgGreen
	colors["yellow"] = color.BgYellow
	colors["blue"] = color.BgBlue
	colors["red"] = color.BgRed
	colors["cyan"] = color.BgCyan
	colors["black"] = color.BgBlack
	colors["white"] = color.BgWhite
	colors["magenta"] = color.BgMagenta
	return colors
}

func TermboxPalette() map[string]termbox.Attribute {
	colors := make(map[string]termbox.Attribute)
	colors["green"] = termbox.ColorGreen
	colors["yellow"] = termbox.ColorYellow
	colors["blue"] = termbox.ColorBlue
	colors["red"] = termbox.ColorRed
	colors["cyan"] = termbox.ColorCyan
	colors["black"] = termbox.ColorBlack
	colors["white"] = termbox.ColorWhite
	colors["magenta"] = termbox.ColorMagenta
	return colors
}

/*Adds world to a world as another layer. Will append tiles if needed. */
func ReturnWithLayers(world [][]Tile, layers []Layer) ([][]Tile, error) {
	var color string
	w := DuplicateWorld(world)
	for i := 0; i < len(layers); i++ {

		/* Appending new tiles if needed */
		for len(layers[i].World)+layers[i].Pos.X > len(w) {
			w = append(w, []Tile{})
		}

		for j := 0; j < len(layers[i].World); j++ {
			for len(layers[i].World[j])+layers[i].Pos.Y > len(w[j]) {
				w[j] = append(w[j], Tile{})
			}
			for k := 0; k < len(layers[i].World[0]); k++ {
				if layers[i].World[j][k].BgColor == "" {
					color = w[layers[i].Pos.X+j][layers[i].Pos.Y+k].BgColor
					w[layers[i].Pos.X+j][layers[i].Pos.Y+k] = layers[i].World[j][k]
					w[layers[i].Pos.X+j][layers[i].Pos.Y+k].BgColor = color
				} else {
					w[layers[i].Pos.X+j][layers[i].Pos.Y+k] = layers[i].World[j][k]
				}
			}
		}
	}
	return w, nil
}

/*Compares if worlds are the same*/

func CompareWorlds(world1, world2 [][]Tile) bool {
	var toReturn bool
	if len(world1) != len(world2) {
		toReturn = false
	} else {
		for i := 0; i < len(world1); i++ {
			if len(world1[i]) != len(world2[i]) {
				toReturn = false
			} else {
				for j := 0; j < len(world1[i]); j++ {
					if world1[i][j] != world2[i][j] {
						toReturn = false
					} else {
						toReturn = true
					}
				}
			}
		}
	}
	return toReturn
}

/* Animates an int. Still doesnt work when duration is smaller, than the difference between the two input numbers*/
func Animate(num1 *int, num2, duration int) {
	frequency := float64(duration) / math.Abs(float64(*num1-num2))
	moveLenght := 1

	var dir int
	if *num1 < num2 {
		dir = 1
	} else {
		dir = -1
	}

	for *num1 != num2 {
		if dir == 1 && *num1 >= num2-moveLenght {
			*num1 = num2
			return
		} else if dir == -1 && *num1 <= num2+moveLenght {
			*num1 = num2
			return
		} else {
			*num1 += moveLenght * dir
			time.Sleep(time.Duration(frequency) * 1000 * time.Millisecond)
		}
	}
}

/* Converts an array of colors into a world of halfblocks. */

func ParsePixMap(pix [][]string) [][]Tile {
	w := [][]Tile{}
	for x := range pix {
		if x%2 == 0 {
			w = append(w, []Tile{})
		}
		for y := range pix[x] {
			if x%2 == 0 {
				if pix[x][y] == "" {
					w[x/2] = append(w[x/2], Tile{Tile: "▀", Color: "000000"})
				} else {
					w[x/2] = append(w[x/2], Tile{Tile: "▀", Color: pix[x][y]})
				}
				
			} else {
				if pix[x][y] == "" {
					w[x/2][y].BgColor = "000000"
				} else {
					w[x/2][y].BgColor = pix[x][y]		
				}
			}
		}
	}
	return w
}

/* Converts a world to pixmap */

func WorldToPixMap(w [][]Tile) [][]string {
	pix := [][]string{}

	for x := range w {
		pix = append(pix, []string{})
		for y := range w[x] {
			pix[x] = append(pix[x], w[x][y].Color)
		}
	}
	return pix
}

/* Same as ReturnWithLayers, but for pixmap. */

func ReturnWithPixLayers(pixmap [][]string, layers []PixLayer) [][]string {
	for i := 0; i < len(layers); i++ {
		for j := 0; j < len(layers[i].PixMap); j++ {
			for k := 0; k < len(layers[i].PixMap[j]); k++ {
				if layers[i].PixMap[j][k] != "" {
					pixmap[layers[i].Pos.X+j][layers[i].Pos.Y+k] = layers[i].PixMap[j][k]		
				}
			}
		}
	}
	return pixmap
}

/* Duplicates a pixmap. */

func DuplicatePix(p1 [][]string) [][]string {
	p2 := [][]string{}
	for i := range p1 {
		p2 = append(p2, []string{})
		for j := range p1[i] {
			p2[i] = append(p2[i], p1[i][j])
		}
	}
	return p2
}

/* Makes a cut from a pixmap. */
func CutPix(pix [][]string, from Ve2, to Ve2) ([][]string, error) {
	if from.X < 0 {
		to.X += from.X * -1
		from.X = 0
	}
	if from.Y < 0 {
		to.Y += from.Y * -1
		from.Y = 0
	}
	if to.X >= len(pix) {
		to.X = len(pix) - 1
	}
	if to.Y >= len(pix[0]) {
		to.Y = len(pix[0]) - 1
	}

	if /*from.X+to.X <= len(pix) && from.Y+to.Y <= len(pix[0])*/ true {
		var toReturn [][]string
		var toAppend []string
		for i := 0; i <= to.X; i++ {
			if i+from.X >= len(pix) {
				break
			}
			toReturn = append(toReturn, toAppend)
			for j := 0; j <= to.Y; j++ {
				if j+from.Y >= len(pix[i]) {
					break
				}
				toReturn[i] = append(toReturn[i], pix[i+from.X][j+from.Y])
			}
		}
		return toReturn, nil
	} else {
		return pix, errors.New("cut: Out of boundaries")
	}
}

/* Applies changes using termbox */
func ApplyChanges(w [][]Tile, changes []Ve2, centered bool) {
	var cTile Tile
	offset := V2(0, 0)
	if centered {
		offset = GetSize()
		offset.X /= 2
		offset.X -= len(w)
		offset.Y /= 2
		offset.Y -= len(w[0])
	}

	for i := range changes {
		cTile = w[changes[i].X][changes[i].Y]
		termbox.SetCell(changes[i].X+offset.X, changes[i].Y+offset.Y, rune(cTile.Tile[0]), TermboxPalette()[cTile.Color], TermboxPalette()[cTile.BgColor])
	}
}

/* Writes a world to Termbox buffer. You have to Init and Sync termbox yourself. If tile is longer than one character, cliw can deal with it. */
func WriteToTermbox(w [][]Tile) {
	termbox.Init()
	defer termbox.Close()
	for i := range w {
		for j := range w[i] {
			/* has to switch x and y */
			termbox.SetCell(i, j, '▀', TermboxPalette()[w[i][j].Color], TermboxPalette()[w[i][j].BgColor])
		}
	}
	termbox.Sync()
}

/* Moves cursor to specified location */
func SetCursor(pos Ve2) {
	fmt.Printf("\033[%d;%df", pos.X, pos.Y)
}

/* Clears the screen */
func Clear() {
	fmt.Println("\033[2J")
}

/* Sets a world column to be drawn at specific terminal column */ 
func ColumnMargin(w [][]Tile, c, cPos int) [][]Tile {
	for i := range w {
		if c < len(w[i]) {
			w[i][c].Tile = fmt.Sprintf("\033[%dG", cPos) + w[i][c].Tile
		}
	}
	return w
}

func LoadPixMap(path string) [][]string {
	dat, _ := ioutil.ReadFile(path)
	var tr [][]string
	json.Unmarshal(dat, &tr)
	return tr
}

func SavePixMap(pm [][]string, path string) {
	dat, _ := json.Marshal(pm)
	ioutil.WriteFile(path, dat, 0644)	
}

func FlipPixMapV(pm [][]string) (tr [][]string) {
	for i := range pm {
		tr = append(tr, []string{})
		for j := range pm[i] {
			tr[i] = append(tr[i], pm[i][len(pm[i]) - 1 - j])
		}
	}
	return
}

func PNGToPixMap(filename string) (tr [][]string, err error) {
	var col struct{
		r uint32
		g uint32
		b uint32
		a uint32
	}
	var ta string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}

	for i:=img.Bounds().Min.Y; i < img.Bounds().Max.Y; i++ {
		tr = append(tr, []string{})
		for j:=img.Bounds().Min.X; j < img.Bounds().Max.X; j++ {
			// this will probably crash with unsymetrical images. we will see :]
			col.r, col.g, col.b, col.a = img.At(j, i).RGBA() 

			if col.a == 0 {
				tr[i] = append(tr[i], "")
				continue
			}
			
			ta = strconv.FormatInt(int64(col.r), 16)
			ta += strconv.FormatInt(int64(col.g), 16)
			ta += strconv.FormatInt(int64(col.b), 16)

			tr[i] = append(tr[i], ta)
		}
	}
	return
}
