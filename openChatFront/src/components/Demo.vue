<template>
  <div class="flex justify-center">
    <div id="log" class="flex-col flex-wrap w-1/2">

    </div>
  </div>
  <div class="flex justify-center py-5 mt-10">
    <div class="fixed bottom-5">
      <input autofocus type="text" id="msg" size="64" class="border-solid rounded-lg shadow-xl py-2 text-center bg-slate-100" @keypress="entryPress"/>
      <button @click="sendMessage" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-2 rounded mx-4">Send</button>
    </div>
  </div>
</template>


<script>
let msg;
let log;
let conn;

function appendLog(message){
  let item = document.createElement("div");
  let breakLine = `<div class="py-4">
        <div class="w-full border-t border-gray-300"></div>
      </div>`
  item.innerHTML =  message + breakLine;
  log.appendChild(item);
  window.scrollTo(0, document.body.scrollHeight);
}

export default{
    methods:{
      // 发送socket消息
      sendMessage(){
          if(conn.readyState === 3 || conn.readyState === 2){    // 当websocket处于关闭状态
            appendLog("<b>Connection closed.<b>")
            console.log("Connection closed")
            return false;
          }
          if(msg.value === ""){     // 不发送空字符串
            return false;
          }
          appendLog('<div class="text-right mt-5"> ' + msg.value + '</div>');   // 自己发送的消息
          conn.send(msg.value);     // 将消息发送的服务器
          msg.value = "";     // 清空聊天框
          return false;
        },
        entryPress(event){
          if(event.key === "Enter"){
            event.preventDefault();
            this.sendMessage();
          }
        },
    },
    mounted(){    // 变量初始化并设置回调函数
        msg = document.getElementById("msg");
        log = document.getElementById("log");
        conn = new WebSocket("ws://localhost:8080/ws");
        conn.onclose = function(evt){
          appendLog("<b>Connection closed.</b>");
        };
        conn.onmessage = function(evt){
            let receiveMessage = evt.data.split('\n')
            for(let i = 0; i < receiveMessage.length; i++) {
              let messageHTML = '<p class="justify-center mt-5"> ' + receiveMessage[i] +' <p>'   // 别入发送过来的消息
              appendLog(messageHTML);
            }
        }
    }
}
</script>



