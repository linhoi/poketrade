package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i,j int) bool{
	return x[i] < x[j]
}

func (x units) Swap (i,j int) {
	x[i],x[j] = x[j], x[i]
}

type Consistent struct {
	circle map[uint32]string
	sortedHashes units
	VirtualNode int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string),
		VirtualNode:  16,
	}
}

func (c *Consistent) generateKey(ele string, index int) string {
	return ele + strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key)  <64 {
		var srcatch [64]byte
		copy(srcatch[:],key)
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle){
		hashes = nil
	}

	for k := range c.circle {
		hashes = append(hashes, k)
	}

	sort.Sort(hashes)
	c.sortedHashes =hashes
}

func (c *Consistent) add(element string){
	for i :=0; i < c.VirtualNode; i++ {
		c.circle[c.hashKey(c.generateKey(element,i))] = element
	}
	c.updateSortedHashes()
}
func (c *Consistent) Add(element string){
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

func (c *Consistent) Search(key uint32) int {
	searchFunc := func(x int ) bool {
		return c.sortedHashes[x] > key
	}

	i := sort.Search(len(c.sortedHashes), searchFunc)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return "", errors.New("not Data in Hash Circle")
	}
	key := c.hashKey(name)
	i := c.Search(key)
	return  c.circle[c.sortedHashes[i]], nil
}
