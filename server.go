package main //server实现
import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

//创建一个server对外接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) Handler(conn net.Conn) {
	//当前链接的业务
	fmt.Println("连接成功建立.")
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
