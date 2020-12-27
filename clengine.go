package clengine

import (
	"bytes"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

type Player struct {
	Hp      int
	Defense int
	Inv     Inventory
	Name    string
	Money   int
	Quests  []Quest
	Tile    Tile
	Pos     Ve2
}

type Inventory struct {
	WeightLimit int
	Items       []Item
}

type Item struct {
	AvgPrice   int
	Weight     int
	Durability int
	Attack     int
	CanBuild   bool
	Stolen     bool
	Legal      bool
}

type Quest struct {
	accepted      int
	end           int
	timeToFinnish int
	requester     Character
	message       string
	toDo          bool
	legal         bool
}

type Character struct {
	name  string
	money int
}

type Tile struct {
	Name    string
	Tile    string
	Damage  int
	Color   string
	BgColor string
}

type Ve2 struct {
	X int
	Y int
}

type BattleUnit struct {
	Name       string
	Tile       Tile
	Pos        Ve2
	FocusPoint Ve2
	Health     int
	Weapon     Attack
	Distance   int
}

type Attack struct {
	Name       string
	Tile       Tile
	Damage     int
	AttackRate int
}

type Layer struct {
	World [][]Tile
	Pos   Ve2
}

//Changes one specific tile in the world. Can append tiles.
func EditTile(world [][]Tile, pos Ve2, t Tile) ([][]Tile, error) {
	w := DuplicateWorld(world)
	if pos.X < 0 || pos.Y < 0{
		return w, errors.New("You entered value smaller than zero")
	} else {
		for len(w) <= pos.X{
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

//Gets size of terminal window
func GetSize() (int, int) {
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

	return h, w
}

//Draws world on the center of the screen
//can get additional blank margin for user input in ui
func DrawCentered(w [][]Tile, additionalRow bool) {
	he, wi := GetSize()
	var rows, colls, wwidth int
	//var toPrint string
	var currentRow [][]Tile
	currentRow = append(currentRow, make([]Tile, 0))

	rows = (he-len(w))/2

	//first row has to be the longest
	//TODO: first row doesnt have to be the longest
	for i := range w[0] {
		wwidth += len(w[0][i].Tile)
	}
	colls = (wi-wwidth)/2

	fmt.Println(strings.Repeat("\n", rows))
	for i:=0; i<len(w); i++ {
		fmt.Print(strings.Repeat(" ", colls))
		currentRow[0] = w[i]
		DrawWorld(currentRow)
	}
	if additionalRow {
		fmt.Print(strings.Repeat(" ", colls))
	}
}

//makes cut from the world
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

//returns new Ve2
func V2(x, y int) Ve2 {
	return Ve2{X: x, Y: y}
}

//Adds on Ve2 to another
func (v1 *Ve2) Add(v2 Ve2) {
	v1.X += v2.X
	v1.Y += v2.Y
}

//changes all tiles in a rectangular shape
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

//returns, how much does the inventory weight
func InventoryWeight(inv Inventory) int {
	var weight int
	for i := 0; i < len(inv.Items); i++ {
		weight += inv.Items[i].Weight
	}
	return weight
}

//adds item to inventory and automaticaly checks weight
func AddToInventory(inv Inventory, toAdd Item) (int, error) {
	if InventoryWeight(inv) < inv.WeightLimit {
		inv.Items = append(inv.Items, toAdd)
		return InventoryWeight(inv), nil
	} else {
		return InventoryWeight(inv), errors.New("The item weights too much for you to cary.")
	}
}

//saves world to a text file
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

func SaveWorldJSON(w [][]Tile, path string) {
	toSave, err := json.Marshal(w)
	if err == nil {
		ioutil.WriteFile(path, toSave, 0644)
	}
}

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

//loads world from file
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

func LoadWorldJSON(path string) [][]Tile {
	var w [][]Tile

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &w)
	return w
}

//prints out the world to terminal
func DrawWorld(world [][]Tile) {
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

//returns color palette
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

//Adds layers to a world
func ReturnWithLayers(world [][]Tile, layers []Layer) ([][]Tile, error) {
	var color string
	w := DuplicateWorld(world)
	for i := 0; i < len(layers); i++ {
		for j := 0; j < len(layers[i].World); j++ {
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

func (u *BattleUnit) Ai(time int) {
	attack := false
	lastAttack := 0
	if attack == false {
		if u.Pos.X < u.FocusPoint.X {
			if time%500 == 0 {
				u.Pos.X += 1
			}
		} else if u.Pos.X > u.FocusPoint.X {
			if time%500 == 0 {
				u.Pos.X -= 1
			}
		}
		if u.Pos.Y < u.FocusPoint.Y {
			if time%500 == 0 {
				u.Pos.Y += 1
			}
		} else if u.Pos.Y > u.FocusPoint.Y {
			if time%500 == 0 {
				u.Pos.Y -= 1
			}
		}
	}
	if float64(u.Distance) <= math.Sqrt(float64(((u.Pos.X-u.FocusPoint.X)*(u.Pos.X-u.FocusPoint.X))-((u.Pos.Y-u.FocusPoint.Y)*(u.Pos.Y-u.FocusPoint.Y)))) {
		attack = true
	}
	if time-lastAttack >= u.Weapon.AttackRate && attack == true {
		u.Weapon.Fire(u.Pos, u.FocusPoint)
	}
}

func (w *Attack) Fire(attackerPos, focusPoint Ve2) {
	//TODO: this whole function
}

//Compares if worlds are the same
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

// Converts an array of colors into a world of halfblocks. Final world should be used as layer.
func ParsePixMap(pix [][]string) [][]Tile {
	w := [][]Tile{}
	for x := range pix {
		w = append(w, []Tile{})
		for y := range pix[x] {
			fmt.Println(w)
			if x % 2 == 0 {
				fmt.Println("odd")
				w[x] = append(w[x], Tile{Tile: "▀", Color: pix[x][y]})	
			} else {
				w[x/2][y].BgColor = pix[x][y]
			}
		}
	}
	return w
}

func WorldToPixMap(w [][]Tile) [][]string {
	pix := [][]string{}

	for x := range w {
		pix = append(pix, []string{})
		for y := range w[x] {
			pix[x*2] = append(pix[x], w[x][y].Color)
			pix[x*2+1] = append(pix[x+1], w[x][y].BgColor)
		}
	}
	return pix
}
