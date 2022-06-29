<script>
var msg
var log
var conn

function init(){
    msg = document.getElementById("msg");
    log = document.getElementById("log");
    conn = new WebSocket("ws://localhost:8080/ws");
}

// 发送socket消息
function sendMessage(){
    if(!conn){
        appendlog("<b>Connection closed.<b>")
        return false;
    }
    appendlog('<p style="color:blue;text-align: right"> ' + msg.value +' <p>');
    conn.send(msg.value);
    msg.value = "";
    return false;
}

function entryPress(event){
    if(event.key == "Enter"){
        event.preventDefault();
        sendMessage();
    }
}

function appendlog(message){
    let item = document.createElement("div");
    item.innerHTML =  message;
    log.appendChild(item);
}

function connect(){
    conn.onclose = function(evt){
        appendlog("<b>Connection closed.</b>");
    };
    conn.onmessage = function(evt){
        var recvMessage = evt.data.split('\n')
        for(var i = 0; i < recvMessage.length; i++){
            appendlog(recvMessage[i]);
        }
    }
}

export default{
    methods:{
        appendlog,
        sendMessage,
        entryPress,
    },
    mounted(){
        init();     // 变量初始化
        connect();      // 连接websockets
    }
}
</script>

<template>
    <div id="log">
    </div>
    <div id="form">
        <input type="text" id="msg" size="64" @keypress="entryPress"/>
        <input type="submit" @click="sendMessage" value="Sent" />
    </div>
</template>

<style>
#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}
#form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}
</style>

