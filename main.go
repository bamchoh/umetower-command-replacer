package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
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
		command := true
		for _, c := range comment {
			if c == 'h' || c == 'j' || c == 'k' || c == 'l' || c == ' ' {
				continue
			}
			command = false
			break
		}

		if command {
			comment = strings.Replace(comment, "h", "4", -1)
			comment = strings.Replace(comment, "j", "2", -1)
			comment = strings.Replace(comment, "k", "8", -1)
			comment = strings.Replace(comment, "l", "6", -1)
			comment = strings.Replace(comment, " ", "5", -1)
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
