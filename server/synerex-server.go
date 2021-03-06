package main

//go:generate protoc -I ../api -I .. --go_out=plugins=grpc:../api simulation.proto
//go:generate protoc -I ../api -I .. --go_out=plugins=grpc:../api synerex.proto
//go:generate protoc -I ../nodeapi -I .. --go_out=plugins=grpc:../nodeapi nodeid.proto

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"path"
	"strings"
	"sync"
	"time"

	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/synerex/synerex_alpha/api"
	napi "github.com/synerex/synerex_alpha/nodeapi"
	"github.com/synerex/synerex_alpha/util"
	"google.golang.org/grpc"
)

const MessageChannelBufferSize = 20

var (
	synerexAddr string
	nodeIdAddr  string
	serverName  string
	logger      *util.Logger
	nodeapi     *napi.NodeAPI
)

type synerexServerInfo struct {
	demandChans        [api.ChannelType_END][]chan *api.Demand // create slices for each ChannelType(each slice contains channels)
	supplyChans        [api.ChannelType_END][]chan *api.Supply
	mbusChans          map[uint64][]chan *api.MbusMsg                       // Private Message bus for each provider
	mbusMap            map[api.IDType]map[uint64]chan *api.MbusMsg          // map from IDtype to Mbus channel
	demandMap          [api.ChannelType_END]map[api.IDType]chan *api.Demand // map from IDtype to Demand channel
	supplyMap          [api.ChannelType_END]map[api.IDType]chan *api.Supply // map from IDtype to Supply channel
	waitConfirms       [api.ChannelType_END]map[api.IDType]chan *api.Target // confirm maps
	dmu, smu, mmu, wmu sync.RWMutex
	messageStore       *MessageStore // message store
}

func init() {
	synerexAddr = os.Getenv("SYNEREX_SERVER")
	if synerexAddr == "" {
		synerexAddr = "127.0.0.1:10000"
	}
	nodeIdAddr = os.Getenv("NODEID_SERVER")
	if nodeIdAddr == "" {
		nodeIdAddr = "127.0.0.1:9000"
	}

	serverName = os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "SynerexServer"
	}

	logger = util.NewLogger()

}

// Implementation of each Protocol API
func (s *synerexServerInfo) RegisterDemand(c context.Context, dm *api.Demand) (r *api.Response, e error) {
	// send demand for desired channels
	okFlag := true
	okMsg := ""
	s.dmu.RLock()
	chs := s.demandChans[dm.GetType()]
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize {
			ch <- dm
		} else {
			okMsg = fmt.Sprintf("RD MessageDrop %v", dm)
			okFlag = false
			log.Printf("RD MessageDrop %v\n", dm)
		}
	}
	s.dmu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}

func (s *synerexServerInfo) RegisterSupply(c context.Context, sp *api.Supply) (r *api.Response, e error) {
	//fmt.Printf("Register Supply!!!")
	okFlag := true
	okMsg := ""
	str := ""
	s.smu.RLock()
	chs := s.supplyChans[sp.GetType()]
	for i := range chs {
		ch := chs[i]
		str = str + fmt.Sprintf("%d ", len(ch))
		if len(ch) < MessageChannelBufferSize { // run under not blocking state.
			ch <- sp
		} else {
			okMsg = fmt.Sprintf("RS MessageDrop %v", sp)
			okFlag = false
			log.Printf("RS MessageDrop %v\n", sp)
		}

	}
	s.smu.RUnlock()
	//fmt.Printf("RS: %d, %s:", len(chs), str)
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}
func (s *synerexServerInfo) ProposeDemand(c context.Context, dm *api.Demand) (r *api.Response, e error) {
	okFlag := true
	okMsg := ""
	s.dmu.RLock()
	chs := s.demandChans[dm.GetType()]
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize {
			ch <- dm
		} else {
			okMsg = fmt.Sprintf("PD MessageDrop %v", dm)
			okFlag = false
			log.Printf("PD MessageDrop %v\n", dm)
		}
	}
	s.dmu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}
