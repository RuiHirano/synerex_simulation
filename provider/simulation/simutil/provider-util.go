package simutil

import (
	"github.com/synerex/synerex_alpha/api/simulation/provider"
)

////////////////////////////////////////////////////////////
//////////////       Provider Manager Class      //////////
///////////////////////////////////////////////////////////

type ProviderManager struct {
	Providers []*provider.Provider
}

func NewProviderManager() *ProviderManager {
	pm := &ProviderManager{
		Providers: make([]*provider.Provider, 0),
	}
	return pm
}

func (pm *ProviderManager) AddProvider(provider *provider.Provider) {
	pm.Providers = append(pm.Providers, provider)
	//log.Printf("Providers: %v\n", pm.Providers)
}

func (pm *ProviderManager) SetProvider(index int, provider *provider.Provider) {
	pm.Providers[index] = provider
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

func (pm *ProviderManager) GetProviderIDs(providerTypes []provider.ProviderType) []uint64 {
	providerIDs := make([]uint64, 0)
	for _, pType := range providerTypes {
		for _, provider := range pm.Providers {
			if provider.Type == pType {
				providerIDs = append(providerIDs, provider.Id)
			}
		}
	}
	return providerIDs
}

func (pm *ProviderManager) GetNeighborIDs(providerTypes []provider.ProviderType) []uint64 {
	providerIDs := make([]uint64, 0)
	for _, pType := range providerTypes {
		for _, provider := range pm.Providers {
			if provider.Type == pType {
				providerIDs = append(providerIDs, provider.Id)
			}
		}
	}
	return providerIDs
}

func (pm *ProviderManager) GetSameIDs(providerTypes []provider.ProviderType) []uint64 {
	providerIDs := make([]uint64, 0)
	for _, pType := range providerTypes {
		for _, provider := range pm.Providers {
			if provider.Type == pType {
				providerIDs = append(providerIDs, provider.Id)
			}
		}
	}
	return providerIDs
}
