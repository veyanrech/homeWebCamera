package client

import "sync"

type RoundBufferQueue struct {
	queue     []string
	nextIndex int
	pointer   int
	snc       sync.RWMutex
	buflen    int
}

func NewRoundBufferQueue(bufferlen int) *RoundBufferQueue {
	return &RoundBufferQueue{
		queue:  make([]string, bufferlen),
		snc:    sync.RWMutex{},
		buflen: bufferlen,
	}
}

func (rbq *RoundBufferQueue) Add(s string) {
	rbq.snc.Lock()
	defer rbq.snc.Unlock()

	rbq.queue[rbq.nextIndex] = s
	nextindex := (rbq.nextIndex + 1) % len(rbq.queue)

	//extend the queue if it si full
	if nextindex == rbq.pointer {
		newq := make([]string, len(rbq.queue)+rbq.buflen)
		copy(newq, rbq.queue)
		rbq.queue = newq
	}

	rbq.nextIndex = (rbq.nextIndex + 1) % len(rbq.queue)
}

func (rbq *RoundBufferQueue) Get() (string, bool) {
	rbq.snc.Lock()
	defer rbq.snc.Unlock()

	if rbq.pointer == rbq.nextIndex {
		return "", false
	}

	s := rbq.queue[rbq.pointer]
	rbq.pointer = (rbq.pointer + 1) % len(rbq.queue)

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
