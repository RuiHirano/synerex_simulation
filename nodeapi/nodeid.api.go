package nodeapi

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"os"
	"os/signal"

	"github.com/bwmarrin/snowflake"
	"google.golang.org/grpc"
)

/////////////////////////////////////
//////////   NodeAPI    ////////////
////////////////////////////////////

type NodeAPI struct {
	Node       *snowflake.Node // package variable for keeping unique ID.
	Nid        *NodeID
	Nupd       *NodeUpdate
	Numu       sync.RWMutex
	MyNodeName string
	Conn       *grpc.ClientConn
	Clt        NodeClient
	FuncSlice  []func()
}

func NewNodeAPI() *NodeAPI {
	na := &NodeAPI{
		FuncSlice: make([]func(), 0),
	}
	return na
}

// register closing functions.
func (api *NodeAPI) RegisterDeferFunction(f func()) {
	api.FuncSlice = append(api.FuncSlice, f)
}

func (api *NodeAPI) CallDeferFunctions() {
	for _, f := range api.FuncSlice {
		log.Printf("Calling %v", f)
		f()
	}
}

func (api *NodeAPI) HandleSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	for range c {
		log.Println("Signal Interrupt!")
		close(c)
	}

	api.CallDeferFunctions()

	log.Println("End at HandleSigInt in sxutil/signal.go")
	os.Exit(1)
}

// InitNodeNum for initialize NodeNum again
func (api *NodeAPI) InitNodeNum(n int) {
	var err error
	api.Node, err = snowflake.NewNode(int64(n))
	if err != nil {
		fmt.Println("Error in initializing snowflake:", err)
	} else {
		fmt.Println("Successfully Initialize node ", n)
	}
}

func (api *NodeAPI) GetNodeName(n int) string {
	ni, err := api.Clt.QueryNode(context.Background(), &NodeID{NodeId: int32(n)})
	if err != nil {
		log.Printf("Error on QueryNode %v", err)
		return "Unknown"
	}
	return ni.NodeName
}

func (api *NodeAPI) SetNodeStatus(status int32, arg string) {
	api.Numu.Lock()
	api.Nupd.NodeStatus = status
	api.Nupd.NodeArg = arg
	api.Numu.Unlock()
}

func (api *NodeAPI) startKeepAlive() {
	for {
		//		fmt.Printf("KeepAlive %s %d\n",nupd.NodeStatus, nid.KeepaliveDuration)
		time.Sleep(time.Second * time.Duration(api.Nid.KeepaliveDuration))
		if api.Nid.Secret == 0 { // this means the node is disconnected
			break
		}
		api.Numu.RLock()
		api.Nupd.UpdateCount++
		api.Clt.KeepAlive(context.Background(), api.Nupd)
		api.Numu.RUnlock()
	}
}

// RegisterNodeName is a function to register node name with node server address
func (api *NodeAPI) RegisterNodeName(nodesrv string, nm string, isServ bool) error { // register ID to server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) // insecure
	var err error
	api.Conn, err = grpc.Dial(nodesrv, opts...)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return err
	}
	//	defer conn.Close()

	api.Clt = NewNodeClient(api.Conn)
	nif := NodeInfo{
		NodeName: nm,
		IsServer: isServ,
	}
	api.MyNodeName = nm
	var ee error
	api.Nid, ee = api.Clt.RegisterNode(context.Background(), &nif)

	if ee != nil { // has error!
		log.Println("Error on get NodeID", ee)
		return ee
	} else {
		var nderr error
		api.Node, nderr = snowflake.NewNode(int64(api.Nid.NodeId))
		if nderr != nil {
			fmt.Println("Error in initializing snowflake:", err)
			return nderr
		} else {
			fmt.Println("Successfully Initialize node ", api.Nid.NodeId)
		}
	}

	api.Nupd = &NodeUpdate{
		NodeId:      api.Nid.NodeId,
		Secret:      api.Nid.Secret,
		UpdateCount: 0,
		NodeStatus:  0,
		NodeArg:     "",
	}
	//node = api.Node
	// start keepalive goroutine
	go api.startKeepAlive()
	//	fmt.Println("KeepAlive started!")
	return nil
}

// UnRegisterNode de-registrate node id
func (api *NodeAPI) UnRegisterNode() {
	log.Println("UnRegister Node ", api.Nid)
	resp, err := api.Clt.UnRegisterNode(context.Background(), api.Nid)
	api.Nid.Secret = 0
	if err != nil || !resp.Ok {
		log.Print("Can't unregister", err, resp)
	}
}
