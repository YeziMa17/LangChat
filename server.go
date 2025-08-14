package main //server实现
import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//增加onlineusermap
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//增加广播管道
	Message chan string
}

//创建一个server对外接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 写一个广播方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + msg

	this.Message <- sendMsg
}

// 监听message广播chan的go程，一旦有消息就发给全部在线user
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//发给全部在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Handler(conn net.Conn) {
	//当前链接的业务
	//fmt.Println("连接成功建立.")
	//用户上线了
	user := NewUser(conn)
	//将用户加到onlinemap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	//，广播消息
	this.BroadCast(user, "已上线")

	select {}
}

// 启动服务器的方法
func (this *Server) Start() {
	//socket listen
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen error:", err)
	}
	//close socket
	defer Listener.Close()

	//启动监听message的go程
	go this.ListenMessager()

	for {
		//accept
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("Listener accept error", err)
			continue
		}

		//业务回调do handler
		go this.Handler(conn)
	}

}

//windows 使用go build -o server.exe main.go server.go
//telnet localhost 8888
