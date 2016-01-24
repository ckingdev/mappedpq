package minpq

import (
	"fmt"
)

type node struct {
	val      interface{}
	priority float32
}

// MinPQ represents a min-priority queue implemented with a heap and a map from
// value to its place in the heap to allow for constant-time checks for
// membership and log-time updates to a value's priority.
type MinPQ struct {
	heapDegree int
	heap       []node
	indices    map[interface{}]int
	size       int
}

// New creates a new MinPQ using heapDegree as the degree of the underlying
// d-ary heap. sizeHint is supplied as the capacity to the slice holding the
// heap and a size hint to the map of value to location in the heap. sizeHint
// may be safely set to zero in the absence of a better choice.
func New(heapDegree, sizeHint int) *MinPQ {
	return &MinPQ{
		heapDegree: heapDegree,
		heap:       make([]node, 0, sizeHint),
		indices:    make(map[interface{}]int, sizeHint),
	}
}

func (p *MinPQ) parent(ind int) int {
	return int((float32(ind) - 1) / float32(p.heapDegree))
}

func (p *MinPQ) firstChild(ind int) int {
	return ind*p.heapDegree + 1
}

func (p *MinPQ) heapUp(ind int) {
	if ind == 0 {
		return
	}
	parentInd := p.parent(ind)
	if p.heap[ind].priority < p.heap[parentInd].priority {
		p.indices[p.heap[ind].val], p.indices[p.heap[parentInd].val] = p.indices[p.heap[parentInd].val], p.indices[p.heap[ind].val]
		p.heap[ind], p.heap[parentInd] = p.heap[parentInd], p.heap[ind]
		p.heapUp(parentInd)
	}
}

// Insert adds val to the priority queue with the given priority. It does NOT
// check whether the value is already in the priority queue. It is the user's
// duty to check this- the data structure will not be valid after inserting a
// value already contained in the queue.
func (p *MinPQ) Insert(val interface{}, priority float32) {
	p.indices[val] = p.size
	if p.size+1 >= len(p.heap) {
		p.heap = append(p.heap, node{
			val:      val,
			priority: priority,
		})
	} else {
		p.heap[p.size] = node{
			val:      val,
			priority: priority,
		}
	}
	p.size++
	p.heapUp(p.size - 1)
}

func (p *MinPQ) heapDown(ind int) {
	fChildInd := p.firstChild(ind)
	if fChildInd >= p.size {
		return
	}
	minInd := fChildInd
	minPri := p.heap[fChildInd].priority
	for i := 1; i < p.heapDegree && (fChildInd+1) < p.size; i++ {
		iterPri := p.heap[fChildInd+i].priority
		if iterPri < minPri {
			minPri = iterPri
			minInd = fChildInd + i
		}
	}
	if p.heap[ind].priority > minPri {
		p.indices[p.heap[ind].val], p.indices[p.heap[minInd].val] = p.indices[p.heap[minInd].val], p.indices[p.heap[ind].val]
		p.heap[ind], p.heap[minInd] = p.heap[minInd], p.heap[ind]
		p.heapDown(minInd)
	}
}

// Pop returns the value with the minimum priority and that priority, removing
// the value from the queue.
func (p *MinPQ) Pop() (interface{}, float32) {
	if p.size == 0 {
		return nil, 0
	}
	ret := p.heap[0]
	delete(p.indices, p.heap[0].val)
	p.size--
	if p.size == 0 {
		return ret.val, ret.priority
	}
	p.heap[0] = p.heap[p.size]
	p.indices[p.heap[0].val] = 0
	p.heapDown(0)
	return ret.val, ret.priority
}

// Empty returns whether the queue is empty or not.
func (p *MinPQ) Empty() bool {
	return p.size == 0
}

// Size returns the size of the queue.
func (p *MinPQ) Size() int {
	return p.size
}

// Contains checks whether the queue contains the given value. Note: if you
// insert a value that is already contained in the heap, the data structure will
// be invalid. Check first unless an invariant allows you to safely ignore this
// possibility.
func (p *MinPQ) Contains(val interface{}) bool {
	_, ok := p.indices[val]
	return ok
}

// CurrentPriority returns the priority of the given value in the queue and
// whether that value was actually found in the queue.
func (p *MinPQ) CurrentPriority(val interface{}) (float32, bool) {
	ind, ok := p.indices[val]
	if !ok {
		return 0, false
	}
	return p.heap[ind].priority, true
}

// UpdatePriority updates the priority associated with the given value to the
// given priority and heapifies as appropriate. An error will be returned if
// the supplied value is not found in the heap.
func (p *MinPQ) UpdatePriority(val interface{}, priority float32) error {
	ind, ok := p.indices[val]
	if !ok {
		return fmt.Errorf("Unable to find value in priority queue.")
	}
	oldPri := p.heap[ind].priority
	p.heap[ind].priority = priority
	if oldPri < priority {
		p.heapDown(ind)
	} else {
		p.heapUp(ind)
	}
	return nil
}
