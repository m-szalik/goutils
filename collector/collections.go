package collector

type Collection[T interface{}] interface {
	Add(element T)
	GetRange() []*T
	Remove(element T) int
	Length() int
	Get(index int) *T
}
