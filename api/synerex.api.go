package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"os"
	"os/signal"

	"github.com/google/uuid"

	"github.com/bwmarrin/snowflake"
	"github.com/golang/protobuf/ptypes"
	"github.com/synerex/synerex_alpha/nodeapi"
	"google.golang.org/grpc"
)

// IDType for all ID in Synergic Exchange
type IDType uint64

var (
	node       *snowflake.Node // package variable for keeping unique ID.
	nid        *nodeapi.NodeID
	nupd       *nodeapi.NodeUpdate
	numu       sync.RWMutex
	myNodeName string
	conn       *grpc.ClientConn
	clt        nodeapi.NodeClient
	funcSlice  []func()
)

// DemandOpts is sender options for Demand
type DemandOpts struct {
	ID        uint64
	Target    uint64
	Name      string
	JSON      string
	SimDemand *SimDemand
}

// SupplyOpts is sender options for Supply
type SupplyOpts struct {
	ID        uint64
	Target    uint64
	Name      string
	JSON      string
	SimSupply *SimSupply
}

func init() {
	fmt.Println("Synergic Exchange Util init() is called!")
	funcSlice = make([]func(), 0)
}

/////////////////////////////////////
//////////   NodeAPI    ////////////
////////////////////////////////////

type NodeAPI struct {
	Node       *snowflake.Node // package variable for keeping unique ID.
	Nid        *nodeapi.NodeID
	Nupd       *nodeapi.NodeUpdate
	Numu       sync.RWMutex
	MyNodeName string
	Conn       *grpc.ClientConn
	Clt        nodeapi.NodeClient
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

	CallDeferFunctions()

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
	ni, err := api.Clt.QueryNode(context.Background(), &nodeapi.NodeID{NodeId: int32(n)})
	if err != nil {
		log.Printf("Error on QueryNode %v", err)
		return "Unknown"
	}
	return ni.NodeName
}

func (api *NodeAPI) SetNodeStatus(status int32, arg string) {
	numu.Lock()
	api.Nupd.NodeStatus = status
	api.Nupd.NodeArg = arg
	numu.Unlock()
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

	api.Clt = nodeapi.NewNodeClient(api.Conn)
	nif := nodeapi.NodeInfo{
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

	api.Nupd = &nodeapi.NodeUpdate{
		NodeId:      api.Nid.NodeId,
		Secret:      api.Nid.Secret,
		UpdateCount: 0,
		NodeStatus:  0,
		NodeArg:     "",
	}
	node = api.Node
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

///////////////////////////////////////////////

// register closing functions.
func RegisterDeferFunction(f func()) {
	funcSlice = append(funcSlice, f)
}

func CallDeferFunctions() {
	for _, f := range funcSlice {
		log.Printf("Calling %v", f)
		f()
	}
}

func HandleSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	for range c {
		log.Println("Signal Interrupt!")
		close(c)
	}

	CallDeferFunctions()

	log.Println("End at HandleSigInt in sxutil/signal.go")
	os.Exit(1)
}

// InitNodeNum for initialize NodeNum again
func InitNodeNum(n int) {
	var err error
	node, err = snowflake.NewNode(int64(n))
	if err != nil {
		fmt.Println("Error in initializing snowflake:", err)
	} else {
		fmt.Println("Successfully Initialize node ", n)
	}
}

func GetNodeName(n int) string {
	ni, err := clt.QueryNode(context.Background(), &nodeapi.NodeID{NodeId: int32(n)})
	if err != nil {
		log.Printf("Error on QueryNode %v", err)
		return "Unknown"
	}
	return ni.NodeName
}

func SetNodeStatus(status int32, arg string) {
	numu.Lock()
	nupd.NodeStatus = status
	nupd.NodeArg = arg
	numu.Unlock()
}

func startKeepAlive() {
	for {
		//		fmt.Printf("KeepAlive %s %d\n",nupd.NodeStatus, nid.KeepaliveDuration)
		time.Sleep(time.Second * time.Duration(nid.KeepaliveDuration))
		if nid.Secret == 0 { // this means the node is disconnected
			break
		}
		numu.RLock()
		nupd.UpdateCount++
		clt.KeepAlive(context.Background(), nupd)
		numu.RUnlock()
	}
}

// RegisterNodeName is a function to register node name with node server address
func RegisterNodeName(nodesrv string, nm string, isServ bool) error { // register ID to server
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) // insecure
	var err error
	conn, err = grpc.Dial(nodesrv, opts...)
	if err != nil {
		log.Printf("fail to dial: %v", err)
		return err
	}
	//	defer conn.Close()

	clt = nodeapi.NewNodeClient(conn)
	nif := nodeapi.NodeInfo{
		NodeName: nm,
		IsServer: isServ,
	}
	myNodeName = nm
	var ee error
	nid, ee = clt.RegisterNode(context.Background(), &nif)

	if ee != nil { // has error!
		log.Println("Error on get NodeID", ee)
		return ee
	} else {
		var nderr error
		node, nderr = snowflake.NewNode(int64(nid.NodeId))
		if nderr != nil {
			fmt.Println("Error in initializing snowflake:", err)
			return nderr
		} else {
			fmt.Println("Successfully Initialize node ", nid.NodeId)
		}
	}

	nupd = &nodeapi.NodeUpdate{
		NodeId:      nid.NodeId,
		Secret:      nid.Secret,
		UpdateCount: 0,
		NodeStatus:  0,
		NodeArg:     "",
	}
	// start keepalive goroutine
	go startKeepAlive()
	//	fmt.Println("KeepAlive started!")
	return nil
}

