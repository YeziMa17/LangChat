package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	//关联到server
	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user chan go程
	go user.ListenMessage()
	return user
}

// 用户的上线业务
func (this *User) Online() {
	//将用户加到onlinemap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//，广播消息
	this.server.BroadCast(this, "is online")
}

// 用户的下线业务
func (this *User) Offline() {
	//将用户从onlinemap中移除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	//，广播消息
	this.server.BroadCast(this, "is offline")
}

// 用户处理消息的业务
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// 监听当前user的channel的方法,一旦有消息就发给客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
