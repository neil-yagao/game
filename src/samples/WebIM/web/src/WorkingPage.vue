<template>
    <div>
        <!-- route outlet -->
        <!-- component matched by the route will render here -->
        <ul class="nav nav-pills">
            <li role="presentation"><a>欢迎来到房间：{{user}}</a></li>
            <!-- <li role="presentation" style="float:right"><a href="#/">退出房间</a></li> -->
        </ul>
          <div class="row" style="height:90%; min-height :300px">
                <div class="col-xs-6" style="padding:0">
                    <template v-for="p in progress">
                        <div>{{p}}</div>
                    </template>
                </div>
                  <div class="col-xs-6">
                    <p v-for="p in chats" :key="p.Timestamp">{{p.Content}}</p>
                </div>
            </div>
        <div class="navbar-fixed-bottom">
            <div class="container-fluid" style="text-align:center">
                <card :value="holdings" :turn="yourTurn"></card>
                <div class="input-group" style="margin-bottom:10px">
                  <input type="text" class="form-control" placeholder="说点什么" v-model="chat" @keyup.enter="sendChat">
                  <span class="input-group-btn">
                    <button class="btn btn-default" type="button"  @click="sendChat">发送</button>
                  </span>
                </div><!-- /input-group -->
            </div>
        </div>
    </div>
</template>
<script>
import card from './components/card.vue'
import Vue from 'vue'
import _ from 'lodash'
var ROOM = 1
export default {
    name: 'working-page',
    data() {
        return {
            user: this.$route.params.user,
            progress: [],
            socket: {},
            yourTurn: true,
            holdings :[],
            chats:[],
            chat:""
        }
    },
    mounted: function() {
        self = this
        this.$socket.onmessage = function(event) {
            var data = JSON.parse(event.data);
            console.log(data);
            switch (data.Type) {
                case 0: // JOIN
                    self.progress.push(data.User + ' has join the room')
                    break;
                case 1: // LEAVE
                    self.progress.push(data.User + " left the chat room.");
                    break;
                case 2: // MESSAGE
                    if (self.progress.length == 10){
                        self.progress.shift()
                    }
                    self.progress.push(data.Content.replace(self.user,"你"));
                    self.$forceUpdate()
                    break;
                case 3: //Play
                    var content = JSON.parse(data.Content)
                    self.holdings = _.sortBy(content.value,["color", "number"])

                    if(self.holdings.length ==14){
                        self.yourTurn = true
                    }else {
                        self.yourTurn = false
                    }
                    self.$forceUpdate()
                    break;
                case 4: //Talk
                    if (self.chats.length == 10){
                        self.chats.shift()
                    }
                    if (self.user == data.User){
                          data.Content = "你说：" + data.Content
                    }else {
                        data.Content = data.User + "说："  + data.Content
                    }
                  
                    self.chats.push(data)
                    self.$forceUpdate()
                    break;
                }       
    }

        
    },
    methods: {
        sendChat: function(){
            var talk = {
                'code':4,
                'value':this.chat
            }
            this.$socket.send(JSON.stringify(talk))
            this.chat = ""
        }
    },
    components:{
        card
    }

}
</script>
<style></style>
