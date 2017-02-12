package general

import "github.com/astaxie/beego"

type GameRoom struct {
	ID          int
	Name        string
	State       string
	Players     map[string]*Player
	CurrentGame IGameLogic
}

type OperationCode int

const (
	CREATE_ROOM OperationCode = iota
	JOIN_ROOM
	LEAVE_ROOM
	PLAY
	TALK
)

type Operation struct {
	Code  OperationCode `json:"code"`
	Value string        `json:"value"`
}

type IGameLogic interface {
	GameStart(players map[string]*Player)
	AssignCardToPlayer(player *Player) error
	IsAbleToStart(players map[string]*Player) bool
	UserPlayed(played *Player, card Card, room *GameRoom) //
}

type PlayLogic interface {
	Turn()
}

type Card struct {
	Color   string `json:"color"`   //W, T, B, (1 - 9)
	Showing string `json:"showing"` // Z (1 - 7) [东，南，西，北，中，發，白]
	Number  int    `json:"number"`  // H (1 - 8) 梅,兰,竹,菊，春,夏,秋,冬
}

func (this *GameRoom) JoinRoom(newPlayer *Player) {
	this.State = "joined"
	beego.Info(newPlayer.Name+" has join the room ", this.ID)
	this.BroadcastContent(newPlayer.Name + " 加入了房间！")
	this.Players[newPlayer.Name] = newPlayer
	if this.CurrentGame.IsAbleToStart(this.Players) {
		this.CurrentGame.GameStart(this.Players)
		this.State = "started"
		this.BroadcastContent("游戏开始！！")
	}
	newPlayer.JoinRoom(this)

}

func (this *GameRoom) LeaveRoom(player *Player) {
	delete(this.Players, player.Name)
	player.LeaveRoom()
}

func (this *GameRoom) BroadcastContent(content string) {
	for _, player := range this.Players {
		player.Send(content)
	}
}

func (this *GameRoom) PlayerPlayed(player *Player, card Card) {
	this.CurrentGame.UserPlayed(player, card, this)
}
