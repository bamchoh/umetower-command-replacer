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

	DownCode  = byte(0)
	LeftCode  = byte(1)
	RightCode = byte(2)
	UpCode    = byte(3)
	BlockCode = byte(4)
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

	return rune(key[0])
}

func setUp(mapper map[rune]byte, key string) {
	runekey := getRuneKey(key, DefaultUpKey)
	mapper[runekey] = UpCode
}

func setDown(mapper map[rune]byte, key string) {
	runekey := getRuneKey(key, DefaultDownKey)
	mapper[runekey] = DownCode
}

func setLeft(mapper map[rune]byte, key string) {
	runekey := getRuneKey(key, DefaultLeftKey)
	mapper[runekey] = LeftCode
}

func setRight(mapper map[rune]byte, key string) {
	runekey := getRuneKey(key, DefaultRightKey)
	mapper[runekey] = RightCode
}

func setBlock(mapper map[rune]byte, key string) {
	runekey := getRuneKey(key, DefaultBlockKey)
	mapper[runekey] = BlockCode
}

func printMapping(mapper map[rune]byte) {
	codes := []byte{UpCode, DownCode, LeftCode, RightCode, BlockCode}
	names := []string{"Up   ", "Down ", "Left ", "Right", "Block"}
	for i, c := range codes {
		for k, v := range mapper {
			if c == v {
				fmt.Println(names[i], "=>", string(k))
			}
		}
	}
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

func isCommand(comment string, mapper map[rune]byte) bool {
	for _, c := range comment {
		if _, ok := mapper[c]; ok == false {
			return false
		}
	}
	return true
}

func buildCommandMessage(basemsg []byte, keystat byte, opcode byte) []byte {
	cmd := make([]byte, 0)
	basemsg = append(basemsg, keystat)
	basemsg = append(basemsg, opcode)
	cmd = append(cmd, byte(len(basemsg)+1))
	cmd = append(cmd, basemsg...)
	return cmd
}

func sendCommand(conn *websocket.Conn, id, comment string, mapper map[rune]byte) (err error) {
	for key, val := range mapper {
		comment = strings.Replace(comment, string(key), string(val), -1)
	}

	msg := make([]byte, 0)
	msg = append(msg, id...)
	for _, c := range comment {
		opcode := byte(c)
		keydownmsg := buildCommandMessage(msg, byte(1), opcode)
		fmt.Println(keydownmsg)
		err = conn.WriteMessage(websocket.BinaryMessage, []byte(keydownmsg))
		if err != nil {
			return err
		}
		keyupmsg := buildCommandMessage(msg, byte(0), opcode)
		fmt.Println(keyupmsg)
		err = conn.WriteMessage(websocket.BinaryMessage, []byte(keyupmsg))
		if err != nil {
			return err
		}
	}
	return nil
}

func sendComment(conn *websocket.Conn, id, comment string) (err error) {
	text := id + "\t" + comment
	err = conn.WriteMessage(websocket.TextMessage, []byte(text))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	configfile := "config.json"
	var config config
	if exists(configfile) {
		config = getConfig(configfile)
	}

	mapper := make(map[rune]byte, 0)
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

	printMapping(mapper)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		comment := scanner.Text()
		command := isCommand(comment, mapper)

		var err error
		if command {
			err = sendCommand(c, id, comment, mapper)
		} else {
			err = sendComment(c, id, comment)
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}
