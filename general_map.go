package grizzly

// IntegerSet represents a set of integers using a map for fast lookup
type IntegerSet struct {
	set map[int]struct{}
}

// newIntegerSet creates a new IntegerSet
func newIntegerSet() *IntegerSet {
	return &IntegerSet{set: make(map[int]struct{})}
}

// addIntegerSet Add inserts an integer into the set
func (s *IntegerSet) addIntegerSet(value int) {
	s.set[value] = struct{}{}
}

// setContainsInteger Contains checks if an integer is in the set
func (s *IntegerSet) setContainsInteger(value int) bool {
	_, exists := s.set[value]
	return exists
}

// FloatSet represents a set of integers using a map for fast lookup
type FloatSet struct {
	set map[int]struct{}
}

// newIntegerSet creates a new IntegerSet
func newFloatSet() *FloatSet {
	return &FloatSet{set: make(map[int]struct{})}
}

// addIntegerSet Add inserts an integer into the set
func (s *FloatSet) addFloatSet(value int) {
	s.set[value] = struct{}{}
}

// setContainsInteger Contains checks if an integer is in the set
func (s *FloatSet) setContainsFloat(value int) bool {
	_, exists := s.set[value]
	return exists
}
