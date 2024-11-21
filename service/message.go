package service

import (
	"encoding/json"
	"fmt"
	"goIM/models"
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

// 读写锁
var rwLock sync.RWMutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
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
}

// sendProc 从node中获取信息并写入websocket中
func sendProc(node *Node) {
	for {
		select { 
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("websocket 写入消息失败", err)
				return
			}
			fmt.Println("websocket 发送消息成功")
		}
	}
}

// recProc 从websocket中将消息体拿出，然后进行解析，再进行信息类型判断， 最后将消息发送至目的用户的node
func recvProc(node *Node) {
	for {
		// 获取消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("websocket 读取消息失败", err)
			return
		}

		// 解析消息
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
		}
	}
}
