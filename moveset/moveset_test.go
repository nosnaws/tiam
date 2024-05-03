package moveset

import "testing"

func TestCreate(t *testing.T) {
	if Create() != 0 {
		panic("did not create empty moveset")
	}
}

func TestSetLeft(t *testing.T) {
	ms := Create()

	newMs := SetLeft(ms)

	if newMs != 0b00001000 {
		panic("Did not set left bit")
	}
}

func TestSetRight(t *testing.T) {
	ms := Create()

	newMs := SetRight(ms)

	if newMs != 0b00000001 {
		panic("Did not set right bit")
	}
}

func TestSetUp(t *testing.T) {
	ms := Create()

	newMs := SetUp(ms)

	if newMs != 0b00000010 {
		panic("Did not set up bit")
	}
}

func TestSetDown(t *testing.T) {
	ms := Create()

	newMs := SetDown(ms)

	if newMs != 0b00000100 {
		panic("Did not set left bit")
	}
}

func TestHasUp(t *testing.T) {
	ms := Create()

	newMs := SetUp(ms)

	if !HasUp(newMs) {
		panic("Should have Up in move set")
	}
}

func TestHasDown(t *testing.T) {
	ms := Create()

	newMs := SetDown(ms)

	if !HasDown(newMs) {
		panic("Should have Down in move set")
	}
}

func TestHasLeft(t *testing.T) {
	ms := Create()

	newMs := SetLeft(ms)

	if !HasLeft(newMs) {
		panic("Should have Left in move set")
	}
}

func TestHasRight(t *testing.T) {
	ms := Create()

	newMs := SetRight(ms)

	if !HasRight(newMs) {
		panic("Should have Right in move set")
	}
}

func TestIsEmpty(t *testing.T) {
	ms := Create()

	if !IsEmpty(ms) {
		panic("default state is not considered empty")
	}

	ms = SetLeft(ms)

	if IsEmpty(ms) {
		panic("should not be considered empty")
	}
}

func TestToString(t *testing.T) {
	ms := Create()

	ms = SetDown(ms)
	ms = SetLeft(ms)

	if len(ToDirs(ms)) != 2 {
		panic("did not output correct number of dirs")
	}
}
