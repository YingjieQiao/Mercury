package main

import (
	"sync"
	"time"
	"log"
)


type Manager struct {
	lock sync.RWMutex
	port uint64           //Peiyuan: local port number. In other words, which server does this node manager reside on.
	portToInfo         map[uint64]NodeInfo // Peiyuan: this is to keep track of the entire ring structure, including dead nodes.
	numOfvNodesPerpNodes    int
	consistent   *Consistent
	timers       map[uint64]*time.Timer
	liveDuration       time.Duration  // Peiyuan: how long will a node be marked alive without receiving heatbeat. 
	 								  // If node A does not receive heartbeat message from node B for liveDuration, node A will mark node B as dead.
}

func NewManager(port uint64, numOfvNodesPerpNodes int, liveDuration time.Duration) *Manager {
	return &Manager{
		port: port,
		portToInfo:   make(map[uint64]NodeInfo),
		numOfvNodesPerpNodes:    numOfvNodesPerpNodes,
		consistent:   NewConsistent(numOfvNodesPerpNodes),
		timers:       make(map[uint64]*time.Timer),
		liveDuration:       liveDuration,
	}
}

func (m *Manager) doUpdateNode(inComingNodeInfo NodeInfo) {

	currAlive := false
	if prev, ok := m.portToInfo[inComingNodeInfo.Port]; ok {
		currAlive = prev.Alive
	}
	m.portToInfo[inComingNodeInfo.Port] = inComingNodeInfo
	if !currAlive && inComingNodeInfo.Alive {
		log.Printf("Node %d has been added to consistent hash ring on Node %d", inComingNodeInfo.Port, m.port)
		m.consistent.AddPhysicalNode(inComingNodeInfo.Port)
	} else if currAlive && !inComingNodeInfo.Alive {
		log.Printf("Node %d has been removed consistent hash ring on Node %d", inComingNodeInfo.Port, m.port)
		m.consistent.Remove(inComingNodeInfo.Port)
	}
}

func (m *Manager) UpdateNode(node NodeInfo) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.doUpdateNode(node)

	nodeport := node.Port
	if timer, ok := m.timers[nodeport]; ok {
		timer.Reset(m.liveDuration)
		return
	}

	m.timers[nodeport] = time.AfterFunc(m.liveDuration, func() {
		m.lock.Lock()
		defer m.lock.Unlock()

		if !m.portToInfo[nodeport].Alive {
			return
		}
		nodeInfo := m.portToInfo[nodeport]
		nodeInfo.Alive = false
		m.portToInfo[nodeport] = nodeInfo
		m.consistent.Remove(nodeport)
		log.Printf("Node %v has been removed from consistent hash ring", node.Port)
	})
}

func (m *Manager) ExportRing() map[uint64]NodeInfo{
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.portToInfo
}

func (m *Manager) ImportRing(portToInfo map[uint64]NodeInfo) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for nodeport, info := range portToInfo {
		curr, ok := m.portToInfo[nodeport]

		if !ok || curr.Version < info.Version || (curr.Version == info.Version && !info.Alive) {
			m.doUpdateNode(info)
		}
	}
}

func (m *Manager) GetPreferenceList(key uint64, num int) ([]uint64, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.consistent.GetN(key, num)
}