package grizzly

// IntegerSet represents a set of integers using a map for fast lookup
type IntegerSet struct {
	set map[int]struct{}
}

// NewIntegerSet creates a new IntegerSet
func NewIntegerSet() *IntegerSet {
	return &IntegerSet{set: make(map[int]struct{})}
}

// AddIntegerSet Add inserts an integer into the set
func (s *IntegerSet) AddIntegerSet(value int) {
	s.set[value] = struct{}{}
}

// SetContainsInteger Contains checks if an integer is in the set
func (s *IntegerSet) SetContainsInteger(value int) bool {
	_, exists := s.set[value]
	return exists
}
