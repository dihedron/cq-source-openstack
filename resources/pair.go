package resources

type Pair[K comparable, V any] struct {
	Key   K `cq-name:"key"`
	Value V `cq-name:"value"`
}
