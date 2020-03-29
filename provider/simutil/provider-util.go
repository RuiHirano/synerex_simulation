package simutil

/*import (
	//"log"
	"sync"

	"github.com/synerex/synerex_alpha/api/simulation/provider"
)

var (
	mu sync.Mutex
)

////////////////////////////////////////////////////////////
//////////////       Provider Manager Class      //////////
///////////////////////////////////////////////////////////

type IDType int

const (
	IDType_SCENARIO      IDType = 1
	IDType_CLOCK         IDType = 2
	IDType_VISUALIZATION IDType = 3
	IDType_AGENT         IDType = 4
	IDType_NEIGHBOR      IDType = 5
	IDType_SAME          IDType = 6
	IDType_GATEWAY       IDType = 7
)

type ProviderManager struct {
	MyProvider    *provider.Provider
	Providers     []*provider.Provider
	MyProviders   []*provider.Provider
	ProviderIDMap map[IDType][]uint64
}

func NewProviderManager(myProvider *provider.Provider) *ProviderManager {
	pm := &ProviderManager{
		MyProvider:  myProvider,
		Providers:   []*provider.Provider{myProvider},
		MyProviders: []*provider.Provider{myProvider},
	}
	return pm
}

func (pm *ProviderManager) AddProvider(p *provider.Provider) {
	mu.Lock()
	pm.Providers = append(pm.Providers, p)
	mu.Unlock()
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) AddMyProvider(p *provider.Provider) {
	mu.Lock()
	pm.MyProviders = append(pm.MyProviders, p)
	mu.Unlock()
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) UpdateProviders(ps []*provider.Provider) {
	mu.Lock()
	pm.Providers = ps
	mu.Unlock()
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) GetProviders() []*provider.Provider {
	mu.Lock()
	providers := pm.Providers
	mu.Unlock()
	return providers
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) SetProvider(index int, provider *provider.Provider) {
	mu.Lock()
	pm.Providers[index] = provider
	mu.Unlock()
}

func (pm *ProviderManager) DeleteProvider(id uint64) {
	newProviders := make([]*provider.Provider, 0)
	for _, provider := range pm.Providers {
		if provider.Id == id {
			continue
		}
		newProviders = append(newProviders, provider)
	}
	pm.Providers = newProviders
}

func (pm *ProviderManager) GetProviderNum() uint64 {
	return uint64(len(pm.Providers))
}

func (pm *ProviderManager) GetIDList(IdTypeList []IDType) []uint64 {
	idList := make([]uint64, 0)
	for _, idType := range IdTypeList {
		for _, id := range pm.ProviderIDMap[idType] {
			idList = append(idList, id)
		}
	}
	return idList
}

func (pm *ProviderManager) CreateIDMap() {
	providerIDMap := make(map[IDType][]uint64)
	sameIDs := make([]uint64, 0)
	neighborIDs := make([]uint64, 0)
	agentIDs := make([]uint64, 0)
	for _, p := range pm.Providers {
		switch p.GetType() {
		case provider.ProviderType_SCENARIO:
			providerIDMap[IDType_SCENARIO] = []uint64{p.GetId()}
		case provider.ProviderType_GATEWAY:
			providerIDMap[IDType_GATEWAY] = []uint64{p.GetId()}
		case provider.ProviderType_CLOCK:
			providerIDMap[IDType_CLOCK] = []uint64{p.GetId()}
		case provider.ProviderType_VISUALIZATION:
			providerIDMap[IDType_VISUALIZATION] = []uint64{p.GetId()}
		case provider.ProviderType_AGENT:
			if p.GetSynerexAddress() == pm.MyProvider.GetSynerexAddress() {
				agentIDs = append(agentIDs, p.GetId())
			}
			// AgentProviderでなければ必要ない
			if pm.MyProvider.GetType() == provider.ProviderType_AGENT {
				//log.Printf("IsNeighbor %v", pm.IsNeighborArea(p))
				if pm.IsNeighborArea(p) && p.GetAgentStatus().GetAgentType() == pm.MyProvider.GetAgentStatus().GetAgentType() {
					// 隣接エリアかつAgentTypeが等しい場合
					neighborIDs = append(neighborIDs, p.GetId())

				} else if pm.IsSameArea(p) && p.GetAgentStatus().GetAgentType() != pm.MyProvider.GetAgentStatus().GetAgentType() {
					// 同じエリアかつAgentTypeが等しくない場合
					sameIDs = append(sameIDs, p.GetId())
				}
			}

		}
	}
	providerIDMap[IDType_NEIGHBOR] = neighborIDs
	providerIDMap[IDType_SAME] = sameIDs
	providerIDMap[IDType_AGENT] = agentIDs
	pm.ProviderIDMap = providerIDMap

}

func (pm *ProviderManager) IsSameArea(p *provider.Provider) bool {
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
