<template>
<div>
    <a class="btn btn-default" @click="play(card)" v-bind:disabled="yourTurn" v-for=" card in value">{{card.showing}}</a>
    </div>
</template>
<script>
export default {
    name :'card',
    data(){
        return{
        }
    },
    props:['value','turn'],
    methods:{
        play: function(card){
            if(this.turn){
                var play = {
                'code':3,
                'value': JSON.stringify(card)
            }
            console.info("user played" + play)
            this.$socket.send(JSON.stringify(play))    
            }
            
        }
    },
    computed:{
        yourTurn: function(){
            if(this.turn){
                return false
            }else {
                return "disabled"
            }
        }
    }
}
</script>
<style>
</style>
