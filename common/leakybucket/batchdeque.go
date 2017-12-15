package leakybucket

import "errors"

type deque struct {
	data []Batch
}

func newDeque() *deque {
	return &deque{
		data: make([]Batch, 0),
	}
}

func (dq *deque) size() int {
	return len(dq.data)
}

func (dq *deque) shift() {
	dq.data = dq.data[1:]
}

func (dq *deque) append(batch Batch) {
	dq.data = append(dq.data, batch)
}

func (dq *deque) roll(batch Batch, limit int) {
	for dq.size() >= limit {
		dq.shift()
	}

	dq.append(batch)
}

func (dq *deque) last() (Batch, error) {
	if dq.size() == 0 {
		return Batch{}, errors.New("Cannot get last on empty deque")
	}

	return dq.data[dq.size()-1], nil
}

func (dq *deque) first() (Batch, error) {
	if dq.size() == 0 {
		return Batch{}, errors.New("Cannot get first on empty deque")
	}

	return dq.data[0], nil
}

func (dq *deque) iterate() []Batch {
	return dq.getData()
}

func (dq *deque) getData() []Batch {
	cpy := make([]Batch, dq.size())
	copy(cpy, dq.data)
	return cpy
}
