package main

import (
	"encoding/binary"
	"errors"
	"hash/fnv"
	"sort"
	"sync"
)

type uints []uint64


func (x uints) Len() int { return len(x) }


func (x uints) Less(i, j int) bool { return x[i] < x[j] }


func (x uints) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

var ErrEmptyCircle = errors.New("empty vNodeID2pNodeID")


type Consistent struct {
	vNodeID2pNodeID           map[uint64]uint64 // hash value to node name
	pNodeID2vNodeID          map[uint64] [] uint64//set of nodes
	sortedHashes     uints //all hash values on the ring in order
	numOfvNodesPerpNodes int // number of vitual node per physical node
	sync.RWMutex
}

func NewConsistent(numOfvNodesPerpNodes int) *Consistent {
	c := new(Consistent)
	c.numOfvNodesPerpNodes = numOfvNodesPerpNodes
	c.vNodeID2pNodeID = make(map[uint64]uint64)
	c.pNodeID2vNodeID = make(map[uint64][]uint64)
	return c
}


func (c *Consistent) AddPhysicalNode(physicalNodeId uint64) {
	c.Lock()
	defer c.Unlock()
	vNodeID := physicalNodeId
	vNodes := make([]uint64, c.numOfvNodesPerpNodes)
	for i := 0; i < c.numOfvNodesPerpNodes; i++ {
		vNodeID = c.hash(vNodeID)
		c.vNodeID2pNodeID[vNodeID] = physicalNodeId
		vNodes[i] = vNodeID
	}
	c.pNodeID2vNodeID[physicalNodeId] = vNodes
	c.updateSortedHashes()
}





func (c *Consistent) Remove(pNodeID uint64) {
	c.Lock()
	defer c.Unlock()
	vNodes, ok := c.pNodeID2vNodeID[pNodeID]
	if !ok{
		return
	}
	for i := 0; i < c.numOfvNodesPerpNodes; i++ {
		delete(c.vNodeID2pNodeID, vNodes[i])
	}
	delete(c.pNodeID2vNodeID, pNodeID)
	c.updateSortedHashes()

}



func (c *Consistent) getVNodes() []uint64 {
	c.RLock()
	defer c.RUnlock()
	var m []uint64
	for k := range c.pNodeID2vNodeID {
		m = append(m, k)
	}
	return m
}

func (c *Consistent) Get(pNodeID uint64) (uint64, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.vNodeID2pNodeID) == 0 {
		return 0, ErrEmptyCircle
	}
	hashValue := c.hash(pNodeID)
	i := c.search(hashValue)

	return c.vNodeID2pNodeID[c.sortedHashes[i]], nil
}

func (c *Consistent) search(pNodeID uint64) (i int) {
	f := func(x int) bool {
		return c.sortedHashes[x] > pNodeID
	}
	i = sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return
}


// Input: physical address
// Output: physical address
func (c *Consistent) GetN(name uint64, n int) ([]uint64, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.vNodeID2pNodeID) == 0 {
		return nil, ErrEmptyCircle
	}

	if len(c.pNodeID2vNodeID) < int(n) {
		n = len(c.pNodeID2vNodeID)
	}

	var (
		hashValue   = c.hash(name)
		i     = c.search(hashValue)
		start = i
		res   = make([]uint64, 0, n)
		elem  = c.vNodeID2pNodeID[c.sortedHashes[i]]
	)

	res = append(res, elem)

	if len(res) == n {
		return res, nil
	}

	for i = start + 1; i != start; i++ {
		if i >= len(c.sortedHashes) {
			i = 0
		}
		elem = c.vNodeID2pNodeID[c.sortedHashes[i]]
		if !sliceContainsMember(res, elem) {
			res = append(res, elem)
		}
		if len(res) == n {
			break
		}
	}

	return res, nil
}



func (c *Consistent) hash(pNodeID uint64) uint64 {
	h := fnv.New64a()
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, pNodeID)
	h.Write(bs)
	return h.Sum64()
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	
	if cap(c.sortedHashes)/(c.numOfvNodesPerpNodes*4) > len(c.vNodeID2pNodeID) {
		hashes = nil
	}
	for k := range c.vNodeID2pNodeID {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

func sliceContainsMember(set []uint64, member uint64) bool {
	for _, m := range set {
		if m == member {
			return true
		}
	}
	return false
}