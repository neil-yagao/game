// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"container/list"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"

	"samples/WebIM/models"
	general "samples/WebIM/models/general"
	"samples/WebIM/models/majiang"
)

func newEvent(ep general.EventType, user, msg string) general.Event {
	return general.Event{ep, user, int(time.Now().Unix()), msg}
}

func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

func Leave(user string) {
	unsubscribe <- user
}

type Subscriber struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.
	unsubscribe = make(chan string, 10)
	// Send events here to publish them.
	publish = make(chan general.Event, 10)
	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()
)
var gameLobby *models.GameLobby = new(models.GameLobby)

// This function handles all incoming chan messages.
func chatroom() {
	gameLobby.Rooms = make(map[int]*general.GameRoom, 8)
	gameLobby.Players = make(map[string]*general.Player, 32)
	for {
		select {
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				gameLobby.BroadcastContent(sub.Name + "加入了大厅")
				newPlayer := new(general.Player)
				newPlayer.Name = sub.Name
				newPlayer.Socket = sub.Conn
				gameLobby.Players[sub.Name] = newPlayer
				//for demo case only one room is used
				if gameLobby.Rooms[1] == nil {
					gameLobby.Rooms[1] = new(majiang.GameMaJiang).CreateGameRoom(1)

				}
				gameLobby.Rooms[1].JoinRoom(newPlayer)
				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:

			if event.Type == general.EVENT_MESSAGE {
				gameLobby.DispatchUserOperation(event.Content, gameLobby.Players[event.User])

				beego.Info("Message from", event.User, ";Content:", event.Content)
			}

		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}

					gameLobby.BroadcastContent(sub.Value.(Subscriber).Name + "离开了大厅")

					player := new(general.Player)
					player.Name = sub.Value.(Subscriber).Name
					player.Socket = sub.Value.(Subscriber).Conn
					gameLobby.Rooms[1].LeaveRoom(player)
					beego.Info(player.Name + " has left!")
				}
			}
		}
	}
}

func init() {
	go chatroom()
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
