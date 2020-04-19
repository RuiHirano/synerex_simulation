package simutil

import (
	//"log"
	"math"
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
	for _, pv := range pm.Providers {
		if pv.Id == p.Id {
			return
		}
	}
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
				if pm.MyProvider.GetType() == api.ProviderType_AGENT {
					myArea := pm.MyProvider.GetAgentStatus().GetArea()
					agentStatus := pm.MyProvider.GetAgentStatus()
					tgtArea := p.GetAgentStatus().GetArea()
					//log.Printf("IsNeighbor %v", pm.IsNeighborArea(p))
					if IsNeighborArea(myArea, tgtArea) && p.GetAgentStatus().GetAgentType() == agentStatus.GetAgentType() {
						// 隣接エリアかつAgentTypeが等しい場合
						//neighborProviders = append(neighborProviders, p)
						providersMap[IDType_NEIGHBOR] = append(providersMap[IDType_NEIGHBOR], p)

					} else if IsSameArea(myArea, tgtArea) && p.GetAgentStatus().GetAgentType() != agentStatus.GetAgentType() {
						// 同じエリアかつAgentTypeが等しくない場合
						//sameProviders = append(sameProviders, p)
						providersMap[IDType_SAME] = append(providersMap[IDType_SAME], p)
					}
				}

			}
		}
	}
	pm.ProvidersMap = providersMap

}

func IsSameArea(area1 *api.Area, area2 *api.Area) bool {
	if area1.GetId() == area2.GetId() {
		// エリアIDが等しければtrue
		return true
	}
	return false
}

// FIX
func IsNeighborArea(area1 *api.Area, area2 *api.Area) bool {
	myControlArea := area1.GetControlArea()
	tControlArea := area2.GetControlArea()
	maxLat, maxLon, minLat, minLon := GetCoordRange(myControlArea)
	tMaxLat, tMaxLon, tMinLat, tMinLon := GetCoordRange(tControlArea)

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

func GetCoordRange(coords []*api.Coord) (float64, float64, float64, float64) {
	maxLon, maxLat := math.Inf(-1), math.Inf(-1)
	minLon, minLat := math.Inf(0), math.Inf(0)
	for _, coord := range coords {
		if coord.Latitude > maxLat {
			maxLat = coord.Latitude
		}
		if coord.Longitude > maxLon {
			maxLon = coord.Longitude
		}
		if coord.Latitude < minLat {
			minLat = coord.Latitude
		}
		if coord.Longitude < minLon {
			minLon = coord.Longitude
		}
	}
	return maxLat, maxLon, minLat, minLon
}
