package collector

// Collection is a mutable generic collection interface.
type Collection[T any] interface {
	// Add inserts elements and returns the number of added entries.
	Add(elements ...T) int
	// Remove deletes matching elements and returns the number of removed entries.
	Remove(elements ...T) int
	// Contains reports whether element exists in the collection.
	Contains(element T) bool
	// AsSlice returns elements as a slice of pointers.
	AsSlice() []*T
	// Length returns current number of elements.
	Length() int
}

// IndexableCollection is a Collection with random access by index.
type IndexableCollection[T any] interface {
	Collection[T]
	Get(index int) *T
}
