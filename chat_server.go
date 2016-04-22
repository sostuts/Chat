package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Chaters map[net.Conn]string
}

func main() {
	//本机ip
	ip := "210.38.196.xxx"
	Start_server(ip)
}

/** 开启服务
*	@params    port  		   string	监听端口
*
**/
func Start_server(ip string) {
	//初始化
	chaters := Server{make(map[net.Conn]string, 0)}

	//监听
	conn, err := net.Listen("tcp", ip+":9988")
	chaters.Check_err(err, nil)
	fmt.Println("Start Ok!")
	for {
		//允许握手
		c, err := conn.Accept()
		chaters.Check_err(err, nil)

		//添加线程
		go chaters.Server_chat(c)
	}
}

/** 检查错误
*	@params    err			   error
*	@params    conn  		   net.Conn  	聊天线程
*
*	@return    string 		   "break"   	判断有没人退出
**/
func (server *Server) Check_err(err error, conn net.Conn) string {
	if err != nil {
		if conn != nil {
			//有人退出
			return "break"
		} else {
			//error
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
	return ""
}

/** 监听多线程聊天
*
*	@params    conn  		   net.Conn  	聊天线程
**/
func (server *Server) Server_chat(conn net.Conn) {
	//第一次加入
	server.Chat(true, conn)

	for {
		//重复监听，判断是否有人退出
		if server.Chat(false, conn) == "break" {
			break
		}
	}
}

/** 多线程聊天
*	@params    bool   mode   		(true为第一次加入)
*	@params    conn   net.Conn  	聊天线程
*
*	@return    string "break"   判断有没人退出
**/
func (server *Server) Chat(mode bool, conn net.Conn) string {
	//接收信息
	bytes, err := server.Receive(mode, conn)

	//有人退出
	if server.Check_err(err, conn) == "break" {
		//broadcast
		server.Broadcast([]byte(server.Chaters[conn]+" OUT!"), conn)
		//delete the map
		delete(server.Chaters, conn)
		conn.Close()
		return "break"
	}

	//广播
	broadcast_text := server.Handle_message(mode, bytes, conn)
	server.Broadcast(broadcast_text, conn)
	return ""
}

/** 整理广播信息
*	@params    mode   		bool   		(true为第一次加入)
*	@params    []byte   	bytes   	服务器接收的信息
*	@params    conn  		net.Conn  	聊天线程
*
*	@return    []byte
**/
func (server *Server) Handle_message(mode bool, bytes []byte, conn net.Conn) (broadcast_text []byte) {
	if mode == true {
		//第一次加入
		broadcast_text = []byte(string(bytes) + " IN !")
		//add to map
		server.Chaters[conn] = string(bytes)
	} else {
		//say something
		broadcast_text = append([]byte(server.Chaters[conn]+" Say: "), bytes...)
	}
	return
}

/** 接收信息
*	@params    mode   	bool   		(true为第一次加入)
*	@params    conn  	net.Conn  	聊天线程
*
*	@return    []byte && error
**/
func (server *Server) Receive(mode bool, conn net.Conn) (bytes []byte, err error) {
	if mode == true {
		//姓名10字节
		bytes = make([]byte, 10)
	} else {
		//内容256字节
		bytes = make([]byte, 256)
	}

	//接收信息
	_, err = conn.Read(bytes)
	return
}

/** 广播
*	@params    broadcast_text  []byte  		广播内容
*	@params    conn  		   net.Conn  	聊天线程
*
**/
func (server *Server) Broadcast(broadcast_text []byte, conn net.Conn) {
	if len(broadcast_text) == 0 {
		return
	}
	for key, _ := range server.Chaters {
		//跳过自己
		if key == conn {
			continue
		}
		_, err := key.Write(broadcast_text)
		server.Check_err(err, nil)
	}
}
