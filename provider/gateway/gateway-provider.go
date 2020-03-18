package main

// main synerex serverからgatewayを介してother synerex serverへ情報を送る
// 基本的に一方通行

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/synerex/synerex_alpha/api"
	simapi "github.com/synerex/synerex_alpha/api/simulation"
	provider "github.com/synerex/synerex_alpha/api/simulation/provider"
	"github.com/synerex/synerex_alpha/provider/simutil"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr           = flag.String("synerex", "127.0.0.1:10000", "The server address in the format of host:port")
	gatewayAddr          = flag.String("gateway", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv              = flag.String("nodeid", "127.0.0.1:9990", "Node ID Server")
	providerJson         = flag.String("provider_json", "", "Provider Json")
	scenarioProviderJson = flag.String("scenario_provider_json", "", "Provider Json")
	mu                   sync.Mutex
	myProvider           *provider.Provider
	scenarioProvider     *provider.Provider
	com1                 *simutil.Communicator
	com2                 *simutil.Communicator
	providerManager1     *simutil.ProviderManager
	providerManager2     *simutil.ProviderManager
	logger               *simutil.Logger
	mes1                 *Message
	mes2                 *Message
)

func flagToProviderInfo(pJson string) *provider.Provider {
	pInfo := &provider.Provider{}
	jsonpb.Unmarshal(strings.NewReader(pJson), pInfo)
	return pInfo
}

func init() {
	flag.Parse()
	logger = simutil.NewLogger()
	myProvider = flagToProviderInfo(*providerJson)
	scenarioProvider = flagToProviderInfo(*scenarioProviderJson)
	mes1 = NewMessage()
	mes2 = NewMessage()
}

// Supplyのコールバック関数
func supplyCallback1(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	//senderId := sp.GetSenderId()

	tid := sp.GetTargetId()
	pid := sp.GetSimSupply().GetPid()
	targets := sp.GetSimSupply().GetTargets()
	senderInfo := sp.GetSimSupply().GetSenderInfo()
	switch sp.GetSimSupply().GetType() {
	case simapi.SupplyType_GET_AGENTS_RESPONSE:
		if com2 != nil && IsContainProviders(providerManager1.Providers, senderInfo) && IsContainTarget(providerManager2.Providers, targets) {
			//logger.Error("GetAgentsResponse")
			agents := sp.GetSimSupply().GetGetAgentsResponse().GetAgents()
			agentType := sp.GetSimSupply().GetGetAgentsResponse().GetAgentType()
			areaId := sp.GetSimSupply().GetGetAgentsResponse().GetAreaId()

			//logger.Error("Send Agent from %v to %v", pid, tid)
			com2.GetAgentsResponse(senderInfo, targets, pid, tid, agents, agentType, areaId)
		}
	case simapi.SupplyType_GET_PROVIDERS_RESPONSE:
		providers := sp.GetSimSupply().GetGetProvidersResponse().GetProviders()
		//logger.Info("mes1")
		if mes1 != nil {
			mes1.Set(providers)
		}
	}

}

// Demandのコールバック関数
func demandCallback1(clt *sxutil.SMServiceClient, dm *pb.Demand) {

	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_GET_AGENTS_REQUEST:
		//logger.Error("GetAgentsRequest2")
	}

	// check if supply is match with my demand.
	pid := dm.GetSimDemand().GetPid()
	targets := dm.GetSimDemand().GetTargets()
	senderInfo := dm.GetSimDemand().GetSenderInfo()
	//senderId := dm.GetSenderId()
	if com2 != nil && IsContainProviders(providerManager1.Providers, senderInfo) && IsContainTarget(providerManager2.Providers, targets) {
		switch dm.GetSimDemand().GetType() {
		case simapi.DemandType_GET_AGENTS_REQUEST:
			//logger.Error("GetAgentsRequest")
			com2.GetAgentsRequest(senderInfo, targets, pid, nil)
		}
	}
}

// Supplyのコールバック関数
func supplyCallback2(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	//senderId := sp.GetSenderId()

	tid := sp.GetTargetId()
	pid := sp.GetSimSupply().GetPid()
	targets := sp.GetSimSupply().GetTargets()
	senderInfo := sp.GetSimSupply().GetSenderInfo()
	switch sp.GetSimSupply().GetType() {
	case simapi.SupplyType_GET_AGENTS_RESPONSE:
		if com1 != nil && IsContainProviders(providerManager2.Providers, senderInfo) && IsContainTarget(providerManager1.Providers, targets) {
			//logger.Error("GetAgentsResponse")
			agents := sp.GetSimSupply().GetGetAgentsResponse().GetAgents()
			agentType := sp.GetSimSupply().GetGetAgentsResponse().GetAgentType()
			areaId := sp.GetSimSupply().GetGetAgentsResponse().GetAreaId()

			//logger.Error("Send Agent from %v to %v", pid, tid)
			com1.GetAgentsResponse(senderInfo, targets, pid, tid, agents, agentType, areaId)
		}

	case simapi.SupplyType_GET_PROVIDERS_RESPONSE:
		providers := sp.GetSimSupply().GetGetProvidersResponse().GetProviders()
		if mes2 != nil {
			mes2.Set(providers)
		}
	}

}