// UnRegisterNode de-registrate node id
func UnRegisterNode() {
	log.Println("UnRegister Node ", nid)
	resp, err := clt.UnRegisterNode(context.Background(), nid)
	nid.Secret = 0
	if err != nil || !resp.Ok {
		log.Print("Can't unregister", err, resp)
	}
}

// SMServiceClient Wrappter Structure for market client
type SMServiceClient struct {
	ClientID   IDType
	ProviderID uint64
	MType      ChannelType
	Client     SynerexClient
	ArgJson    string
	MbusID     IDType
}

// NewSMServiceClient Creates wrapper structre SMServiceClient from SynerexClient
func NewSMServiceClient(clt SynerexClient, mtype ChannelType, providerID uint64, argJson string) *SMServiceClient {
	uid, _ := uuid.NewRandom()
	s := &SMServiceClient{
		ClientID:   IDType(uid.ID()),
		ProviderID: providerID,
		MType:      mtype,
		Client:     clt,
		ArgJson:    argJson,
	}
	return s
}

// GenerateIntID for generate uniquie ID
func GenerateIntID() uint64 {
	uid, _ := uuid.NewRandom()
	return uint64(uid.ID())
}

func (clt SMServiceClient) getChannel() *Channel {
	return &Channel{ClientId: uint64(clt.ClientID), Type: clt.MType, ProviderId: clt.ProviderID, ArgJson: clt.ArgJson}
}

// IsSupplyTarget is a helper function to check target
func (clt *SMServiceClient) IsSupplyTarget(sp Supply, idlist []uint64) bool {
	spid := sp.TargetId
	for _, id := range idlist {
		if id == spid {
			return true
		}
	}
	return false
}

// IsDemandTarget is a helper function to check target
func (clt *SMServiceClient) IsDemandTarget(dm Demand, idlist []uint64) bool {
	dmid := dm.TargetId
	for _, id := range idlist {
		if id == dmid {
			return true
		}
	}
	return false
}

