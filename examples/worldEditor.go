package main

import(
	"fmt"
	"clengine"
	"strings"
	"strconv"
	"os/exec"
	"os"
)

func main(){
	var in string = "0"
	var pos1, pos2 clengine.Ve2
	var t clengine.Tile
	var world [][]clengine.Tile
	var err error

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout

	fmt.Println("insert world file location:")
	fmt.Scanln(&in)
	cmd.Run()
	clengine.LoadWorld(in, &world)
	fmt.Println("world loaded")
	fmt.Println("h: help\ne: edit tile\nq: quit\ns: save\nd: draw")

	for{
		fmt.Scanln(&in)
		if strings.ToLower(in) == "e"{
			cmd.Run()
			//getting tile position
			fmt.Println("Set x position:")
			fmt.Scanln(&in)
			pos1.X, _ = strconv.Atoi(in)
			fmt.Println("Set y position:")
			fmt.Scanln(&in)
			pos1.Y, _ = strconv.Atoi(in)

			fmt.Println("Set color:")
			fmt.Scanln(&in)
			t.Color = in

			fmt.Println("Set tile:")
			fmt.Scanln(&in)
			t.Tile = in

			cmd.Run()

			world, err = clengine.EditTile(world, pos1, t)
			if err != nil{
				fmt.Println(err)
				panic(err)
			}
			clengine.DrawWorld(world)
		} else if strings.ToLower(in) == "q"{
			break
		} else if strings.ToLower(in) == "s"{
			fmt.Println("set filename:")
			fmt.Scanln(&in)
			if in != ""{
				clengine.SaveWorld(world, in)
			} else {
				clengine.SaveWorld(world, "out.txt")
			}
			fmt.Println("world saved")
		} else if strings.ToLower(in) == "d"{
			clengine.DrawWorld(world)
		} else if strings.ToLower(in) == "o"{
			fmt.Println("set file location:")
			fmt.Scanln(&in)
			if in != ""{
				clengine.LoadWorld(in, &world)
			} else {
				fmt.Println("no location set")
			}
			fmt.Println("world loaded")
		} else if strings.ToLower(in) == "r"{
			cmd.Run()
			//getting tile position
			fmt.Println("Set start x position:")
			fmt.Scanln(&in)
			pos1.X, _ = strconv.Atoi(in)
			fmt.Println("Set start y position:")
			fmt.Scanln(&in)
			pos1.Y, _ = strconv.Atoi(in)

			fmt.Println("Set height:")
			fmt.Scanln(&in)
			pos2.X, _ = strconv.Atoi(in)
			fmt.Println("Set width:")
			fmt.Scanln(&in)
			pos2.Y, _ = strconv.Atoi(in)

			fmt.Println("Set color:")
			fmt.Scanln(&in)
			t.Color = in

			fmt.Println("Set tile:")
			fmt.Scanln(&in)
			t.Tile = in

			cmd.Run()
			world, err = clengine.EditWorld(world, pos1, pos2, t)
			if err != nil {
				fmt.Println("Editing failed")
			} else {
				fmt.Println("world edited succsessfully")
				clengine.DrawWorld(world)
			}

		} else {
			cmd.Run()
			fmt.Println("h: help\ne: edit tile\nq: quit\ns: save\nd: draw\no: open")
		}
	}
}
