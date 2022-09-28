package mmm

import (
	"github.com/nosnaws/tiam/board"
	z "github.com/nosnaws/tiam/zobrist"
)

const tableSize uint64 = 10067 // prime numbers are better apparently

type Cache struct {
	m        [tableSize]CacheEntry
	hasher   z.ZobristTable
	adjuster int
}

func createCache(b *board.FastBoard, depthAdjuster int) *Cache {
	return &Cache{
		m:        [tableSize]CacheEntry{},
		hasher:   z.InitializeZobristTable(int(b.Height), int(b.Width)),
		adjuster: depthAdjuster,
	}
}

func (c *Cache) setAdjuster(adj int) {
	c.adjuster = adj
}

func (c *Cache) getIndex(k z.Key) uint64 {
	return uint64(k) % tableSize
}

func (c *Cache) getEntry(b *board.FastBoard, id board.SnakeId, depth int) (*CacheEntry, bool) {
	key := z.GetZobristKey(c.hasher, b)
	adjDepth := c.adjuster - depth
	index := c.getIndex(key)

	entry := c.m[index]

	if !entry.isHit || entry.minPlayer != id {
		return nil, false
	}

	if adjDepth >= entry.depth {
		return nil, false
	}

	return &entry, true
}

func (c *Cache) isCollision(index uint64) bool {
	return c.m[index].isHit
}

func (c *Cache) addLowerBound(b *board.FastBoard, value float64, id board.SnakeId, depth int) {
	key := z.GetZobristKey(c.hasher, b)
	index := c.getIndex(key)

	c.m[index] = CacheEntry{
		value:     value,
		depth:     c.adjuster - depth,
		minPlayer: id,
		etype:     lowType,
		isHit:     true,
	}
}

func (c *Cache) addUpperBound(b *board.FastBoard, value float64, id board.SnakeId, depth int) {
	key := z.GetZobristKey(c.hasher, b)
	index := c.getIndex(key)

	c.m[index] = CacheEntry{
		value:     value,
		depth:     c.adjuster - depth,
		minPlayer: id,
		etype:     upType,
		isHit:     true,
	}
}

func (c *Cache) addExact(b *board.FastBoard, value float64, id board.SnakeId, depth int) {
	key := z.GetZobristKey(c.hasher, b)
	index := c.getIndex(key)

	c.m[index] = CacheEntry{
		value:     value,
		depth:     c.adjuster - depth,
		minPlayer: id,
		etype:     exactType,
		isHit:     true,
	}
}

type entryType string

const (
	upType    entryType = "u"
	lowType             = "l"
	exactType           = "e"
)

type CacheEntry struct {
	value     float64
	minPlayer board.SnakeId
	etype     entryType
	depth     int
	isHit     bool
}

func (ce *CacheEntry) isUpperBound() bool {
	if ce.etype == upType || ce.etype == exactType {
		return true
	}

	return false
}

func (ce *CacheEntry) isLowerBound() bool {
	if ce.etype == lowType || ce.etype == exactType {
		return true
	}

	return false

}
