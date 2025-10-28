package collector

type Collection[T any] interface {
	Add(elements ...T) int
	Remove(elements ...T) int
	Contains(element T) bool
	AsSlice() []*T
	Length() int
}

type IndexableCollection[T any] interface {
	Collection[T]
	Get(index int) *T
}