func (s *synerexServerInfo) ProposeSupply(c context.Context, sp *api.Supply) (r *api.Response, e error) {
	okFlag := true
	okMsg := ""
	s.smu.RLock()
	chs := s.supplyChans[sp.GetType()]
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize {
			ch <- sp
		} else {
			okMsg = fmt.Sprintf("PS MessageDrop %v", sp)
			okFlag = false
			log.Printf("PS MessageDrop %v\n", sp)
		}
	}
	s.smu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}

////////////////////////////////////////////////////////////
//////////////////        Sync           ///////////////////
///////////////////////////////////////////////////////////
func (s *synerexServerInfo) SyncDemand(c context.Context, dm *api.Demand) (r *api.Response, e error) {
	okFlag := true
	okMsg := ""
	s.dmu.RLock()
	chs := s.demandChans[dm.GetType()]
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize {

			ch <- dm
		} else {
			okMsg = fmt.Sprintf("PD MessageDrop %v", dm)
			okFlag = false
			log.Printf("PD MessageDrop %v\n", dm)
		}
	}
	s.dmu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}
func (s *synerexServerInfo) SyncSupply(c context.Context, sp *api.Supply) (r *api.Response, e error) {
	okFlag := true
	okMsg := ""
	s.smu.RLock()
	chs := s.supplyChans[sp.GetType()]
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize {
			ch <- sp
		} else {
			okMsg = fmt.Sprintf("PS MessageDrop %v", sp)
			okFlag = false
			log.Printf("PS MessageDrop %v\n", sp)
		}
	}
	s.smu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}

// go routine which wait demand channel and sending demands to each providers.
func syncDemandServerFunc(ch chan *api.Demand, stream api.Synerex_SubscribeSyncDemandServer, id uint64) error {
	for {
		select {
		case dm := <-ch: // may block until receiving info
			//targets = dm.GetSimDemand().GetTargets()
			//dmId = dm.GetId()
			//s.SyncMap[dmId] = targets

			targets := dm.GetSimDemand().GetTargets()   // id is channelID, targets is target channelIDs
			senderId := dm.GetSimDemand().GetSenderId() // not to recieve own message
			//fmt.Println("get dm %v targets: %v\n", targets)
			if contains(id, targets) || (len(targets) == 0 && id != senderId) {
				//fmt.Println("dm %v\n", dm)
				err := stream.Send(dm)

				if err != nil {
					log.Printf("Error in DemandServer Error %v", err, dm, id)
					return err
				}
			}
		}
	}
}

// SubscribeDemand is called form client to subscribe channel
func (s *synerexServerInfo) SubscribeSyncDemand(ch *api.Channel, stream api.Synerex_SubscribeSyncDemandServer) error {
	// TODO: we can check the duplication of node id here! (especially 1024 snowflake node ID)
	idt := api.IDType(ch.GetClientId())
	s.dmu.RLock()
	_, ok := s.demandMap[ch.Type][idt]
	s.dmu.RUnlock()
	if ok { // check the availability of duplicated client ID
		return errors.New(fmt.Sprintf("duplicated SubscribeDemand ClientID %d", idt))
	}

	subCh := make(chan *api.Demand, MessageChannelBufferSize)
	// We should think about thread safe coding.
	tp := ch.GetType()
	s.dmu.Lock()
	s.demandChans[tp] = append(s.demandChans[tp], subCh)
	s.demandMap[tp][idt] = subCh // mapping from clientID to channel
	s.dmu.Unlock()
	pid := ch.ProviderId
	syncDemandServerFunc(subCh, stream, pid) // infinite go routine?
	// if this returns, stream might be closed.
	// we should remove channel
	s.dmu.Lock()
	delete(s.demandMap[tp], idt) // remove map from idt
	s.demandChans[tp] = removeDemandChannelFromSlice(s.demandChans[tp], subCh)
	log.Printf("Remove Demand Stream Channel %v", ch)
	s.dmu.Unlock()
	return nil
}

