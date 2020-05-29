package clengine

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type Player struct {
	hp        int
	attack    int
	defense   int
	inventory Inventory
	name      string
	money     int
	quests    []Quest
}

type Inventory struct {
	weightLimit int
	items       []Item
}

type Item struct {
	avgPrice   int
	weight     int
	durability int
	attack     int
	canBuild   bool
	stolen     bool
	legal      bool
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

func NewVe2(x, y int) Ve2 {
	return Ve2{X: x, Y: y}
}

//changes one specific tile in the world
func EditTile(world [][]Tile, pos Ve2, t Tile) ([][]Tile, error) {
	if pos.X <= 0 || pos.Y <= 0 || pos.X > len(world) || pos.Y > len(world) {
		return nil, errors.New("You entered value smaller than zero")
	} else {
		world[pos.X][pos.Y] = t
		return world, nil
	}
}

//changes all tiles in a rectangular shape
func EditWorld(world [][]Tile, fromX, fromY, toX, toY int, tile Tile) ([][]Tile, error) {
	if fromX < 0 || fromY < 0 || toX < fromX || toY < fromY {
		return nil, errors.New("Invalid number")
	} else {
		for len(world) <= fromX+toX {
			world = append(world, make([]Tile, 0))
		}
		for i := 0; i <= toX; i++ {
			for len(world[fromX+i]) <= fromY+toY {
				//fmt.Println("x")
				world[fromX+i] = append(world[fromX+i], Tile{})
			}
		}

		for i := 0; i < toX; i++ {
			for r := 0; r <= toY; r++ {
				world[fromX+i][fromY+r] = tile
			}
		}
		return world, nil
	}
}

//returns, how much does the inventory weight
func InventoryWeight(inv Inventory) int {
	var weight int
	for i := 0; i < len(inv.items); i++ {
		weight += inv.items[i].weight
	}
	return weight
}

//adds item to inventory and automaticaly checks weight
func AddToInventory(inv Inventory, toAdd Item) (int, error) {
	if InventoryWeight(inv) < inv.weightLimit {
		inv.items = append(inv.items, toAdd)
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

//loads world from file
func LoadWorld(path string, world *[][]Tile) {
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
		for len(*world) <= x {
			(*world) = append(*world, make([]Tile, 0))
		}
		for len((*world)[x]) <= y {
			(*world)[x] = append((*world)[x], Tile{})
		}
		(*world)[x][y] = Tile{Name: text[i+2], Tile: text[i+3], Damage: damage, Color: text[i+5]}
		i += 5
	}
}

//prints out the world to terminal
func DrawWorld(world [][]Tile) {
	var c Tile
	palette := palette()
	col := color.New(color.FgWhite)
	for i := 0; i < len(world); i++ {
		for j := 0; j < len(world[0]); j++ {
			c = world[i][j]
			col = color.New(palette[c.Color])
			col.Print(c.Tile)
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
