package fbReport

import (
	"sync"
)

var staticMapLock = sync.Mutex{}
var staticMap = make(map[string]string)

func GetStaticData(mapKey string) string {
	staticMapLock.Lock()
	defer staticMapLock.Unlock()
	return staticMap[mapKey]
}

func SetStaticData(key string, val string) {
	staticMapLock.Lock()
	defer staticMapLock.Unlock()
	staticMap[key] = val
}

func DelStaticData(key string) {
	staticMapLock.Lock()
	defer staticMapLock.Unlock()
	delete(staticMap, key)
}