// Demandのコールバック関数
func demandCallback2(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	//asenderId := dm.GetSenderId()
	pid := dm.GetSimDemand().GetPid()

	targets := dm.GetSimDemand().GetTargets()
	senderInfo := dm.GetSimDemand().GetSenderInfo()
	switch dm.GetSimDemand().GetType() {
	case simapi.DemandType_GET_AGENTS_REQUEST:
		//logger.Error("GetAgentsRequest2")
	}
	// 相手のプロバイダーへのDemandかつ自分のプロバイダーであること
	if com1 != nil && IsContainProviders(providerManager2.Providers, senderInfo) && IsContainTarget(providerManager1.Providers, targets) {
		switch dm.GetSimDemand().GetType() {
		case simapi.DemandType_GET_AGENTS_REQUEST:
			//logger.Error("GetAgentsRequest")
			com1.GetAgentsRequest(senderInfo, targets, pid, nil)
		}
	}
}

func IsContainProviders(providers []*provider.Provider, senderInfo *provider.Provider) bool {
	for _, pr := range providers {
		if pr.Id == senderInfo.Id {
			return true
		}
	}
	return false
}

func IsContainTarget(providers []*provider.Provider, targets []uint64) bool {
	for _, pr := range providers {
		for _, tgt := range targets {
			if pr.Id == tgt {
				return true
			}
		}
	}
	return false
}

func IsContainSameClientID(clients *simutil.Clients, id uint64) bool {
	if uint64(clients.AgentClient.ClientID) == id {
		return true
	}
	if uint64(clients.ClockClient.ClientID) == id {
		return true
	}
	if uint64(clients.ProviderClient.ClientID) == id {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////
////////////            Message Class           ///////////
///////////////////////////////////////////////////////////

type Message struct {
	ready     chan struct{}
	providers []*provider.Provider
}

func NewMessage() *Message {
	return &Message{ready: make(chan struct{})}
}
func (m *Message) Set(a []*provider.Provider) {
	m.providers = a
	logger.Info("Close")
	close(m.ready)
}

func (m *Message) Get() []*provider.Provider {
	<-m.ready
	return m.providers
}

func main() {
	logger.Info("StartUp Provider")

	// ProviderManager
	providerManager1 = simutil.NewProviderManager(myProvider)
	providerManager1.AddProvider(scenarioProvider)
	providerManager1.CreateIDMap()

	providerManager2 = simutil.NewProviderManager(myProvider)
	providerManager2.CreateIDMap()

	//////////////////////////////////////////////////
	//////////        node server        ////////////
	////////////////////////////////////////////////
	sxutil.RegisterNodeName(*nodesrv, "GatewayProvider", false)
	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	//////////////////////////////////////////////////
	//////////      main synerex server      ////////
	////////////////////////////////////////////////
	go func() {
		for {
			var opts []grpc.DialOption
			opts = append(opts, grpc.WithInsecure())
			conn, err := grpc.Dial(*serverAddr, opts...)
			if err != nil {
				log.Fatalf("fail to dial: %v", err)
				logger.Error("Fail to dial, Connect again...")
				time.Sleep(500 * time.Millisecond)
			} else {
				sxutil.RegisterDeferFunction(func() { conn.Close() })
				client := pb.NewSynerexClient(conn)
				argJson := fmt.Sprintf("{Client:Gateway}")

				// Communicator
				com1 = simutil.NewCommunicator()
				com1.RegistClients(client, argJson)                 // channelごとのClientを作成
				com1.SubscribeAll(demandCallback1, supplyCallback1) // ChannelにSubscribe
				logger.Info("Success to Connect! ServerAddr: %v", *serverAddr)

				return
			}
		}
	}()

	//////////////////////////////////////////////////
	//////////       other synerex server    ////////
	////////////////////////////////////////////////

	ch := make(chan *grpc.ClientConn)
	go func() {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithBlock())
		conn, err := grpc.Dial(*gatewayAddr, opts...)
		if err != nil {
			logger.Error("fail to dial: %v", err)
		} else {
			ch <- conn
			return

		}
	}()

	go func() {
		for {
			select {
			case conn := <-ch:
				sxutil.RegisterDeferFunction(func() { conn.Close() })
				client := pb.NewSynerexClient(conn)
				argJson := fmt.Sprintf("{Client:Gateway}")

				// Communicator
				com2 = simutil.NewCommunicator()
				com2.RegistClients(client, argJson) // channelごとのClientを作成
				// Subscribeは必要ない?
				com2.SubscribeAll(demandCallback2, supplyCallback2) // ChannelにSubscribe
				logger.Info("Success to Connect! GatewayAddr: %v", *gatewayAddr)

				// notify success connection to each synerex server
				if com1 != nil {
					pid := providerManager1.MyProvider.Id
					targets := []uint64{}
					senderInfo := providerManager2.MyProvider
					com1.GetProvidersRequest(senderInfo, targets, pid, nil)
					com2.GetProvidersRequest(senderInfo, targets, pid, nil)
					providers1 := mes1.Get()
					mes1 = nil
					providers2 := mes2.Get()
					mes2 = nil
					providerManager1.UpdateProviders(providers1)
					providerManager1.CreateIDMap()
					providerManager2.UpdateProviders(providers2)
					providerManager2.CreateIDMap()
					com1.SetProvidersRequest(senderInfo, targets, pid, nil, providers2)
					com2.SetProvidersRequest(senderInfo, targets, pid, nil, providers1)
				}

				return
			case <-time.After(3 * time.Second):
				logger.Error("fail to connect gateway. connect again...")
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
