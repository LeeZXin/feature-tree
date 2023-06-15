package manage

import (
	"github.com/LeeZXin/feature-tree/util"
	"sync"
)

var (
	scriptFeatureConfigMap = make(map[string]*util.CachedScript, 8)
	scriptFeatureConfigMu  = sync.RWMutex{}
)

func RefreshScriptFeatureConfigMap(configMap map[string]*util.CachedScript) {
	if configMap == nil {
		return
	}
	scriptFeatureConfigMu.Lock()
	defer scriptFeatureConfigMu.Unlock()
	scriptFeatureConfigMap = configMap
}

func LoadScriptFeatureConfig(featureKey string) (*util.CachedScript, bool) {
	scriptFeatureConfigMu.RLock()
	defer scriptFeatureConfigMu.RUnlock()
	val, ok := scriptFeatureConfigMap[featureKey]
	return val, ok
}
