package majiang

import (
	"errors"
	"github.com/astaxie/beego"
	"math/rand"
	m "samples/WebIM/models/general"
)

type GameMaJiang struct {
	RemainingCards []m.Card
	AssignedCards  map[int]bool
	UsedCards      int
}

const maxPlayer = 2

func (this *GameMaJiang) GameStart(players map[string]*m.Player) {
	beego.Info("start the game")
	this.RemainingCards = make([]m.Card, 0)
	this.UsedCards = 0
	for _, c := range []string{"W", "T", "B"} {
		for i := 1; i <= 9; i++ {
			card := m.Card{c, getShowing(c, i), i}
			this.RemainingCards = append(this.RemainingCards, card, card, card, card)
		}
	}
	dong := m.Card{"Z", "东风", 1}
	nan := m.Card{"Z", "南风", 2}
	xi := m.Card{"Z", "西风", 3}
	bei := m.Card{"Z", "北风", 4}
	zhong := m.Card{"Z", "中", 5}
	fa := m.Card{"Z", "發", 6}
	bai := m.Card{"Z", "白", 7}
	this.RemainingCards = append(this.RemainingCards,
		dong, dong, dong, dong,
		nan, nan, nan, nan,
		xi, xi, xi, xi,
		bei, bei, bei, bei,
		zhong, zhong, zhong, zhong,
		fa, fa, fa, fa,
		bai, bai, bai, bai)
	chun := m.Card{"H", "春", 1}
	xia := m.Card{"H", "夏", 2}
	qiu := m.Card{"H", "秋", 3}
	don := m.Card{"H", "东", 4}
	mei := m.Card{"H", "梅", 5}
	lan := m.Card{"H", "兰", 6}
	zhu := m.Card{"H", "竹", 7}
	ju := m.Card{"H", "菊", 8}
	this.RemainingCards = append(this.RemainingCards,
		chun, chun, chun, chun,
		xia, xia, xia, xia,
		qiu, qiu, qiu, qiu,
		don, don, don, don,
		mei, mei, mei, mei,
		lan, lan, lan, lan,
		zhu, zhu, zhu, zhu,
		ju, ju, ju, ju)
	//assign to randomly player
	this.AssignedCards = make(map[int]bool, 144)
	for _, player := range players {
		player.HoldingCards = make([]*m.Card, 0)
	}

	for _, player := range players {
		for i := 0; i < 13; i++ {
			assignCardToPlayer(this, player)
		}

	}
	for _, player := range players {
		assignCardToPlayer(this, player)
		break
	}
	for _, player := range players {
		player.StartGame()
	}

}

func (this *GameMaJiang) AssignCardToPlayer(player *m.Player) error {
	return assignCardToPlayer(this, player)
}

func (this *GameMaJiang) IsAbleToStart(players map[string]*m.Player) bool {
	return len(players) == maxPlayer
}

func (this *GameMaJiang) CreateGameRoom(id int) *m.GameRoom {

	newRoom := new(m.GameRoom)
	newRoom.Name = "majiang"
	newRoom.ID = id
	newRoom.Players = make(map[string]*m.Player, maxPlayer)
	newRoom.CurrentGame = this
	return newRoom
}

func (this *GameMaJiang) UserPlayed(played *m.Player, card m.Card, room *m.GameRoom) {
	next := determineNext(played, room.Players)
	next.Turn(this)
}

func determineNext(played *m.Player, all map[string]*m.Player) *m.Player {
	for _, play := range all {
		if play.Name != played.Name {
			return play
		}
	}
	return nil
}

func assignCardToPlayer(game *GameMaJiang, player *m.Player) error {
	var assign int = -1
	if game.UsedCards == 144 {
		return errors.New("card used up")
	}
	for assign = rand.Intn(144); true; assign++ {
		if assign >= 144 {
			assign = 0
		}
		if !game.AssignedCards[assign] {
			//beego.Info("assign " + fmt.Sprint(game.RemainingCards[assign]) + " to " + player.Name)
			player.HoldingCards = append(player.HoldingCards, &game.RemainingCards[assign])
			game.UsedCards++
			break
		}
	}
	return nil
}

func getShowing(color string, number int) string {
	var result string = ""
	switch number {
	case 1:
		result += "一"
	case 2:
		result += "二"
	case 3:
		result += "三"
	case 4:
		result += "四"
	case 5:
		result += "五"
	case 6:
		result += "六"
	case 7:
		result += "七"
	case 8:
		result += "八"
	case 9:
		result += "九"
	}
	switch color {
	case "W":
		result += "萬"
	case "T":
		result += "条"
	case "B":
		result += "饼"
	}
	return result

}
