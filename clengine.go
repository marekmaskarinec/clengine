package clengine

import
(
	"io"
	"errors"
	"os"
	"io/ioutil"
	"strings"
	"strconv"
)

type player struct {
	hp        int
	attack    int
	defense   int
	inventory inventory
	name      string
	money     int
	quests    []quest
}

type inventory struct {
	weightLimit int
	items []item
}

type item struct {
	avgPrice   int
	weight     int
	durability int
	attack     int
	canBuild   bool
	stolen     bool
	legal      bool
}

type quest struct {
	accepted      int
	end           int
	timeToFinnish int
	requester     character
	message       string
	toDo          bool
	legal         bool
}

type character struct {
	name         string
	money        int
}

type tile struct {
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

func editTile(world [][]tile, posX, posY int, t tile) ([][]tile, error) {
	if(posX <= 0 || posY <= 0 || posX > len(world) || posY > len(world)){
		return nil, errors.New("You entered value smaller than zero")
	} else {
		world[posX][posY] = t
		return world, nil
	}
}
func editWorld(world [][]tile, fromX, fromY, toX, toY int, tile string) ([][]tile, error) {
	if (fromX < 0 || fromY < 0 || toX < fromX || toY < fromY) {
		return nil, errors.New("Invalid number")
	} else {
		for i := 0; i <= toX; i++ {
			world[fromX + i][fromY].tile = tile
			for r := 0; r <= toY; r++ {
			world[fromX + i][fromY + r].tile = tile
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
func inventoryWeight(inv inventory) int {
	var weight int
	for i:=0; i < len(inv.items); i++{
		weight += inv.items[i].weight
	}
	return weight
}
func addToInventory(inv inventory, toAdd item) (int, error) {
	if inventoryWeight(inv) < inv.weightLimit {
		inv.items = append(inv.items, toAdd)
		return inventoryWeight(inv), nil
	} else {
		return inventoryWeight(inv), errors.New("The item weights too much for you to cary.")
	}
}
func saveWorld(world [][]tile, path string){
	var c tile
	file, _ := os.Open(path)

	for i:=0; i<len(world); i++{
		for j:=0; j<len(world[1]); j++{
			c = world[i][j]
			io.WriteString(file, string(i) + "\n" + string(j) + "\n" + c.name + "\n" + c.tile + "\n" + string(c.damage) + "\n" + c.color)
		}
	}
}
func loadWorld(path string, world [][]tile){
	var readed []byte
	var text []string
	var damage, x, y int

	readed, _ = ioutil.ReadFile(path)
	text = strings.Split(string(readed), "\n")

	for i:=0;i<len(text); i++{
		damage, _ = strconv.Atoi(text[i+4])
		x, _ = strconv.Atoi(text[i])
		y, _ = strconv.Atoi(text[i+1])
		world[x][y] = tile{name: text[i+2], tile: text[i+3], damage: damage, color: text[i+5]}
		i+=6
	}
}
