package simutil

import (
	"log"
	"sync"

	"github.com/synerex/synerex_alpha/api"
)

var (
	mu sync.Mutex
)

////////////////////////////////////////////////////////////
//////////////       Provider Manager Class      //////////
///////////////////////////////////////////////////////////

type IDType int

const (
	IDType_MASTER        IDType = 1
	IDType_WORKER        IDType = 2
	IDType_VISUALIZATION IDType = 3
	IDType_AGENT         IDType = 4
	IDType_NEIGHBOR      IDType = 5
	IDType_SAME          IDType = 6
	IDType_GATEWAY       IDType = 7
)

type ProviderManager struct {
	MyProvider   *api.Provider
	Providers    []*api.Provider
	ProvidersMap map[IDType][]*api.Provider
}

func NewProviderManager(myProvider *api.Provider) *ProviderManager {
	pm := &ProviderManager{
		MyProvider:   myProvider,
		Providers:    []*api.Provider{},
		ProvidersMap: make(map[IDType][]*api.Provider),
	}
	return pm
}

func (pm *ProviderManager) AddProvider(p *api.Provider) {
	mu.Lock()
	pm.Providers = append(pm.Providers, p)
	pm.CreateProvidersMap()
	mu.Unlock()
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) SetProviders(ps []*api.Provider) {
	mu.Lock()
	pm.Providers = ps
	pm.CreateProvidersMap()
	mu.Unlock()
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) GetProviders() []*api.Provider {
	mu.Lock()
	providers := pm.Providers
	mu.Unlock()
	return providers
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) DeleteProvider(id uint64) {
	newProviders := make([]*api.Provider, 0)
	for _, provider := range pm.Providers {
		if provider.Id == id {
			continue
		}
		newProviders = append(newProviders, provider)
	}
	pm.Providers = newProviders
	pm.CreateProvidersMap()
}

func (pm *ProviderManager) GetProviderIds(IdTypeList []IDType) []uint64 {
	idList := make([]uint64, 0)
	for _, idType := range IdTypeList {
		for _, p := range pm.ProvidersMap[idType] {
			id := p.GetId()
			idList = append(idList, id)
		}
	}
	return idList
}

func (pm *ProviderManager) CreateProvidersMap() {
	providersMap := make(map[IDType][]*api.Provider)
	//sameProviders := make([]*api.Provider, 0)
	//neighborProviders := make([]*api.Provider, 0)
	//agentProviders := make([]*api.Provider, 0)
	log.Printf("providers: %v", pm.Providers)
	for _, p := range pm.Providers {
		if p.GetId() != pm.MyProvider.GetId() { // 自分は含まない
			switch p.GetType() {
			case api.ProviderType_MASTER:
				providersMap[IDType_MASTER] = append(providersMap[IDType_MASTER], p)
			case api.ProviderType_WORKER:
				providersMap[IDType_WORKER] = append(providersMap[IDType_WORKER], p)
			case api.ProviderType_GATEWAY:
				providersMap[IDType_GATEWAY] = append(providersMap[IDType_GATEWAY], p)
			case api.ProviderType_VISUALIZATION:
				providersMap[IDType_VISUALIZATION] = append(providersMap[IDType_VISUALIZATION], p)
			case api.ProviderType_AGENT:
				providersMap[IDType_AGENT] = append(providersMap[IDType_AGENT], p)
				// AgentProviderでなければ必要ない
				/*if pm.MyProvider.GetType() == api.ProviderType_AGENT {
					//log.Printf("IsNeighbor %v", pm.IsNeighborArea(p))
					if pm.IsNeighborArea(p) && p.GetAgentStatus().GetAgentType() == pm.MyProvider.GetAgentStatus().GetAgentType() {
						// 隣接エリアかつAgentTypeが等しい場合
						neighborProviders = append(neighborProviders, p)

					} else if pm.IsSameArea(p) && p.GetAgentStatus().GetAgentType() != pm.MyProvider.GetAgentStatus().GetAgentType() {
						// 同じエリアかつAgentTypeが等しくない場合
						sameProviders = append(sameProviders, p)
					}
				}*/

			}
		}
	}
	//providersMap[IDType_NEIGHBOR] = neighborProviders
	//providersMap[IDType_SAME] = sameProviders
	//providersMap[IDType_AGENT] = agentProviders
	pm.ProvidersMap = providersMap

}

/*func (pm *ProviderManager) IsSameArea(p *provider.Provider) bool {
	myAreaID := pm.MyProvider.GetAgentStatus().GetArea().GetId()
	opAreaID := p.GetAgentStatus().GetArea().GetId()
	if myAreaID == opAreaID {
		// エリアIDが等しければtrue
		return true
	}
	return false
}

// FIX
func (pm *ProviderManager) IsNeighborArea(p *provider.Provider) bool {
	myControlArea := pm.MyProvider.GetAgentStatus().GetArea().GetControlArea()
	tControlArea := p.GetAgentStatus().GetArea().GetControlArea()
	maxLat, maxLon, minLat, minLon := GetCoordRange(myControlArea)
	tMaxLat, tMaxLon, tMinLat, tMinLon := GetCoordRange(tControlArea)
	//log.Printf("latlon %v, %v, %v, %v", maxLat, maxLon, minLat, minLon)
	//log.Printf("latlon %v, %v, %v, %v", tMaxLat, tMaxLon, tMinLat, tMinLon)
	if maxLat == tMinLat && (minLon <= tMaxLon && tMaxLon <= maxLon || minLon <= tMinLon && tMinLon <= maxLon) {
		return true
	}
	if minLat == tMaxLat && (minLon <= tMaxLon && tMaxLon <= maxLon || minLon <= tMinLon && tMinLon <= maxLon) {
		return true
	}
	if maxLon == tMinLon && (minLat <= tMaxLat && tMaxLat <= maxLat || minLat <= tMinLat && tMinLat <= maxLat) {
		return true
	}
	if minLon == tMaxLon && (minLat <= tMaxLat && tMaxLat <= maxLat || minLat <= tMinLat && tMinLat <= maxLat) {
		return true
	}
	return false
}
*/
