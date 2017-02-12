package general

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"time"
)

type Player struct {
	Name         string
	HoldingCards []*Card
	Socket       *websocket.Conn
	IsTurn       func() bool
	CurrentRoom  *GameRoom
}

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
	EVENT_PLAY
	EVENT_TALK
)

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      string
	Timestamp int // Unix timestamp (secs)
	Content   string
}

type ICommunicate interface {
	Send(content string)
}

func (this *Player) JoinRoom(game *GameRoom) {
	this.CurrentRoom = game
}

func (this *Player) LeaveRoom() {
	this.CurrentRoom = nil
}

func (this *Player) PlayCard(card Card) Card {
	var index int = 0
	beego.Info("user play card", card)
	for ; index < len(this.HoldingCards)-1; index++ {
		comparingCard := this.HoldingCards[index]
		if comparingCard.Color == card.Color && comparingCard.Number == card.Number {
			break
		}
	}
	var current string = ""
	for _, card := range this.HoldingCards {
		current += fmt.Sprint(*card)
	}
	beego.Info("before slice" + current)
	returning := this.HoldingCards[index]
	beego.Info("play card: ", *returning)
	this.HoldingCards = append(this.HoldingCards[0:index], this.HoldingCards[index+1:]...)
	current = ""
	for _, card := range this.HoldingCards {
		current += fmt.Sprint(*card)
	}
	beego.Info("after slice" + current)
	this.Socket.WriteMessage(websocket.TextMessage, this.wrapCurrentHolding("UPDATE"))
	this.CurrentRoom.BroadcastContent(this.Name + " 打出了" + card.Showing)
	this.CurrentRoom.PlayerPlayed(this, card)
	return *returning
}

func (this *Player) Turn(game IGameLogic) {
	err := game.AssignCardToPlayer(this)
	if err != nil {
		this.Send("GameEnd")
	}
	this.Socket.WriteMessage(websocket.TextMessage, this.wrapCurrentHolding("UPDATE"))
}

func (this *Player) Send(content string) {
	data := Event{EVENT_MESSAGE, this.Name, int(time.Now().Unix()), content}
	translated, _ := json.Marshal(data)
	this.Socket.WriteMessage(websocket.TextMessage, translated)
}

func (this *Player) Talk(content string) {
	data := Event{EVENT_TALK, this.Name, int(time.Now().Unix()), content}
	translated, _ := json.Marshal(data)
	for _, player := range this.CurrentRoom.Players {
		player.Socket.WriteMessage(websocket.TextMessage, translated)
	}
}

func (this *Player) StartGame() {
	this.Socket.WriteMessage(websocket.TextMessage, this.wrapCurrentHolding("START"))
}

func (this *Player) wrapCurrentHolding(operation string) []byte {
	cards := map[string]interface{}{"operation": "START", "value": this.HoldingCards}
	translated, err := json.Marshal(cards)
	if err != nil {
		beego.Error(err)
	}
	data := Event{EVENT_PLAY, this.Name, int(time.Now().Unix()), string(translated)}
	content, err := json.Marshal(data)
	if err != nil {
		beego.Error(err)
	}
	return content
}