func syncSupplyServerFunc(ch chan *api.Supply, stream api.Synerex_SubscribeSyncSupplyServer, id uint64) error {
	for {
		select {
		case sp := <-ch:
			targets := sp.GetSimSupply().GetTargets()   // id is channelID, targets is target channelIDs
			senderId := sp.GetSimSupply().GetSenderId() // not to recieve own message
			//fmt.Println("get sp %v targets: %v\n", sp, targets)
			if contains(id, targets) || (len(targets) == 0 && id != senderId) { //len(targets) == 0: broadcast, else unicast
				//fmt.Println("sp %v\n", sp)
				err := stream.Send(sp) // send to client
				if err != nil {
					log.Printf("Error SupplyServer Error %v", err, sp, id)
					return err
				}
			}
		}
	}
}

func (s *synerexServerInfo) SubscribeSyncSupply(ch *api.Channel, stream api.Synerex_SubscribeSyncSupplyServer) error {
	idt := api.IDType(ch.GetClientId())
	tp := ch.GetType()
	s.smu.RLock()
	_, ok := s.supplyMap[tp][idt]
	s.smu.RUnlock()
	if ok { // check the availability of duplicated client ID
		return errors.New(fmt.Sprintf("duplicated SubscribeSupply for ClientID %v", idt))
	}

	subCh := make(chan *api.Supply, MessageChannelBufferSize)

	s.smu.Lock()
	s.supplyChans[tp] = append(s.supplyChans[tp], subCh)
	s.supplyMap[tp][idt] = subCh // mapping from clientID to channel
	s.smu.Unlock()
	pid := ch.ProviderId // providerIDで送信先を決定する
	//fmt.Printf("ProviderID: %v\n", pid)
	err := syncSupplyServerFunc(subCh, stream, pid)
	// this supply stream may closed. so take care.

	s.smu.Lock()
	delete(s.supplyMap[tp], idt) // remove map from idt
	s.supplyChans[tp] = removeSupplyChannelFromSlice(s.supplyChans[tp], subCh)
	log.Printf("Remove Supply Stream Channel %v", ch)
	s.smu.Unlock()

	return err
}