// ProposeSupply send proposal Supply message to server
func (clt *SMServiceClient) ProposeSupply(spo *SupplyOpts) uint64 {
	pid := GenerateIntID()
	ts := ptypes.TimestampNow()
	sp := &Supply{
		Id:         pid,
		SenderId:   uint64(clt.ClientID),
		TargetId:   spo.Target,
		Type:       clt.MType,
		Ts:         ts,
		SupplyName: spo.Name,
		ArgJson:    spo.JSON,
	}

	if spo.SimSupply != nil {
		sp.WithSimSupply(spo.SimSupply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := clt.Client.ProposeSupply(ctx, sp)
	//log.Printf("%v.Test err %v, [%v]", clt, err, sp)
	if err != nil {
		log.Printf("%v.ProposeSupply err %v, [%v]", clt, err, sp)
		return 0 // should check...
	}
	return pid
}

// SelectSupply send select message to server
func (clt *SMServiceClient) SelectSupply(sp Supply) (uint64, error) {
	tgt := &Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: sp.Id, /// Message Id of Supply (not SenderId),
		Type:     sp.Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.SelectSupply(ctx, tgt)
	if err != nil {
		log.Printf("%v.SelectSupply err %v", clt, err)
		return 0, err
	}
	log.Println("SelectSupply Response:", resp)
	// if mbus is OK, start mbus!
	clt.MbusID = IDType(resp.MbusId)
	if clt.MbusID != 0 {
		//		clt.SubscribeMbus()
	}

	return uint64(clt.MbusID), nil
}

// SelectDemand send select message to server
func (clt *SMServiceClient) SelectDemand(dm Demand) error {
	tgt := &Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: dm.Id,
		Type:     dm.Type,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.SelectDemand(ctx, tgt)
	if err != nil {
		log.Printf("%v.SelectDemand err %v", clt, err)
		return err
	}
	log.Println("SelectDemand Response:", resp)
	return nil
}

// SubscribeSupply  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeSupply(ctx context.Context, spcb func(*SMServiceClient, *Supply)) error {
	ch := clt.getChannel()
	smc, err := clt.Client.SubscribeSyncSupply(ctx, ch)
	//log.Printf("Test3 %v", ch)
	//wg.Done()
	if err != nil {
		log.Printf("SubscribeSupply Error...\n")
		return err
	} else {
		log.Print("Connect Synerex Server!\n")
	}
	for {
		var sp *Supply
		sp, err = smc.Recv() // receive Demand
		//log.Printf("\x1b[30m\x1b[47m SXUTIL: SUPPLY\x1b[0m\n")
		if err != nil {
			if err == io.EOF {
				log.Print("End Supply subscribe OK")
			} else {
				log.Printf("SMServiceClient SubscribeSupply error %v\n", err)
			}
			break
		}
		//		log.Println("Receive SS:", sp)
		// call Callback!
		spcb(clt, sp)
	}
	return err
}

// SubscribeDemand  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeDemand(ctx context.Context, dmcb func(*SMServiceClient, *Demand)) error {
	ch := clt.getChannel()
	dmc, err := clt.Client.SubscribeSyncDemand(ctx, ch)
	//log.Printf("Test3 %v", ch)
	//wg.Done()
	if err != nil {
		log.Printf("SubscribeDemand Error...\n")
		return err // sender should handle error...
	} else {
		log.Print("Connect Synerex Server!\n")
	}
	for {
		var dm *Demand
		dm, err = dmc.Recv() // receive Demand
		//log.Printf("\x1b[30m\x1b[47m SXUTIL: DEMAND\x1b[0m\n")
		if err != nil {
			if err == io.EOF {
				log.Print("End Demand subscribe OK")
			} else {
				log.Printf("SMServiceClient SubscribeDemand error %v\n", err)
			}
			break
		}
		//		log.Println("Receive SD:",*dm)
		// call Callback!
		dmcb(clt, dm)
	}
	return err
}

// SubscribeMbus  Wrapper function for SMServiceClient
func (clt *SMServiceClient) SubscribeMbus(ctx context.Context, mbcb func(*SMServiceClient, *MbusMsg)) error {

	mb := &Mbus{
		ClientId: uint64(clt.ClientID),
		MbusId:   uint64(clt.MbusID),
	}

	smc, err := clt.Client.SubscribeMbus(ctx, mb)
	if err != nil {
		log.Printf("%v Synerex_SubscribeMbusClient Error %v", clt, err)
		return err // sender should handle error...
	}
	for {
		var mes *MbusMsg
		mes, err = smc.Recv() // receive Demand
		if err != nil {
			if err == io.EOF {
				log.Print("End Mbus subscribe OK")
			} else {
				log.Printf("%v SMServiceClient SubscribeMbus error %v", clt, err)
			}
			break
		}
		//		log.Printf("Receive Mbus Message %v", *mes)
		// call Callback!
		mbcb(clt, mes)
	}
	return err
}

func (clt *SMServiceClient) SendMsg(ctx context.Context, msg *MbusMsg) error {
	if clt.MbusID == 0 {
		return errors.New("No Mbus opened!")
	}
	msg.MsgId = GenerateIntID()
	msg.SenderId = uint64(clt.ClientID)
	msg.MbusId = uint64(clt.MbusID)
	_, err := clt.Client.SendMsg(ctx, msg)

	return err
}

func (clt *SMServiceClient) CloseMbus(ctx context.Context) error {
	if clt.MbusID == 0 {
		return errors.New("No Mbus opened!")
	}
	mbus := &Mbus{
		ClientId: uint64(clt.ClientID),
		MbusId:   uint64(clt.MbusID),
	}
	_, err := clt.Client.CloseMbus(ctx, mbus)
	if err == nil {
		clt.MbusID = 0
	}
	return err
}

