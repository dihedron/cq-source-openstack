package utils

type Single[V any] struct {
	Name V `cq-name:"name"`
}

type Pair[K comparable, V any] struct {
	Key   K `cq-name:"key"`
	Value V `cq-name:"value"`
}

type Tag struct {
	Value string `cq-name:"value"`
}
