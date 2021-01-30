package cliw

import (
	"io/ioutil"
	"encoding/json"
	"strings"
)

func LoadFont(path string) map[string]interface{} {
	dat, _ := ioutil.ReadFile(path)

	var tr map[string]interface{}
	json.Unmarshal(dat, &tr)
	return tr
}

func TextToPixMap(font map[string]interface{}, color, bgcolor, text string) [][]string {
	var chars [][][]string
	var split [][]string
	var tr [][]string
	for i:=0; i < 3; i++ {
		tr = append(tr, []string{})
		for j:=0; j < len(text)*4; j++ {
			tr[i] = append(tr[i], "black")
		}
	}

	for i := range tr {
		for j := range tr[i] {
			tr[i][j] = "000000"
		}
	}

	for i := range text {
		split = [][]string{
			strings.Split(font[string(text[i])].(string), "")[0:3],
			strings.Split(font[string(text[i])].(string), "")[3:6],
			strings.Split(font[string(text[i])].(string), "")[6:9],
		}

		for j := range split {
			for k := range split {
				if split[j][k] != " " {
					split[j][k] = color
				} else {
					split[j][k] = bgcolor
				}
			}
		}
		
		chars = append(chars, split)
	}

	for i := range chars {
		tr = ReturnWithPixLayers(tr, []PixLayer{
			PixLayer{
				Pos: Ve2{0, i*4}, PixMap: chars[i],
			},
		})
	}

	return tr
}
