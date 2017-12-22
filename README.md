# umetower-command-replacer
This tool replaces Umetower command 2 4 6 8 to h j k l

I respects T. Umezawa who are Youtuber.
He is developing an online game "UMETUKUSHI no TOU(Tower of UMETUKUSHI) which is a social game based on Youtube.
You can see his live coding from below URL:
https://www.youtube.com/channel/UCzXL5v5_-L-s4cDhK3I35Dg

You can play the game of "UMETUKUSHI no TOW" by commenting on his live streaming. It is a purpose of this game that you leads a player to a flag. you can operates the player by typing 2/4/6/8 as his youtube streaming comment.
2 is Down
4 is Left
6 is Right
8 is Up
5 is Take / Put a Block

On the other hands, I love Vim. Vim moves a cursor by h/j/k/l.

This tool is for Vimmer. This tool replaces h/j/k/l to 4/2/8/6 and send a comment via Websocket.

# build
```
$ go build 
```

# usage
```
$ umetower-command-replacer <websocket url> <user id>
```

websocket url and user id is not public.
When he will public them, I update this README
