package clengine

import
(
	"io"
	"errors"
	"os"
	"strconv"
	"fmt"
	"bufio"
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
	items []Item
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
	name         string
	money        int
}

type Tile struct {
	name string
	tile string
	damage int
	color string
}

/*func newWorld(scaleX, scaleY int) ([][]tile , error) {
	world := make([][]string, scaleY)
	for i := 0; i < scaleY; i++ {
		world[i] = make([]tile, scaleX)
	}
	//var world [scaleX][scaleY]string
	if(scaleX <= 0 || scaleY <= 0){
		return nil, errors.New("You entered value smaller than zero")
	} else {
		return world, nil
	}
}*/

func EditTile(world [][]Tile, posX, posY int, t Tile) ([][]Tile, error) {
	if(posX <= 0 || posY <= 0 || posX > len(world) || posY > len(world)){
		return nil, errors.New("You entered value smaller than zero")
	} else {
		world[posX][posY] = t
		return world, nil
	}
}
func EditWorld(world [][]Tile, fromX, fromY, toX, toY int, tile Tile) ([][]Tile, error) {
	if (fromX < 0 || fromY < 0 || toX < fromX || toY < fromY) {
		return nil, errors.New("Invalid number")
	} else {
		for i := 0; i <= toX; i++ {
			world[fromX + i][fromY] = tile
			for r := 0; r <= toY; r++ {
			world[fromX + i][fromY + r] = tile
			}
		}
		return world, nil
	}
}
/*func loadWorld(fileAddress string, scaleY int) ([][]tile, error) {
	world := make([][]tile, 0)
	text := make([]string, 0)
	file, err := os.Open(fileAddress)
	if(err != nil){
		return nil, errors.New("Problem occured, while opening the file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		//text = append(text, scanner.Text()...)
		text[i] = scanner.Text()
	}
	for i := 0; i <= scaleY; i++{
		//world = append(world, text[i])
		world[0][i].tile = text[i]
	}
	return world, nil
}
func saveWorld(fileAdress string, world [][]tile, height int) error {
	file, err := os.Open(fileAdress)
	if(err != nil){
		return errors.New("Problem occured, while saving the file")
	}
	defer file.Close()

	for i := 0; i < height; i++ {
		_, err = io.WriteString(file, world[0][i].tile)
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}
func worldToString(world [][]tile) [][]string{
	stringWorld := make([][]string, len([][]world))
	for i:=0; i < len(world); i++{
		for i:=0; i < len([]world); i++{
			stringWorld[int(i/len([]world))][i] = append(stringWorld, world[i].tile)
		}
	}
	return stringWorld
}
func drawWorld(world [][]tile, worldHeight, playerPosX, playerPosY int, playerTile string) error {
	worldToDraw := world
	worldToDraw[playerPosX][playerPosY] = playerTile
	if (worldHeight <= 0){
		return errors.New("You entered value smaller than zero")
	} else {
		for i := 0; i <= worldHeight; i++ {
			fmt.Println(worldToDraw[i][0])
		}
		return nil
	}
}*/
func InventoryWeight(inv Inventory) int {
	var weight int
	for i:=0; i < len(inv.items); i++{
		weight += inv.items[i].weight
	}
	return weight
}
func AddToInventory(inv Inventory, toAdd Item) (int, error) {
	if InventoryWeight(inv) < inv.weightLimit {
		inv.items = append(inv.items, toAdd)
		return InventoryWeight(inv), nil
	} else {
		return InventoryWeight(inv), errors.New("The item weights too much for you to cary.")
	}
}
func SaveWorld(world [][]Tile, path string){
	var c Tile
	//row := []world
	file, _ := os.Open(path)
	for i:=0; i<len(world); i++{
		for j:=0; j<len(world[0]); j++{
			c = world[i][j]
			io.WriteString(file, string(i) + "\n" + string(j) + "\n" + c.name + "\n" + c.tile + "\n" + string(c.damage) + "\n" + c.color)
			fmt.Println(string(i) + "\n" + string(j) + "\n" + c.name + "\n" + c.tile + "\n" + string(c.damage) + "\n" + c.color)
		}
	}
}
func LoadWorld(path string, world *[][]Tile){
	//var readed []byte
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
		//fmt.Println(scanner.Text())
  }

	/*
	readed, _ = ioutil.ReadFile(path)
	text = append(text, "")
	row = 0
	for i:=0;i<len(readed);i++{
		if string(readed[i]) == ","{
			row += 1
			text = append(text, "")
		} else {
			text[row] = string(readed[i])
		}
	}*/

	//fmt.Println(text)

	for i:=0;i<len(text)-5; i++{
		damage, _ = strconv.Atoi(text[i+4])
		x, _ = strconv.Atoi(text[i])
		y, _ = strconv.Atoi(text[i+1])
		for len(*world) <= x{
			(*world) = append(*world, make([]Tile, 0))
		}
		for len((*world)[x]) <= y{
			(*world)[x] = append((*world)[x], Tile{})
		}
		(*world)[x][y] = Tile{name: text[i+2], tile: text[i+3], damage: damage, color: text[i+5]}
		i+=5
		//fmt.Println(*world)
	}
}
func DrawWorld(world [][]Tile){
	var c Tile
	var toPrint string
	for i:=0; i<len(world); i++{
		//fmt.Println(world[i])
		for j:=0; j<len(world[0]); j++{
			c = world[i][j]
			toPrint = c.color + c.tile
			fmt.Print(toPrint)
		}
		fmt.Println("")
	}
}