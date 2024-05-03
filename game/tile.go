package game

type SnakeId uint16

type Tile struct {
	flags uint8
	id    SnakeId
	idx   uint16

	// for use in turn based evaluation
	tailId  SnakeId
	tailIdx uint16
}

const SNAKE_BODY_PART uint8 = 0x01
const DOUBLE_STACK_PART uint8 = 0x02
const TRIPLE_STACK_PART uint8 = 0x03
const FOOD uint8 = 0x04
const EMPTY uint8 = 0x05
const SNAKE_HEAD uint8 = 0x06
const KIND_MASK uint8 = 0x07
const IS_HAZARD uint8 = 0x10

func CreateEmptyTile() Tile {
	return Tile{
		flags: EMPTY,
		id:    0,
		idx:   0,
	}
}

func CreateHeadTile(id SnakeId, tailIdx uint16) Tile {
	return Tile{
		flags: SNAKE_HEAD,
		id:    id,
		idx:   tailIdx,
	}
}

func CreateBodyTile(id SnakeId, nextIdx uint16) Tile {
	return Tile{
		flags: SNAKE_BODY_PART,
		id:    id,
		idx:   nextIdx,
	}
}

func CreateDoubleStackTile(id SnakeId, nextIdx uint16) Tile {
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
	return t.IsSnakeHead() || t.IsSnakeBodyPart() || t.IsDoubleStack() || t.IsTripleStack()
}

func (t *Tile) IsNonHeadSegment() bool {
	return t.IsSnakeBodyPart() || t.IsDoubleStack() || t.IsTripleStack()
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

func (t *Tile) IsHeadTail() bool {
	//return t.flags&KIND_MASK == HEAD_TAIL || t.IsHeadTailDoubleStack()
	return t.tailId != 0
}

func (t *Tile) SetHeadTail(headSnakeId SnakeId, headSnakeTail uint16, tailSnakeId SnakeId, tailSnakeNext uint16) {
	t.SetHead(headSnakeId, headSnakeTail)
	t.tailId = tailSnakeId
	t.tailIdx = tailSnakeNext
}

func (t *Tile) SetBodyTail(bodyId SnakeId, bodyTail uint16, tailId SnakeId, bodyNext uint16) {
	t.SetBodyPart(bodyId, bodyTail)
	t.tailId = tailId
	t.tailIdx = bodyNext
}

func (t *Tile) SetHead(id SnakeId, tailIdx uint16) {
	t.flags = (t.flags & ^KIND_MASK) | SNAKE_HEAD
	t.id = id
	t.idx = tailIdx
}

func (t *Tile) SetBodyPart(id SnakeId, nextIdx uint16) {
	t.flags = (t.flags & ^KIND_MASK) | SNAKE_BODY_PART
	t.id = id
	t.idx = nextIdx
}

func (t *Tile) SetDoubleStack(id SnakeId, nextIdx uint16) {
	t.flags = (t.flags & ^KIND_MASK) | DOUBLE_STACK_PART
	t.id = id
	t.idx = nextIdx
}

func (t *Tile) SetTripleStack(id SnakeId) {
	t.flags = (t.flags & ^KIND_MASK) | TRIPLE_STACK_PART
	t.id = id
}

func (t *Tile) Clear() {
	t.flags = (t.flags & ^KIND_MASK) | EMPTY
	t.id = 0
	t.idx = 0
	t.tailId = 0
	t.tailIdx = 0
}

func (t *Tile) GetSnakeId() (SnakeId, bool) {
	if t.IsSnakeSegment() || t.IsSnakeHead() {
		return t.id, true
	}
	return 0, false
}

func (t *Tile) GetIdx() uint16 {
	return t.idx
}
