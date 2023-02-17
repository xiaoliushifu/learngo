package main

/**
一个简单的基于TCP的聊天室
四种goroutine:
	一个main：goroutine
	一个广播消息：goroutine
	一个读取goroutine （每个客户端连接一个）
	一个写goroutine （每个客户端连接一个）

三个通道变量
	一个存储全局（在线）用户列表的通道
	一个存储消息的通道
	一个离开消息的通道

每来一个客户端连接，都需要创建一个读和写goroutine
*/
import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

//用户对象
type User struct {
	ID             int
	Addr           string
	EnterAt        time.Time
	MessageChannel chan string
}

func (u User) String() string {
	return strconv.Itoa(u.ID)
}

func GenUserID() int {
	return rand.Intn(10)
}

var (
	// 全局用户列表，用channel变量
	enteringChannel = make(chan *User)
	//离开用户变量，写入消息
	leavingChannel = make(chan *User)
	//服务端消息通道，用来给所有在线用户广播
	messageChannel = make(chan string, 8)
)

func main() {
	listener, err := net.Listen("tcp", ":2020")
	if err != nil {
		panic(err)
	}

	go broadcaster()

	for {
		//主goroutine阻塞在这，等待客户端的连接到来
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//每来一个连接，就对应创建一个用户
		go handleConn(conn)
	}
}

//每个进来的连接，创建一个goroutine，用来读取专属这个用户的连接
func handleConn(conn net.Conn) {
	defer conn.Close()
	user := &User{
		ID:             GenUserID(),
		Addr:           conn.RemoteAddr().String(),
		EnterAt:        time.Now(),
		MessageChannel: make(chan string, 8),
	}

	//服务端向该用户发送消息，也要启动专门的写goroutine
	go sendMessage(conn, user.MessageChannel)

	user.MessageChannel <- "Welcome, " + user.String()
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has enter"

	//加入全局用户在线列表，用channel作为这个所有goroutine可以共享的变量，非常合适
	enteringChannel <- user

	//服务端始终阻塞接收当前用户发送的消息
	input := bufio.NewScanner(conn)
	for input.Scan() { //直到false，才退出循环
		messageChannel <- strconv.Itoa(user.ID) + ":" + input.Text()
	}

	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}
	//用户离开
	leavingChannel <- user
	messageChannel <- "user: `" + strconv.Itoa(user.ID) + "` has left"
}

//每个用户（连接）专属的写goroutine
func sendMessage(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

//服务端处理沟通各个连接的一个专属goroutine
//负责各个goroutine的交互，就是通道变量的操作
func broadcaster() {
	//一个map
	users := make(map[*User]struct{})
	for {
		select {
		case user := <-enteringChannel:
			//新用户进入
			users[user] = struct{}{} //注意这种类似于interface{}的写法：临时的空结构体类型 struct{},后面的{}是实例化
		case user := <-leavingChannel:
			delete(users, user)
			close(user.MessageChannel) //关闭用户身上的通道，避免goroutine泄漏？
		case msg := <-messageChannel:
			for user := range users {
				user.MessageChannel <- msg //给所有用户发消息，当然也包含该消息的原始发送者
			}
		}
	}

}
