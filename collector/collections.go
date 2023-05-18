package collector

type Collection[T interface{}] interface {
	Add(elements ...T)
	AsSlice() []*T
	Remove(elements ...T) int
	Length() int
	Get(index int) *T
}
