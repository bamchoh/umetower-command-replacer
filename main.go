package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	DefaultUpKey    = "k"
	DefaultDownKey  = "j"
	DefaultLeftKey  = "h"
	DefaultRightKey = "l"
	DefaultBlockKey = " "
)

type config struct {
	Down  string `json:"down"`
	Up    string `json:"up"`
	Left  string `json:"left"`
	Right string `json:"right"`
	Block string `json:"block"`
}

func getRuneKey(key string, defaultKey string) rune {
	if len(key) == 0 {
		key = defaultKey
	}

	fmt.Println(key)

	return rune(key[0])
}

func setUp(mapper map[rune]rune, key string) {
	runekey := getRuneKey(key, DefaultUpKey)
	mapper[runekey] = '8'
}

func setDown(mapper map[rune]rune, key string) {
	runekey := getRuneKey(key, DefaultDownKey)
	mapper[runekey] = '2'
}

func setLeft(mapper map[rune]rune, key string) {
	runekey := getRuneKey(key, DefaultLeftKey)
	mapper[runekey] = '4'
}

func setRight(mapper map[rune]rune, key string) {
	runekey := getRuneKey(key, DefaultRightKey)
	mapper[runekey] = '6'
}

func setBlock(mapper map[rune]rune, key string) {
	runekey := getRuneKey(key, DefaultBlockKey)
	mapper[runekey] = '5'
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getConfig(filename string) config {
	var config config
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func isCommand(comment string, mapper map[rune]rune) bool {
	for _, c := range comment {
		if _, ok := mapper[c]; ok == false {
			return false
		}
	}
	return true
}

func main() {
	configfile := "config.json"
	var config config
	if exists(configfile) {
		config = getConfig(configfile)
	}

	mapper := make(map[rune]rune, 0)
	setUp(mapper, config.Up)
	setDown(mapper, config.Down)
	setLeft(mapper, config.Left)
	setRight(mapper, config.Right)
	setBlock(mapper, config.Block)

	rawurl := ""
	id := ""
	if len(os.Args) > 2 {
		rawurl = os.Args[1]
		id = os.Args[2]
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer c.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comment := scanner.Text()
		command := isCommand(comment, mapper)

		if command {
			for key, val := range mapper {
				comment = strings.Replace(comment, string(key), string(val), -1)
			}
		}

		fmt.Println("send:" + comment)
		text := id + "\t" + comment

		err = c.WriteMessage(websocket.TextMessage, []byte(text))
		if err != nil {
			log.Fatal("write:", err)
			return
		}
	}

}
