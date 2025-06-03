package executors

import "sync"

type atomicInt struct {
	mutex sync.Mutex
	val   int
}

func (this *atomicInt) SetValue(newValue int) {
	this.mutex.Lock()
	this.val = newValue
	this.mutex.Unlock()
}

func (this *atomicInt) GetValue() int {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.val
}
