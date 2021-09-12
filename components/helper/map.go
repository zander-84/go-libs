package helper

import "sync"

type ConcurrentMap struct {
	lock sync.RWMutex
	val  map[string]interface{}
}

func NewConcurrentMap() *ConcurrentMap {
	data := new(ConcurrentMap)
	data.val = make(map[string]interface{}, 0)
	return data
}

func (cm *ConcurrentMap) Set(key string, val interface{}) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.val[key] = val
}
func (cm *ConcurrentMap) Get(key string) (interface{}, bool) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	data, ok := cm.val[key]
	return data, ok
}

func (cm *ConcurrentMap) ShouldGetString(key string) string {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(string)
		return res
	}
	return ""
}

func (cm *ConcurrentMap) ShouldGetInt32(key string) int32 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(int32)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetInt64(key string) int64 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(int64)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetInt(key string) int {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(int)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetFloat64(key string) float64 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(float64)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetFloat32(key string) float32 {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(float32)
		return res
	}
	return 0
}

func (cm *ConcurrentMap) ShouldGetBool(key string) bool {
	data, ok := cm.Get(key)
	if ok {
		res, _ := data.(bool)
		return res
	}
	return false
}
