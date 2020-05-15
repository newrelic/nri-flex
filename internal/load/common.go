package load

import (
	"sync"
	"time"
)

// TimestampMs create a timestamp in milliseconds
func TimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// FlexStatusCounter count internal metrics
var FlexStatusCounter = struct {
	sync.RWMutex
	M map[string]int
}{M: make(map[string]int)}

// StatusCounterIncrement increment the status counter for a particular key
func StatusCounterIncrement(key string) {
	FlexStatusCounter.Lock()
	FlexStatusCounter.M[key]++
	FlexStatusCounter.Unlock()
}

// StatusCounterRead the status counter for a particular key
func StatusCounterRead(key string) int {
	FlexStatusCounter.Lock()
	value := FlexStatusCounter.M[key]
	FlexStatusCounter.Unlock()
	return value
}

// Refresh Helper function used for testing
func Refresh() {
	FlexStatusCounter.M = make(map[string]int)
	FlexStatusCounter.M["EventCount"] = 0
	FlexStatusCounter.M["EventDropCount"] = 0
	FlexStatusCounter.M["ConfigsProcessed"] = 0
	Args.ConfigDir = ""
	Args.ConfigFile = ""
	Args.ContainerDiscovery = false
	Args.ContainerDiscoveryDir = ""
}

// ConfigMutex anything writing to core Flex config, lock and unlock
var ConfigMutex = sync.RWMutex{}
