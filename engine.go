package cliEngine

import
(
	"bufio"
	"os"
)

func newWorld(scaleX, scaleY int) ([][]string , err) {
	world := [x][y]string{}
	return world
	if(x <= 0 || y <= 0){
		return err
	}
}

func editWorld(world [][]string, posX, posY int, tile string) ([][]string, err) {
	if(posX <= 0 || posY <= 0 || posX > len(world) || posY > len(world)){
		return err
	}
	world[posX, posY] = tile
	return world
}
func loadWorld(file string, scaleY int) ([][]string, err) {
	text := []string{}
	file, err := os.Open(file)
	if(err != nil){
		return nil
	}
	defer file.Close()

	scanner := bufioNewScanner(files)
	for scanner.Scan(){
		text = append(text, scanner.Text())
	}
	for i; i <= scaleY; i++{
		world = append(world, [i]text)
	}
}
func saveWorld(file string, world [][]string) err {
	file, err := os.Open(file)
	if(err != nil){
		return nil
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range world {
		fmt.Fprintln(w, world)
	}
	return w.Flush()
}
