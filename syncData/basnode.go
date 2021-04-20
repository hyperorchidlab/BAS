package syncData

import (
	"encoding/json"
	"errors"
	"github.com/hyperorchidlab/pirate_contract/util"
	"net"
	"strconv"
	"sync"
	"time"
)

const(
	RcvBufLen int = 1 << 21
)

type Node struct {
	IPAddr string
	Port int
}

func (n *Node)DialWithTimeOut(timeout int) (net.Conn, error)  {
	dialer:=net.Dialer{
		Timeout: time.Second*2,
	}

	addr:=n.IPAddr + ":" + strconv.Itoa(n.Port)

	var (
		err error
		conn net.Conn
	)

	conn,err = dialer.Dial("udp4",addr)
	if err!=nil{
		return nil,err
	}
	defer conn.Close()

	if timeout > 0{
		conn.SetDeadline(time.Now().Add(time.Second*time.Duration(timeout)))
	}

	return conn,nil
}

func (n *Node)Send(msg interface{}) error  {

	conn,err:=n.DialWithTimeOut(0)
	if err!=nil{
		return err
	}

	defer conn.Close()

	j,_:=json.Marshal(msg)

	_, err = conn.Write(j)
	if err!=nil{
		return err
	}

	return nil
}

func (n *Node)SendSync(msg interface{}) ([]byte,error)  {
	conn,err:=n.DialWithTimeOut(4)
	if err!=nil{
		return nil,err
	}

	defer conn.Close()


	j,_:=json.Marshal(msg)

	_,err = conn.Write(j)
	if err!=nil{
		return nil,err
	}

	var nr int

	buf:=make([]byte,RcvBufLen)
	nr,err = conn.Read(buf)
	if err!=nil{
		return nil,err
	}

	return buf[:nr],nil
}


type NodeStatus struct {
	Node
	LastAccessTime int64
    CurrentVersion int64
}


type NodeGroup struct {
	lock sync.Mutex
	nodes []*NodeStatus
	self  *NodeStatus
}

func (ng *NodeGroup)Init(addr string, port int, version int64)  {
	ns:=&NodeStatus{
		Node:Node{
			IPAddr: addr,
			Port: port,
		},
		LastAccessTime: util.GetNowMsTime(),
		CurrentVersion: version,
	}

	ng.self = ns
}

func (ng *NodeGroup)UpdateSelfVersion(version int64)  {
	ng.lock.Lock()
	defer ng.lock.Unlock()

	ng.self.CurrentVersion = version
	ng.self.LastAccessTime = util.GetNowMsTime()
}

func (ng *NodeGroup)add(addr string, port int)  {
	n:=&Node{
		IPAddr: addr,
		Port: port,
	}

	ns:=&NodeStatus{
		Node:*n,
	}

	ng.nodes = append(ng.nodes,ns)
}

func (ng *NodeGroup)Add(addr string, port int)  error {
	ng.lock.Lock()
	defer ng.lock.Unlock()

	for i:=0;i<len(ng.nodes);i++{
		n:=ng.nodes[i]
		if n.Port == port && n.IPAddr == addr{
			return errors.New("duplicate node")
		}
	}

	ng.add(addr,port)

	return nil
}

func (ng *NodeGroup)Del(addr string, port int)  {
	ng.lock.Lock()
	defer ng.lock.Unlock()

	idx := -1
	l:=len(ng.nodes)
	for i:=0;i<l;i++{
		n:=ng.nodes[i]
		if n.Port == port && n.IPAddr == addr{
			idx = i
		}
	}

	if idx != -1{
		ng.nodes[idx]=ng.nodes[l-1]
		ng.nodes = ng.nodes[:l-1]
	}

	return
}

func (ng *NodeGroup)Update(addr string, port int, version int64)  {
	ng.lock.Lock()
	defer ng.lock.Unlock()

	for i:=0;i<len(ng.nodes);i++{
		n:=ng.nodes[i]
		if n.Port == port && n.IPAddr == addr{
			n.LastAccessTime = util.GetNowMsTime()
			n.CurrentVersion = version
		}
	}
}

func (ng *NodeGroup)FindVersion(addr string, port int) (int64, int64, error)  {
	ng.lock.Lock()
	defer ng.lock.Unlock()

	for i:=0;i<len(ng.nodes);i++{
		n:=ng.nodes[i]
		if n.Port == port && n.IPAddr == addr{
			return n.CurrentVersion,n.LastAccessTime,nil
		}
	}

	return 0,0,errors.New("not found node")
}

func (ng *NodeGroup)BroadCast(msg interface{})  {
	ng.lock.Lock()
	defer ng.lock.Unlock()

	for i:=0;i<len(ng.nodes);i++{
		go ng.nodes[i].Send(msg)
	}
}

