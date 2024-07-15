package podsmanagement

import (
	"time"
	"vortex_test/internal/model"
	"vortex_test/pkg/logging"
)

type Storage interface {
	GetCurrentAlgoritmStatus() ([]model.Algorithm, error)
}

type Deployer interface {
	CreatePod(name string) error
	DeletePod(name string) error
	GetPodList() ([]string, error)
}

type PodsManager struct {
	db         Storage
	logger     *logging.Logger
	activePods map[int]Deployer
}

func New(url Storage, logger *logging.Logger) *PodsManager {
	return &PodsManager{
		db:         url,
		logger:     logger,
		activePods: make(map[int]Deployer),
	}
}

func (pm *PodsManager) Start() {
	pm.logger.Info("Pods manager is active")
	pm.managerLoop()
}

func (pm *PodsManager) managerLoop() {
	for {
		time.Sleep(5 * time.Minute)
		pm.synchronizeAlgorithms()
	}
}

func (pm *PodsManager) synchronizeAlgorithms() {
	clients, err := pm.db.GetCurrentAlgoritmStatus()
	if err != nil {
		pm.logger.Error(err)
		return
	}
	for _, client := range clients {
		if _, ok := pm.activePods[client.ClientID]; !ok {
			pm.activePods[client.ClientID] = &PodList{}
		}

		if client.VWAP {
			pm.activePods[client.ClientID].CreatePod("algorithm vwap")
		} else {
			pm.activePods[client.ClientID].DeletePod("algorithm vwap")
		}

		if client.TWAP {
			pm.activePods[client.ClientID].CreatePod("algorithm twap")
		} else {
			pm.activePods[client.ClientID].DeletePod("algorithm twap")
		}

		if client.HFT {
			pm.activePods[client.ClientID].CreatePod("algorithm hft")
		} else {
			pm.activePods[client.ClientID].DeletePod("algorithm hft")
		}

		if !client.VWAP && !client.TWAP && !client.HFT {
			delete(pm.activePods, client.ClientID)
		}

		pm.logger.Infof("%d has pods: %v", client.ClientID, pm.activePods[client.ClientID])
	}
}
