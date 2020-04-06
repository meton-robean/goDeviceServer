package gotcp

import (
	"time"
	// "log"
	// "log/syslog"

	"net"
	"sync"
	// "sync"

	log "github.com/Sirupsen/logrus"
)

//Server ...
type Server struct {
	// 连接的客户端socket
	ClientSocket map[*net.TCPConn]*Conn

	callback ConnCallbackInterface
	Lock     sync.RWMutex
}

//NewServer ...
func NewServer(callback ConnCallbackInterface) *Server {
	return &Server{
		callback:     callback,
		ClientSocket: make(map[*net.TCPConn]*Conn),
	}
}

func (s *Server) heatBeatDeal() {
	for {
		for k, v := range s.ClientSocket {
			v.heartTimeCount++
			//60*5
			if v.heartTimeCount == 180 {
				delete(s.ClientSocket, k)
				v.Close()
				log.Debug("心跳超时断开:", v.gateWayID)
			} else {
				s.ClientSocket[k] = v
			}
		}
		//log.Info("len:", len(s.ClientSocket))
		time.Sleep(1 * time.Second)
	}
}

//StartServer 开始服务
func (s *Server) StartServer(addr string, flag string, callback GatewayCallBack) error {
	var err error
	var tcpAddr *net.TCPAddr
	tcpAddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	var listener *net.TCPListener
	listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	go s.heatBeatDeal()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Error("accpet Error ")
			continue
		}

		connVal := s.StoreClientSocket(conn, callback)
		log.Info("new connection:", connVal.GetRemoteAddr(), ",flag1:", flag)
		connVal.SetClientFlag(flag)
		go connVal.Do()
	}
}

//DeleteClientSocket ...
func (s *Server) DeleteClientSocket(conn *net.TCPConn) {
	s.Lock.Lock()
	delete(s.ClientSocket, conn)
	s.Lock.Unlock()
}

//StoreClientSocket ...
func (s *Server) StoreClientSocket(conn *net.TCPConn, callback GatewayCallBack) *Conn {
	s.Lock.Lock()
	nconn := newConn(conn, s, callback)
	s.ClientSocket[conn] = nconn
	s.Lock.Unlock()
	return nconn
}

//StopServer 停止服务
func (s *Server) StopServer() {
	s.callback.Close()
}
