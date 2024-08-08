package client

import "sync"

type RoundBufferQueue struct {
	queue       []string
	itemsUnique map[string]bool
	nextIndex   int
	pointer     int
	snc         sync.RWMutex
	buflen      int
}

func NewRoundBufferQueue(bufferlen int) *RoundBufferQueue {
	return &RoundBufferQueue{
		queue:       make([]string, bufferlen, bufferlen),
		snc:         sync.RWMutex{},
		buflen:      bufferlen,
		itemsUnique: make(map[string]bool),
	}
}

func (rbq *RoundBufferQueue) Add(s string) {
	rbq.snc.Lock()
	defer rbq.snc.Unlock()

	//check if the item is already in the queue
	if _, ok := rbq.itemsUnique[s]; ok {
		return
	}

	rbq.queue[rbq.nextIndex] = s
	rbq.itemsUnique[s] = true
	nextindex := (rbq.nextIndex + 1) % len(rbq.queue)

	//extend the queue if it si full
	if nextindex == rbq.pointer {
		newq := make([]string, len(rbq.queue)+rbq.buflen)
		copy(newq, rbq.queue)
		rbq.queue = newq
	}

	rbq.nextIndex = (rbq.nextIndex + 1) % len(rbq.queue)
}

func (rbq *RoundBufferQueue) TryGet() (string, bool) {
	rbq.snc.RLock()
	defer rbq.snc.RUnlock()

	if rbq.pointer == rbq.nextIndex {
		return "", false
	}

	return rbq.queue[rbq.pointer], true
}

func (rbq *RoundBufferQueue) Get() (string, bool) {
	rbq.snc.Lock()
	defer rbq.snc.Unlock()

	if rbq.pointer == rbq.nextIndex {
		return "", false
	}

	s := rbq.queue[rbq.pointer]
	rbq.pointer = (rbq.pointer + 1) % len(rbq.queue)

	//remove the item from the unique map
	delete(rbq.itemsUnique, s)

	return s, true
}

func (rbq *RoundBufferQueue) UnprocessedLen() int {
	rbq.snc.RLock()
	defer rbq.snc.RUnlock()

	if rbq.pointer == rbq.nextIndex {
		return 0
	}

	if rbq.pointer < rbq.nextIndex {
		return rbq.nextIndex - rbq.pointer
	}

	return len(rbq.queue) - rbq.pointer + rbq.nextIndex
}
