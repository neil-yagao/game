package models

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"math/rand"
	g "samples/WebIM/models/general"
	"samples/WebIM/models/majiang"
	"strconv"
	"strings"
)

type GameLobby struct {
	Rooms   map[int]*g.GameRoom
	Players map[string]*g.Player
}

const MAX_PLAYER_NUMBER int = 2

func (this *GameLobby) DispatchUserOperation(action string, player *g.Player) {
	var operation g.Operation
	beego.Info("received action:" + action)
	err := json.Unmarshal([]byte(action), &operation)
	if err != nil {
		beego.Error(err)
		return
	}
	value := strings.Split(operation.Value, "$")
	switch operation.Code {

	case g.CREATE_ROOM:
		var id int = rand.Intn(1024)
		for this.Rooms[id] != nil {
			id = rand.Intn(1024)
		}
		if value[0] == "majiang" {
			majiang := new(majiang.GameMaJiang)
			this.Rooms[id] = majiang.CreateGameRoom(id)
		}
	case g.JOIN_ROOM:
		roomId, _ := strconv.Atoi(value[0])
		if this.Rooms[roomId].State != "started" {
			this.Rooms[roomId].JoinRoom(player)
		}
	case g.LEAVE_ROOM:
		roomId, _ := strconv.Atoi(value[0])
		this.Rooms[roomId].LeaveRoom(player)
	case g.PLAY:
		var card g.Card
		beego.Info("user playing " + value[0])
		json.Unmarshal([]byte(value[0]), &card)
		player.PlayCard(card)
	case g.TALK:
		beego.Info("user talk" + value[0])
		player.Talk(value[0])
	}

}

func (this *GameLobby) BroadcastContent(content string) {
	for _, player := range this.Players {
		player.Send(content)
	}
}
