// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"openChat/utils"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

// token: Hub
var hubs = map[string]*Hub{}

const tokenLength = 8

// serveWs 处理Websocket连接
func serveWs(ctx *gin.Context) {
	token := ctx.Query("token")
	hub, ok := hubs[token]
	if !ok {
		ctx.Status(http.StatusNotAcceptable)
		return
	}
	if hub == nil || !hub.isRunning {
		ctx.Status(http.StatusNotAcceptable) // 该聊天框已经被关闭
		return
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan messageType, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// 处理聊天框
func chatFrame(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" { // 如果token不存在则新建一个chat frame
		token = utils.RandStringRunes(tokenLength)
		_, ok := hubs[token]
		if ok { // 如果生成的token已经存在则重新生成
			token = utils.RandStringRunes(tokenLength)
			_, ok = hubs[token]
		}
		hub := newHub()
		hub.token = token
		hubs[token] = hub
		go hub.run()
		ctx.Redirect(302, "/chatFrame?token="+token) // 网页重定向
		return
	}
	ctx.File("./frontEnd/src/chatFrame.html") // 返回前端页面
}

// 获取在线人数
func onlineNumber(ctx *gin.Context) {
	for _, hub := range hubs {
		log.Println(len(hub.clients)) // 打印出各hub的在线人数
	}
}

// 路由注册
func register(router *gin.Engine) {
	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./frontEnd/src/404NotFound.html")
	}) // 处理404页面

	router.GET("/", func(ctx *gin.Context) { // 处理首页
		ctx.File("./frontEnd/src/index.html")
	}) // 处理首页

	router.GET("/chatFrame", chatFrame) // 处理聊天框
	router.GET("/ws", serveWs)          // 处理websocket

	router.StaticFS("/dist/", http.Dir("./frontEnd/dist")) // 处理静态文件
	router.GET("/api", onlineNumber)
}

func main() {
	flag.Parse()
	router := gin.Default()
	register(router) // 注册路由
	log.Printf("open chat start successful, address is %s", *addr)
	err := router.Run() // 启动web服务器
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
