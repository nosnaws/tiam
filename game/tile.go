package game

type SnakeId uint8
type TileIndex uint16

type Tile struct {
	flags uint8
	id    SnakeId
	idx   TileIndex
}

const SNAKE_HEAD uint8 = 0x06
const SNAKE_BODY_PART uint8 = 0x01
const DOUBLE_STACK_PART uint8 = 0x02
const TRIPLE_STACK_PART uint8 = 0x03
const FOOD uint8 = 0x04
const EMPTY uint8 = 0x05
const KIND_MASK uint8 = 0x07
const IS_HAZARD uint8 = 0x10

func CreateEmptyTile() Tile {
	return Tile{
		flags: EMPTY,
		id:    0,
		idx:   0,
	}
}

func CreateHeadTile(id SnakeId, tailIdx TileIndex) Tile {
	return Tile{
		flags: SNAKE_HEAD,
		id:    id,
		idx:   tailIdx,
	}
}

func CreateBodyTile(id SnakeId, nextIdx TileIndex) Tile {
	return Tile{
		flags: SNAKE_BODY_PART,
		id:    id,
		idx:   nextIdx,
	}
}

func CreateDoubleStackTile(id SnakeId, nextIdx TileIndex) Tile {
	return Tile{
		flags: DOUBLE_STACK_PART,
		id:    id,
		idx:   nextIdx,
	}
}

func CreateTripleStackTile(id SnakeId) Tile {
	return Tile{
		flags: TRIPLE_STACK_PART,
		id:    id,
		idx:   0,
	}
}

// BOARD RELATED

func (t *Tile) IsEmpty() bool {
	return t.flags&KIND_MASK == EMPTY
}

func (t *Tile) IsFood() bool {
	return t.flags&KIND_MASK == FOOD
}

func (t *Tile) IsHazard() bool {
	return t.flags&IS_HAZARD != 0
}

func (t *Tile) SetHazard() {
	t.flags |= IS_HAZARD
}

func (t *Tile) SetFood() {
	t.flags = (t.flags & ^KIND_MASK) | FOOD
}

func (t *Tile) ClearHazard() {
	t.flags &= ^IS_HAZARD
}

// SNAKE RELATED
func (t *Tile) IsSnakeHead() bool {
	return t.flags&KIND_MASK == SNAKE_HEAD || t.IsTripleStack()
}

func (t *Tile) IsSnakeSegment() bool {
	return t.IsSnakeBody() || t.IsDoubleStack() || t.IsTripleStack()
}

func (t *Tile) IsSnakeBodyPart() bool {
	return t.flags&KIND_MASK == SNAKE_BODY_PART
}

func (t *Tile) IsDoubleStack() bool {
	return t.flags&KIND_MASK == DOUBLE_STACK_PART
}

func (t *Tile) IsTripleStack() bool {
	return t.flags&KIND_MASK == TRIPLE_STACK_PART
}

func (t *Tile) IsSnakeBody() bool {
	return t.flags&KIND_MASK == SNAKE_BODY_PART || t.flags&KIND_MASK == DOUBLE_STACK_PART
}

func (t *Tile) IsStacked() bool {
	return t.IsDoubleStack() || t.IsTripleStack()
}

func (t *Tile) SetHead(id SnakeId, tailIdx TileIndex) {
	t.flags = (t.flags & ^KIND_MASK) | SNAKE_HEAD
	t.id = id
	t.idx = tailIdx
}

func (t *Tile) SetBodyPart(id SnakeId, nextIdx TileIndex) {
	t.flags = (t.flags & ^KIND_MASK) | SNAKE_BODY_PART
	t.id = id
	t.idx = nextIdx
}

func (t *Tile) SetDoubleStack(id SnakeId, nextIdx TileIndex) {
	t.flags = (t.flags & ^KIND_MASK) | DOUBLE_STACK_PART
	t.id = id
	t.idx = nextIdx
}

func (t *Tile) Clear() {
	t.flags = (t.flags & ^KIND_MASK) | EMPTY
	t.id = 0
	t.idx = 0
}

func (t *Tile) GetSnakeId() (SnakeId, bool) {
	if t.IsSnakeSegment() || t.IsSnakeHead() {
		return t.id, true
	}
	return 0, false
}

func (t *Tile) GetIdx() TileIndex {
	return t.idx
}
