package clengine

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"

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
	Name   string
	Tile   string
	Damage int
	Color  string
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

//changes one specific tile in the world
func EditTile(world [][]Tile, pos Ve2, t Tile) ([][]Tile, error) {
	w := DuplicateWorld(world)
	if pos.X < 0 || pos.Y < 0 || pos.X > len(w) || pos.Y > len(w) {
		return w, errors.New("You entered value smaller than zero")
	} else {
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

//prints out the world to terminal
func DrawWorld(world [][]Tile) {
	var c Tile
	palette := palette()
	col := color.New(color.FgWhite)
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			c = world[i][j]
			if c.Color == "blnk" {
				fmt.Print(c.Tile)
			} else {
				col = color.New(palette[c.Color])
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
	return colors
}

func ReturnWithLayers(world [][]Tile, layers []Layer) ([][]Tile, error) {
	for i := 0; i < len(layers); i++ {
		for j := 0; j < len(layers[i].World); j++ {
			for k := 0; k < len(layers[i].World[0]); k++ {
				world[layers[i].Pos.X+j][layers[i].Pos.Y+k] = layers[i].World[j][k]
			}
		}
	}
	return world, nil
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
