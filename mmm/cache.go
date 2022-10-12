package mmm

import (
	"github.com/nosnaws/tiam/board"
	z "github.com/nosnaws/tiam/zobrist"
)

// medium size 1000003 - about 200 MBs
// larger size 3000017
const tableSize uint64 = 3000017 // prime numbers are better apparently

type Cache struct {
	m       [tableSize]CacheEntry
	hasher  z.ZobristTable
	curMax  int
	curTurn int
}

func CreateCache(b *board.FastBoard, depthAdjuster int) *Cache {
	return &Cache{
		m:      [tableSize]CacheEntry{},
		hasher: z.InitializeZobristTable(int(b.Height), int(b.Width)),
		curMax: depthAdjuster,
	}
}

func (c *Cache) setCurMax(adj int) {
	c.curMax = adj
}

func (c *Cache) SetCurTurn(t int) {
	c.curTurn = t
}

func (c *Cache) getIndex(k z.Key) uint64 {
	return uint64(k) % tableSize
}

func (c *Cache) getEntry(b *board.FastBoard, minId board.SnakeId, depth int) (*CacheEntry, bool) {
	key := z.GetZobristKey(c.hasher, b)
	adjDepth := c.curMax - depth
	index := c.getIndex(key)

	entry := c.m[index]

	if !entry.isHit || entry.minPlayer != minId {
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
	c.addEntry(b, value, id, depth, lowType)
}

func (c *Cache) addUpperBound(b *board.FastBoard, value float64, id board.SnakeId, depth int) {
	c.addEntry(b, value, id, depth, upType)
}

func (c *Cache) addExact(b *board.FastBoard, value float64, id board.SnakeId, depth int) {
	c.addEntry(b, value, id, depth, exactType)
}

func (c *Cache) addEntry(b *board.FastBoard, value float64, id board.SnakeId, depth int, eType entryType) {
	key := z.GetZobristKey(c.hasher, b)
	index := c.getIndex(key)
	adjDepth := c.curMax - depth

	if !c.shouldReplace(index, adjDepth) {
		return
	}

	c.m[index] = CacheEntry{
		value:     value,
		depth:     c.curMax - depth,
		minPlayer: id,
		etype:     eType,
		isHit:     true,
		age:       c.curTurn,
	}
}

func (c *Cache) shouldReplace(index uint64, adjDepth int) bool {
	if c.m[index].isHit {
		if c.m[index].age < c.curTurn {
			return true
		}

		if c.m[index].depth < adjDepth {
			return false
		}
	}

	return true
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
	age       int
}

func (ce *CacheEntry) isUpperBound() bool {
	if ce.etype == upType {
		return true
	}

	return false
}

func (ce *CacheEntry) isLowerBound() bool {
	if ce.etype == lowType {
		return true
	}

	return false
}

func (ce *CacheEntry) isExact() bool {
	if ce.etype == exactType {
		return true
	}
	return false
}
