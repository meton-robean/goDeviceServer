package gotcp

import (

	//"log"

	"net"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

//SendChanSize ...
const SendChanSize int = 10
const ReceiveChanSize int = 10

type ConnCallbackInterface interface {
	HandleMsg(*Conn, []byte) error
	Close()
}

//GatewayCallBack 网关断开的回调函数
type GatewayCallBack func(gatewayID string) error

type Conn struct {
	Srv            *Server
	conn           *net.TCPConn
	Lock           sync.RWMutex
	closeOnce      sync.Once
	closeChan      chan struct{}
	ReceiveChan    chan []byte
	SendChan       chan []byte
	heartTimeCount int

	//标识，哪个客户端
	clientFlag string

	//网关断开回调
	CallBack  GatewayCallBack
	gateWayID string
}

func newConn(conn *net.TCPConn, srv *Server, callBack GatewayCallBack) *Conn {
	return &Conn{
		Srv:            srv,
		conn:           conn,
		closeChan:      make(chan struct{}),
		SendChan:       make(chan []byte, SendChanSize),
		ReceiveChan:    make(chan []byte, ReceiveChanSize),
		heartTimeCount: 0,
		CallBack:       callBack,
	}
}

func (c *Conn) SetClientFlag(flag string) {
	c.clientFlag = flag
}

func (c *Conn) GetRawConn() *net.TCPConn {
	return c.conn
}

func (c *Conn) GetListenAddr() string {
	return c.conn.LocalAddr().String()
}

func (c *Conn) GetRemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Conn) SetGatwayID(gwID string) {
	for k, v := range c.Srv.ClientSocket {
		if v.gateWayID == gwID {
			c.Srv.ClientSocket[k].gateWayID = ""
		}
	}
	c.gateWayID = gwID
}

func (c *Conn) Close() {
	log.Debug("Close socket")
	c.closeOnce.Do(func() {
		//c.Srv.DeleteClientSocket(c.conn)
		close(c.closeChan)
		close(c.ReceiveChan)
		close(c.SendChan)
		c.conn.Close()
		c.CallBack(c.gateWayID)
	})
}

func (c *Conn) Do() {
	go c.readLoop()
	go c.handleLoop()
	go c.writeLoop()
}

func (c *Conn) updateHeartTimer() {
	c.heartTimeCount = 0
}

func (c *Conn) readLoop() {
	defer func() {
		log.Debug("readLoop end:", c.gateWayID)
		//recover()
		c.Close()
	}()
	for {
		select {
		case <-c.closeChan:
			return
		default:
			//c.conn.Write([]byte("hb"))
			c.updateHeartTimer()
			var buf = make([]byte, 1024)
			len, err := c.conn.Read(buf)
			if err != nil {
				if !strings.Contains(err.Error(), "EOF") {
					log.Error("strFlag:", c.clientFlag, ",read pack Eror:", err, c.GetRemoteAddr(), ",gwID:", c.gateWayID)
				}
				log.Error("err:", err, ",addr:", c.GetRemoteAddr(), ",gwID:", c.gateWayID)
				return
			}
			//log.Debug("len:", len, ",buf:", string(buf[:len]))
			c.ReceiveChan <- buf[:len]
		}
	}
}

func (c *Conn) writeLoop() {
	defer func() {
		log.Debug("writeLoop end:", c.gateWayID)
		//recover()
		c.Close()
	}()
	for {
		select {
		case <-c.closeChan:
			return
		case data := <-c.SendChan:
			if _, err := c.conn.Write(data); err != nil {
				log.Error("clientFlag:", c.clientFlag, ",conn write err: ", err, c.GetRemoteAddr(), ",gwID:", c.gateWayID)
				return
			}
		}
	}
}

func (c *Conn) handleLoop() {
	defer func() {
		log.Debug("handleLoop end:", c.gateWayID)
		//recover()
		c.Close()
	}()

	for {
		select {
		case <-c.closeChan:
			return
		case p := <-c.ReceiveChan:
			if p != nil {
				c.Srv.callback.HandleMsg(c, p)
			}
		}
	}
}

func (c *Conn) SetKeepAlivePeriod(d time.Duration) {
	c.conn.SetKeepAlive(true)
	c.conn.SetKeepAlivePeriod(d)
}
