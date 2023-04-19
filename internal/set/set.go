package set

type Set[T comparable] map[T]struct{}

func New[T comparable](v ...T) Set[T] {
	s := Set[T]{}
	for _, x := range v {
		s.Add(x)
	}
	return s
}

func (s Set[T]) Has(v T) bool {
	if s == nil {
		return false
	}
	_, ok := s[v]
	return ok
}

func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}
	for x := range s {
		if !other.Has(x) {
			return false
		}
	}
	return true
}

func (s Set[T]) AddSet(other Set[T]) {
	for x := range other {
		s.Add(x)
	}
}
