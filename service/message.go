package service

import (
	"encoding/json"
	"fmt"
	"goIM/dao"
	"goIM/models"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
)

type messageService struct{}

var MessageService = new(messageService)

type Node struct {
	Conn      *websocket.Conn
	Addr      string
	DataQueue chan []byte
	GroupSet  set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 全局channel
var globalSendChan chan []byte = make(chan []byte, 1024)

// 读写锁
var rwLock sync.RWMutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

func (messageService *messageService) Chat(ctx *gin.Context) {
	userIdStr := ctx.Query("userId")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Info("ownerId 类型转换失败", err)
		jsonOutPut(ctx, -1, "参数错误")
		return
	}

	// 升级websocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		zap.S().Info("升级协议失败", err)
		jsonOutPut(ctx, -1, "链接失败请重试!")
		return
	}

	// 构造节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSet:  set.New(set.ThreadSafe),
	}

	rwLock.Lock()
	clientMap[userId] = node
	rwLock.Unlock()

	// 开协程 收发服务器消息
	go sendProc(node)

	go recvProc(node)

	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

// sendProc 从node中获取信息并写入websocket中
func sendProc(node *Node) {

	/* for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("websocket 写入消息失败", err)
				return
			}
			fmt.Println("websocket 发送消息成功")
		}
	} */

	for data := range node.DataQueue {
		err := node.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			zap.S().Info("websocket 写入消息失败", err)
			return
		}
		fmt.Println("websocket 发送消息成功")
	}
}

// recvProc 从websocket中将消息体拿出，然后进行解析，再进行信息类型判断， 最后将消息发送至目的用户的node
// recvProc 逻辑调整 取出消息不做消息业务处理 做广播 放入全局管道 通过udp广播来增加吞吐
func recvProc(node *Node) {
	for {
		// 获取消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("websocket 读取消息失败", err)
			return
		}

		/* // 解析消息
		msg := models.Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			zap.S().Info("消息解析失败", err)
			return
		}
		fmt.Println("消息解析: ", msg)

		switch msg.Type {
		case 1:
			zap.S().Info("这是一条私信：", msg.Content)
			targetNode, ok := clientMap[msg.TargetId]
			if !ok {
				zap.S().Info("不存在对应的玩家Node: ", msg.TargetId)
				continue
			}
			targetNode.DataQueue <- data
			fmt.Println("发送成功：", string(data))
		} */

		// 将消息放入全局channel中
		broadCastMsg(data)
	}
}

func broadCastMsg(data []byte) {
	globalSendChan <- data
}

// UdpSendProc 完成udp数据发送 将受到的全局channel广播到全部服务器
func udpSendProc() {
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 255),
		Port: 3000,
	})

	if err != nil {
		zap.S().Info("拨号UDP端口失败", err)
		return
	}
	defer udpConn.Close()

	/* for {
		select {
		case data := <-globalSendChan:
			_, err := udpConn.Write(data)
			if err != nil {
				zap.S().Info("写入Udp消息失败", err)
				return
			}
		}
	} */

	for data := range globalSendChan {
		_, err := udpConn.Write(data)
		if err != nil {
			zap.S().Info("写入Udp消息失败", err)
			return
		}
	}
}

// UdpRecvProc 完成udp数据的接受 启动udp服务端 获取udp客户端写入的消息
func udpRecvProc() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 255),
		Port: 3000,
	})
	if err != nil {
		zap.S().Info("监听udp端口失败", err)
		return
	}
	defer udpConn.Close()

	for {
		var buff [1024]byte
		n, err := udpConn.Read(buff[0:])
		if err != nil {
			zap.S().Info("读取udp数据失败", err)
			return
		}

		// 处理消息发送逻辑
		dispatchMsg(buff[0:n])
	}
}

func dispatchMsg(data []byte) {
	// 解析消息
	msg := models.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("消息解析失败", err)
		return
	}
	fmt.Println("消息解析: ", msg)
	fmt.Println("解析数据:", msg, "msg.FormId", msg.FromId, "targetId:", msg.TargetId, "type:", msg.Type)

	// 判断消息类型
	switch msg.Type {
	case 1:
		zap.S().Info("这是一条私信：", msg.Content)
		sendMsg(msg.TargetId, data)
		zap.S().Info("发送成功：", string(data))
	case 2:
		sendGroupMsg(uint(msg.FromId), uint(msg.TargetId), data)
	}
}

func sendMsg(targetId int64, data []byte) {
	targetNode, ok := clientMap[targetId]
	if !ok {
		zap.S().Info("不存在对应的玩家Node: ", targetId)
		return
	}
	targetNode.DataQueue <- data
}

func sendGroupMsg(fromId, groupId uint, data []byte) (int, error) {
	userIds, err := dao.GetGroupUsers(groupId)
	if err != nil {
		return -1, err
	}
	for _, userId := range *userIds {
		if userId == fromId {
			continue
		}
		sendMsg(int64(userId), data)
	}

	return 0, nil
}
