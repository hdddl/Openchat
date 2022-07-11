window.onload = function (){
	let conn = new WebSocket("ws://" + document.location.host + "/ws")
	let msg = document.getElementById("msg");
	let log = document.getElementById("log");

	// 显示消息
	function appendLog(message){
		let item = document.createElement("div");
		let breakLine = `<div class="py-4">
                <div class="w-full border-t border-gray-300"></div>
              </div>`       // 用一个横行隔开不同的消息
		item.innerHTML =  message + breakLine;      // 构建需要插入的元素
		log.appendChild(item);      // 添加子元素
		window.scrollTo(0, document.body.scrollHeight);     // 保持滚动
	}

	// 读取input里面的内容，并发送给服务器
	function sendMessage() {
		if (conn.readyState === 3 || conn.readyState === 2) {    // 当websocket处于关闭状态
			appendLog("<b>Connection closed.<b>")
			console.log("Connection closed")
			return false;
		}
		if (msg.value === "") {     // 不发送空字符串
			return false;
		}
		appendLog('<div class="text-right mt-5"> ' + msg.value + '</div>');   // 自己发送的消息
		conn.send(msg.value);     // 将消息发送的服务器
		msg.value = "";     // 清空聊天框
		return false;
	}
	// 注册消息接受与发送函数
	conn.onclose = function (e){
		appendLog("<b>Connection closed.</b>")
	};
	conn.onmessage = function (evt){
		let receiveMessage = evt.data.split('\n')
		for(let i = 0; i < receiveMessage.length; i++){
			let messageHTML = '<p class="justify-center mt-5"> ' + receiveMessage[i] + '</p>'
			appendLog(messageHTML)
		}
	}

	// 消息发送动作处理
	document.getElementById("submit").onclick = function (){
		sendMessage()
	}
	msg.addEventListener("keypress", function (ev){
		if(ev.key === 'Enter'){
			ev.preventDefault();
			sendMessage();
		}
	})

}