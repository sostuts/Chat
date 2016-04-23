package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

type Conn struct {
	conn net.Conn
}

func main() {
	//获取IP
	ip := Input_Ip()
	//连接服务器
	conn := Connect(ip + ":9988")
	//begin
	go conn.Receive()
	conn.Chat()
}

//输入服务器IP
func Input_Ip() string {
	for {
		fmt.Println("Please input server IP .")
		line, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		//正确的IP格式才返回值
		if net.ParseIP(string(line)) != nil {
			return string(line)
		}
	}
}

//检查错误
func Check_err(err error) {
	if err != nil {
		fmt.Println(err.Error())
		bufio.NewReader(os.Stdin).ReadLine()
		os.Exit(1)
	}
}

//连接服务器
func Connect(ip string) *Conn {
	//30s则连接超时
	c, err := net.DialTimeout("tcp", ip, 30*time.Second)
	Check_err(err)
	return &Conn{c}
}

//输入以及发送信息
func (c *Conn) Chat() {
	var Reader = func(mode bool) {
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		Check_err(err)
		if mode {
			//加密后发送
			c.conn.Write(Encypt(line, []byte("www.zeffee.com")))
		} else {
			c.conn.Write(line)
		}
	}
	fmt.Println("Who are u?")
	Reader(false)
	fmt.Println("Welcome!")
	for {
		Reader(true)
	}
}

//接收信息
func (c *Conn) Receive() {
	for {
		bytes := make([]byte, 256)
		_, err := c.conn.Read(bytes)
		Check_err(err)
		//解密之后输出,不解密前16个字符
		fmt.Println(string(bytes[:16]) + Decrypt(bytes[16:], []byte("www.zeffee.com")))
	}
}

//加密信息
func Encypt(source, the_key []byte) (mcrypt []byte) {
	today := time.Now().Day()
	//添加文本总长度
	mcrypt = append(mcrypt, byte(len(source)))
	for i, val := range source {
		if i%2 == 0 {
			mcrypt = append(mcrypt, val^byte(today))
		} else {
			mcrypt = append(mcrypt, val^the_key[i%len(the_key)])
		}
	}

	return mcrypt
}

//解密信息
func Decrypt(mcrypt, the_key []byte) (source string) {
	today := time.Now().Day()
	//获取文本总长度
	count := int(mcrypt[0])
	for i, val := range mcrypt[1:] {
		if i >= count {
			break
		}
		if i%2 == 0 {
			source += string(val ^ byte(today))
		} else {
			source += string(val ^ the_key[i%len(the_key)])
		}
	}
	return source
}
