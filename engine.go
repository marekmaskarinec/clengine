package cliEngine

import
(
	"io"
	"errors"
	"bufio"
	"os"
	"fmt"
)

func newWorld(scaleX, scaleY int) ([][]string , error) {
	world := make([][]string, scaleY)
	for i := 0; i < scaleY; i++ {
		world[i] = make([]string, scaleX)
	}
	//var world [scaleX][scaleY]string
	if(scaleX <= 0 || scaleY <= 0){
		return nil, errors.New("You entered value smaller than zero")
	} else {
		return world, nil
	}
}

func editWorld(world [][]string, posX, posY int, tile string) ([][]string, error) {
	if(posX <= 0 || posY <= 0 || posX > len(world) || posY > len(world)){
		return nil, errors.New("You entered value smaller than zero")
	} else {
		world[posX][posY] = tile
		return world, nil
	}
}
func loadWorld(fileAddress string, scaleY int) ([][]string, error) {
	world := make([][]string, 0)
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
		world[0][i] = text[i]
	}
	return world, nil
}
func saveWorld(fileAdress string, world [][]string, height int) error {
	file, err := os.Open(fileAdress)
	if(err != nil){
		return errors.New("Problem occured, while saving the file")
	}
	defer file.Close()

	for i := 0; i < height; i++ {
		_, err = io.WriteString(file, world[0][i])
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}
func drawWorld(world [][]string, worldHeight int) error {
	if (worldHeight <= 0){
		return errors.New("You entered value smaller than zero")
	} else {
		for i := 0; i <= worldHeight; i++ {
			fmt.Println(world[i][0])
		}
		return nil
	}
}