func contains(id uint64, targets []uint64) bool {
	for _, tgt := range targets {
		if tgt == id {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////
//////////////////        End          ///////////////////
///////////////////////////////////////////////////////////

func (s *synerexServerInfo) ReserveSupply(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	//	chs := s.demandChans[tg.GetType()]
	//	dm := &api.Demand{}

	r = &api.ConfirmResponse{Ok: true, Err: ""}
	return r, nil
}

func (s *synerexServerInfo) SelectSupply(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	targetSender := s.messageStore.getSrcId(tg.GetTargetId()) // find source from Id
	s.dmu.RLock()
	ch, ok := s.demandMap[tg.Type][api.IDType(targetSender)]
	s.dmu.RUnlock()
	if !ok {
		r = &api.ConfirmResponse{Ok: false, Err: "Can't find demand target from SelectSupply"}
		log.Printf("Can't find SelectSupply target ID %d, src %d", tg.GetTargetId(), targetSender)
		e = errors.New("Cant find channel in SelectSupply")
		return
	}
	id := api.GenerateIntID()
	dm := &api.Demand{
		Id:       id, // generate ID from market server
		SenderId: tg.SenderId,
		TargetId: tg.TargetId,
		Type:     tg.Type,
		MbusId:   id, // mbus id is a message id for select.
	}

	tch := make(chan *api.Target)
	s.wmu.Lock()
	s.waitConfirms[tg.Type][api.IDType(id)] = tch
	s.wmu.Unlock()

	ch <- dm // send select message

	// wait for confim...
	select {

	case tb := <-tch: // got confirm!
		s.wmu.Lock() // remove waitChannel
		delete(s.waitConfirms[tg.Type], api.IDType(id))
		s.wmu.Unlock()

		if tb.TargetId == id {
			if tb.MbusId == id {
				r = &api.ConfirmResponse{Ok: true, Err: "", MbusId: id}
				return r, nil
			} else {
				r = &api.ConfirmResponse{Ok: true, Err: "no mbus id"}
				return r, nil
			}
		}

	case <-time.After(30 * time.Second): // timeout!
		//args := idToNode(tg.SenderId) + "->" + idToNode(tg.TargetId)
		//go monitorapi.SendMessage("notConfirm", int(tg.Type), dm.Id, tg.SenderId, tg.TargetId, tg.TargetId, args)
		r = &api.ConfirmResponse{Ok: false, Err: "waitConfirm Timeout!"}

	}

	return r, errors.New("Should not happen")

}

func (s *synerexServerInfo) ReserveDemand(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	r = &api.ConfirmResponse{Ok: true, Err: ""}
	return r, nil
}
func (s *synerexServerInfo) SelectDemand(c context.Context, tg *api.Target) (r *api.ConfirmResponse, e error) {
	// select!

	r = &api.ConfirmResponse{Ok: true, Err: ""}
	return r, nil
}

func (s *synerexServerInfo) Confirm(c context.Context, tg *api.Target) (r *api.Response, e error) {
	// check waitConfirms
	s.wmu.RLock()
	ch, ok := s.waitConfirms[tg.Type][api.IDType(tg.TargetId)]
	s.wmu.RUnlock()

	if !ok {
		r = &api.Response{Ok: false, Err: "Can't find channel"}
		return r, errors.New("can't find channels for Confirm")
	}
	ch <- tg // send OK
	r = &api.Response{Ok: true, Err: ""}
	return r, nil
}

// go routine which wait demand channel and sending demands to each providers.
func demandServerFunc(ch chan *api.Demand, stream api.Synerex_SubscribeDemandServer, id api.IDType) error {
	for {
		select {
		case dm := <-ch: // may block until receiving info
			fmt.Println("dm %v\n", dm)
			err := stream.Send(dm)
			if err != nil {
				log.Printf("Error in DemandServer Error %v", err, dm)
				return err
			}
		}
	}
}

// remove channel from slice

func removeDemandChannelFromSlice(sl []chan *api.Demand, c chan *api.Demand) []chan *api.Demand {
	for i, ch := range sl {
		if ch == c {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	log.Printf("Cant find channel %v in removeChannel", c)
	return sl
}

func removeSupplyChannelFromSlice(sl []chan *api.Supply, c chan *api.Supply) []chan *api.Supply {
	for i, ch := range sl {
		if ch == c {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	log.Printf("Cant find channel %v in removeChannel", c)
	return sl
}

// SubscribeDemand is called form client to subscribe channel
func (s *synerexServerInfo) SubscribeDemand(ch *api.Channel, stream api.Synerex_SubscribeDemandServer) error {
	// TODO: we can check the duplication of node id here! (especially 1024 snowflake node ID)
	idt := api.IDType(ch.GetClientId())
	s.dmu.RLock()
	_, ok := s.demandMap[ch.Type][idt]
	s.dmu.RUnlock()
	if ok { // check the availability of duplicated client ID
		return errors.New(fmt.Sprintf("duplicated SubscribeDemand ClientID %d", idt))
	}

	subCh := make(chan *api.Demand, MessageChannelBufferSize)
	// We should think about thread safe coding.
	tp := ch.GetType()
	s.dmu.Lock()
	s.demandChans[tp] = append(s.demandChans[tp], subCh)
	s.demandMap[tp][idt] = subCh // mapping from clientID to channel
	s.dmu.Unlock()
	demandServerFunc(subCh, stream, idt) // infinite go routine?
	// if this returns, stream might be closed.
	// we should remove channel
	s.dmu.Lock()
	delete(s.demandMap[tp], idt) // remove map from idt
	s.demandChans[tp] = removeDemandChannelFromSlice(s.demandChans[tp], subCh)
	log.Printf("Remove Demand Stream Channel %v", ch)
	s.dmu.Unlock()
	return nil
}

// This function is created for each subscribed provider
// This is not efficient if the number of providers increases.
func supplyServerFunc(ch chan *api.Supply, stream api.Synerex_SubscribeSupplyServer) error {
	for {
		select {
		case sp := <-ch:
			err := stream.Send(sp)
			if err != nil {

				log.Printf("Error SupplyServer Error %v", err, sp)
				return err
			}
		}
	}
}

func (s *synerexServerInfo) SubscribeSupply(ch *api.Channel, stream api.Synerex_SubscribeSupplyServer) error {
	idt := api.IDType(ch.GetClientId())
	tp := ch.GetType()
	s.smu.RLock()
	_, ok := s.supplyMap[tp][idt]
	s.smu.RUnlock()
	if ok { // check the availability of duplicated client ID
		return errors.New(fmt.Sprintf("duplicated SubscribeSupply for ClientID %v", idt))
	}

	subCh := make(chan *api.Supply, MessageChannelBufferSize)

	s.smu.Lock()
	s.supplyChans[tp] = append(s.supplyChans[tp], subCh)
	s.supplyMap[tp][idt] = subCh // mapping from clientID to channel
	s.smu.Unlock()
	err := supplyServerFunc(subCh, stream)
	// this supply stream may closed. so take care.

	s.smu.Lock()
	delete(s.supplyMap[tp], idt) // remove map from idt
	s.supplyChans[tp] = removeSupplyChannelFromSlice(s.supplyChans[tp], subCh)
	log.Printf("Remove Supply Stream Channel %v", ch)
	s.smu.Unlock()

	return err
}

// This function is created for each subscribed provider
// This is not efficient if the number of providers increases.
func mbusServerFunc(ch chan *api.MbusMsg, stream api.Synerex_SubscribeMbusServer, id api.IDType) error {
	for {
		select {
		case msg := <-ch:
			if msg.GetMsgId() == 0 { // close message
				return nil // grace close
			}
			if api.IDType(msg.GetSenderId()) != id { // do not send msg from myself
				tgt := api.IDType(msg.GetTargetId())
				if tgt == 0 || tgt == id { // =0 broadcast , = tgt unicast
					err := stream.Send(msg)
					if err != nil {
						//				log.Printf("Error mBus Error %v", err)
						return err
					}
				}
			}
		}
	}
}

func removeMbusChannelFromSlice(sl []chan *api.MbusMsg, c chan *api.MbusMsg) []chan *api.MbusMsg {
	for i, ch := range sl {
		if ch == c {
			return append(sl[:i], sl[i+1:]...)
		}
	}
	log.Printf("Cant find channel %v in removeMbusChannel", c)
	return sl
}
func (s *synerexServerInfo) SubscribeMbus(mb *api.Mbus, stream api.Synerex_SubscribeMbusServer) error {

	mbusCh := make(chan *api.MbusMsg, MessageChannelBufferSize) // make channel for each mbus
	id := api.IDType(mb.GetClientId())
	mbid := mb.MbusId
	s.mmu.Lock()
	chans := s.mbusChans[mbid]
	s.mbusChans[mbid] = append(chans, mbusCh)
	mm, ok := s.mbusMap[id]
	if ok {
		//		mm[mbid] = mbusCh
	} else {
		mm = make(map[uint64]chan *api.MbusMsg)
		mm[mbid] = mbusCh
		s.mbusMap[id] = mm
	}
	s.mmu.Unlock()

	err := mbusServerFunc(mbusCh, stream, id)

	s.mmu.Lock()
	s.mbusChans[mbid] = removeMbusChannelFromSlice(s.mbusChans[mbid], mbusCh)
	delete(s.mbusMap, id)
	//	log.Printf("Remove Mbus Stream Channel %v", ch)
	s.mmu.Unlock()

	return err
}

func (s *synerexServerInfo) SendMsg(c context.Context, msg *api.MbusMsg) (r *api.Response, err error) {
	// FIXME: wait until all subscriber is comming
	for {
		chans, ok := s.mbusChans[msg.GetMbusId()]
		if ok && len(chans) == 2 {
			log.Printf("##### All subscriber comming!! [MbusID: %d]\n", msg.GetMbusId())
			break
		}
		log.Printf("##### Another Subscriber wating... [MbusId: %d, len(chans): %d]\n", msg.GetMbusId(), len(chans))
		time.Sleep(1 * time.Second)
	}
	okFlag := true
	okMsg := ""
	s.mmu.RLock()
	chs := s.mbusChans[msg.GetMbusId()] // get channel slice from mbus_id
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize { // run under not blocking state.
			ch <- msg
		} else {
			okMsg = fmt.Sprintf("MBus MessageDrop %v", msg)
			okFlag = false
			log.Printf("Mbus MessageDrop %v\n", msg)
		}
	}
	s.mmu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}

func (s *synerexServerInfo) CloseMbus(c context.Context, mb *api.Mbus) (r *api.Response, err error) {
	okFlag := true
	okMsg := ""
	s.mmu.RLock()
	chs := s.mbusChans[mb.GetMbusId()] // get channel slice from mbus_id
	cmsg := &api.MbusMsg{              // this is close message
		MsgId: 0,
	}
	for i := range chs {
		ch := chs[i]
		if len(ch) < MessageChannelBufferSize { // run under not blocking state.
			ch <- cmsg
		} else {
			okMsg = fmt.Sprintf("MBusClose MessageDrop %v", cmsg)
			okFlag = false
			log.Printf("MBusClose MessageDrop %v\n", cmsg)
		}
	}
	s.mmu.RUnlock()
	r = &api.Response{Ok: okFlag, Err: okMsg}
	return r, nil
}

func newServerInfo() *synerexServerInfo {
	var ms synerexServerInfo
	s := &ms
	for i := 0; i < int(api.ChannelType_END); i++ {
		s.demandMap[i] = make(map[api.IDType]chan *api.Demand)
		s.supplyMap[i] = make(map[api.IDType]chan *api.Supply)
		s.waitConfirms[i] = make(map[api.IDType]chan *api.Target)
	}
	s.mbusChans = make(map[uint64][]chan *api.MbusMsg)
	s.mbusMap = make(map[api.IDType]map[uint64]chan *api.MbusMsg)
	s.messageStore = CreateLocalMessageStore()

	return s
}

// synerex ID system
var (
	NodeBits uint8 = 10
	StepBits uint8 = 12

	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask  int64 = nodeMax << StepBits
	nodeShift uint8 = StepBits
	nodeMap         = make(map[int]string)
)

func idToNode(id uint64) string {
	nodeNum := int(int64(id) & nodeMask >> nodeShift) // snowflake node ID:
	var ok bool
	var str string
	if str, ok = nodeMap[nodeNum]; !ok {
		str = nodeapi.GetNodeName(nodeNum)
	}
	rs := strings.Replace(str, "Provider", "", -1)
	rs2 := strings.Replace(rs, "Server", "", -1)
	return rs2 + ":" + strconv.Itoa(nodeNum)
}

func unaryServerInterceptor(logger *logrus.Logger, s *synerexServerInfo) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var err error
		var args string
		var msgType int
		var srcId, tgtId, mid uint64
		method := path.Base(info.FullMethod)
		switch method {
		// Demand
		case "RegisterDemand", "ProposeDemand":
			dm := req.(*api.Demand)
			msgType = int(dm.Type)
			srcId = dm.SenderId
			tgtId = dm.TargetId
			mid = dm.Id
			//			args = "Type:" + strconv.Itoa(int(dm.Type)) + ":" + strconv.FormatUint(dm.Id, 16) + ":" + idToNode(dm.SenderId) + "->" + strconv.FormatUint(dm.TargetId, 16)
			//args = idToNode(dm.SenderId) + "->" + idToNode(dm.TargetId)
			// 変更
			cdm := dm // copy
			//cdm.GetSimDemand().Data = &simapi.SimDemand_GetAgentsRequest{&agent.GetAgentsRequest{}}
			//logger.Info("dm: %v\n", cdm.GetSimDemand())
			dm_json, _ := json.Marshal(cdm)
			args = string(dm_json)
			// Supply
		case "RegisterSupply", "ProposeSupply":
			sp := req.(*api.Supply)
			msgType = int(sp.Type)
			srcId = sp.SenderId
			tgtId = sp.TargetId
			mid = sp.Id
			//			args = "Type:" + strconv.Itoa(int(sp.Type)) + ":" + strconv.FormatUint(sp.Id, 16) + ":" + idToNode(sp.SenderId) + "->" + strconv.FormatUint(sp.TargetId, 16)
			//args = idToNode(sp.SenderId) + "->" + idToNode(sp.TargetId)
			csp := sp // copy
			//csp.GetSimSupply().Data = &simapi.SimSupply_SetAgentsResponse{&agent.SetAgentsResponse{}}
			//logger.Info("sp: %v\n", csp.GetSimSupply())
			sp_json, _ := json.Marshal(csp)
			args = string(sp_json)
		// Target
		case "SelectSupply", "Confirm", "SelectDemand":
			tg := req.(*api.Target)
			msgType = int(tg.Type)
			mid = tg.Id
			srcId = tg.SenderId
			tgtId = tg.TargetId
			args = idToNode(tg.SenderId) + "->" + idToNode(tg.TargetId)
			//			args = "Type:" + strconv.Itoa(int(tg.Type)) + ":" + strconv.FormatUint(tg.Id, 16) + ":" + idToNode(tg.Id) + "->" + strconv.FormatUint(tg.TargetId, 16)
		case "SendMsg":
			msg := req.(*api.MbusMsg)
			msgType = int(msg.MsgType)
			mid = msg.MsgId
			srcId = msg.SenderId
			tgtId = msg.TargetId
			args = idToNode(msg.SenderId) + "->" + idToNode(msg.TargetId)

		}

		dstId := s.messageStore.getSrcId(tgtId)

		// register for messageStore
		s.messageStore.AddMessage(method, msgType, mid, srcId, dstId, args)

		// Obtain log using defer
		defer func(begin time.Time) {
			// Obtain method name from info
			method := path.Base(info.FullMethod)
			took := time.Since(begin)
			fields := logrus.Fields{
				"method": method,
				"took":   took,
			}
			if err != nil {
				fields["error"] = err
				logger.WithFields(fields).Error("Failed")
			} else {
				//				logger.WithFields(fields).Info("Succeeded")
			}
		}(time.Now())

		// handler = RPC method
		reply, hErr := handler(ctx, req)
		if hErr != nil {
			err = hErr
		}

		return reply, err
	}
}

// Stream Interceptor
func streamServerInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var err error
		//		var args string
		log.Printf("streamserver intercept...")
		method := path.Base(info.FullMethod)
		switch method {
		case "SubscribeDemand":
		case "SubscribeSupply":
		}

		defer func(begin time.Time) {
			// Obtain method name from info
			method := path.Base(info.FullMethod)
			took := time.Since(begin)
			fields := logrus.Fields{
				"method": method,
				"took":   took,
			}
			if err != nil {
				fields["error"] = err
				logger.WithFields(fields).Error("Failed")
			} else {
				logger.WithFields(fields).Info("Succeeded")
			}
		}(time.Now())

		// handler = RPC method
		if hErr := handler(srv, stream); err != nil {
			err = hErr
		}
		log.Printf("streamserver intercept..end .")
		return err
	}
}

func prepareGrpcServer(s *synerexServerInfo, opts ...grpc.ServerOption) *grpc.Server {
	gcServer := grpc.NewServer(opts...)
	api.RegisterSynerexServer(gcServer, s)
	return gcServer
}

func main() {
	flag.Parse()
	nodeapi = napi.NewNodeAPI()
	for {
		err := nodeapi.RegisterNodeName(nodeIdAddr, serverName, true)
		if err == nil {
			logger.Info("connected NodeID server!")
			go nodeapi.HandleSigInt()
			nodeapi.RegisterDeferFunction(nodeapi.UnRegisterNode)
			break
		} else {
			logger.Warn("NodeID Error... reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}

	lis, err := net.Listen("tcp", synerexAddr)
	defer lis.Close()
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	logger := logrus.New()

	s := newServerInfo()
	opts = append(opts, grpc.UnaryInterceptor(unaryServerInterceptor(logger, s)))

	grpcServer := prepareGrpcServer(s, opts...)
	log.Printf("Start Synergic Exchange Server, connection waiting at %d ...", synerexAddr)
	serr := grpcServer.Serve(lis)
	log.Printf("Should not arrive here.. server closed. %v", serr)

}
