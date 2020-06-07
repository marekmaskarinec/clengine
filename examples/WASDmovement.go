package main

import (
	"clengine"
	"fmt"
	"os"
	"os/exec"
)

type game struct {
	world  [][]clengine.Tile
	wtd    [][]clengine.Tile
	player clengine.Tile
	pp     clengine.Ve2
}

func run(com, arg string) {
	cmd := exec.Command(com)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func returnWorld(x [][]clengine.Tile) [][]clengine.Tile {
	return x
}

func main() {
	var err error
	//defining game variables
	g := game{}                               //creating new game
	clengine.LoadWorld("world.txt", &g.world) //loading world
	g.player = clengine.Tile{Tile: "X"}       //creating player tile
	g.pp = clengine.NewVe2(0, 0)              //setting player's position
	//g.world, err = clengine.EditWorld(g.world, clengine.NewVe2(0, 0), clengine.NewVe2(2, 2), clengine.Tile{Tile: "^", Color: "blue"})

	g.wtd = g.world

	//bu := clengine.BattleUnit{Name: "unit", Tile: clengine.Tile{Tile: ":", Color: "red"}, Pos: clengine.NewVe2(2, 3), FocusPoint: g.pp, Health: 10, Distance: 2}

	var refresh, tick int

	//game runtime
	var key string
	for {
		go fmt.Scan(&key)
		if key == "q" {
			break
		} else if key == "w" {
			if g.pp.X > 0 {
				g.pp.X--
				refresh = 4000
			}
		} else if key == "s" {
			if g.pp.X < len(g.world)-1 {
				g.pp.X++
				refresh = 4000
			}
		} else if key == "a" {
			if g.pp.Y > 0 {
				g.pp.Y--
				refresh = 4000
			}
		} else if key == "d" {
			if g.pp.Y < len(g.world[g.pp.X])-1 {
				g.pp.Y++
				refresh = 4000
			}
		}
		key = ""

		g.wtd, err = clengine.EditTile(g.wtd, g.pp, g.player)
		check(err)
		if refresh == 4000 {
			run("clear", "")
			clengine.DrawWorld(g.wtd)
			refresh = 0
		} else {
			refresh++
		}
		clengine.LoadWorld("world.txt", &g.wtd) //loading world
		tick++
	}
}