// RegisterDemand sends Typed Demand to Server
func (clt *SMServiceClient) RegisterDemand(dmo *DemandOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	dm := &Demand{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		DemandName: dmo.Name,
		Ts:         ts,
		ArgJson:    dmo.JSON,
	}

	if dmo.SimDemand != nil {
		dm.WithSimDemand(dmo.SimDemand)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := clt.Client.RegisterDemand(ctx, dm)

	//	resp, err := clt.Client.RegisterDemand(ctx, &dm)
	if err != nil {
		log.Printf("%v.RegisterDemand err %v", clt, err)
		return 0
	}
	//	log.Println(resp)
	dmo.ID = id // assign ID
	return id
}

// RegisterSupply sends Typed Supply to Server
func (clt *SMServiceClient) RegisterSupply(spo *SupplyOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	sp := &Supply{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		SupplyName: spo.Name,
		Ts:         ts,
		ArgJson:    spo.JSON,
	}

	if spo.SimSupply != nil {
		sp.WithSimSupply(spo.SimSupply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//	resp , err := clt.Client.RegisterSupply(ctx, &dm)

	_, err := clt.Client.RegisterSupply(ctx, sp)
	if err != nil {
		log.Printf("Error for sending:RegisterSupply to  Synerex Server as %v ", err)
		return 0
	}
	//	log.Println("RegiterSupply:", smo, resp)
	spo.ID = id // assign ID
	return id
}

//////////////////////////
// add sync function////////
/////////////////////////
// RegisterDemand sends Typed Demand to Server
func (clt *SMServiceClient) SyncDemand(dmo *DemandOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	dm := &Demand{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		DemandName: dmo.Name,
		Ts:         ts,
		ArgJson:    dmo.JSON,
	}

	if dmo.SimDemand != nil {
		dm.WithSimDemand(dmo.SimDemand)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := clt.Client.SyncDemand(ctx, dm)

	//	resp, err := clt.Client.SyncDemand(ctx, &dm)
	if err != nil {
		log.Printf("%v.SyncDemand err %v", clt, err)
		return 0
	}
	//	log.Println(resp)
	dmo.ID = id // assign ID
	return id
}

// SyncSupply sends Typed Supply to Server
func (clt *SMServiceClient) SyncSupply(spo *SupplyOpts) uint64 {
	id := GenerateIntID()
	ts := ptypes.TimestampNow()
	sp := &Supply{
		Id:         id,
		SenderId:   uint64(clt.ClientID),
		Type:       clt.MType,
		SupplyName: spo.Name,
		Ts:         ts,
		ArgJson:    spo.JSON,
	}

	if spo.SimSupply != nil {
		sp.WithSimSupply(spo.SimSupply)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//	resp , err := clt.Client.SyncSupply(ctx, &dm)

	_, err := clt.Client.SyncSupply(ctx, sp)
	if err != nil {
		log.Printf("Error for sending:SyncSupply to  Synerex Server as %v ", err)
		return 0
	}
	//	log.Println("RegiterSupply:", smo, resp)
	spo.ID = id // assign ID
	return id
}

// Confirm sends confirm message to sender
func (clt *SMServiceClient) Confirm(id IDType) error {
	tg := &Target{
		Id:       GenerateIntID(),
		SenderId: uint64(clt.ClientID),
		TargetId: uint64(id),
		Type:     clt.MType,
		MbusId:   uint64(id),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := clt.Client.Confirm(ctx, tg)
	if err != nil {
		log.Printf("%v Confirm Failier %v", clt, err)
		return err
	}
	clt.MbusID = id
	log.Println("Confirm Success:", resp)
	return nil
}

// Demand
// NewDemand returns empty Demand.
func NewDemand() *Demand {
	return &Demand{}
}

// NewSupply returns empty Supply.
func NewSupply() *Supply {
	return &Supply{}
}

func (dm *Demand) WithSimDemand(r *SimDemand) *Demand {
	dm.ArgOneof = &Demand_SimDemand{r}
	return dm
}

func (sp *Supply) WithSimSupply(c *SimSupply) *Supply {
	sp.ArgOneof = &Supply_SimSupply{c}
	return sp
}
