package collector

type Collection[T any] interface {
	Add(elements ...T) int
	AsSlice() []*T
	Remove(elements ...T) int
	Length() int
}

type IndexableCollection[T any] interface {
	Collection[T]
	Get(index int) *T
}
