package util

import (
	"github.com/LeeZXin/feature-tree/cast"
	"reflect"
	"strings"
)

// Copy2Bindings 复制map
func Copy2Bindings(bindings map[string]any) Bindings {
	if bindings == nil {
		return make(Bindings)
	}
	ret := make(Bindings, len(bindings))
	for k, v := range bindings {
		ret[k] = v
	}
	return ret
}

type Bindings map[string]any

func (b Bindings) Get(key string) (any, error) {
	ret := b.findInMap(key)
	return ret, nil
}

func (b Bindings) Set(key string, val any) {
	b[key] = val
}

// findInMap 支持jsonPath获取map值 大量反射 path过长 性能低
func (b Bindings) findInMap(key string) any {
	sp := strings.Split(key, ".")
	var (
		ret any
		ok  bool
		m   map[string]any
	)
	for _, str := range sp {
		if ret != nil {
			ref := reflect.ValueOf(ret)
			if ref.Kind() != reflect.Map {
				return nil
			}
			if ref.Type().Key().Kind() != reflect.String {
				return nil
			}
			keys := ref.MapKeys()
			m = make(map[string]any, ref.Len())
			for _, k := range keys {
				value := ref.MapIndex(k).Interface()
				m[k.String()] = value
			}
		} else {
			m = b
		}
		ret, ok = m[str]
		if !ok {
			return nil
		}
	}
	return ret
}

func (b Bindings) GetInt(key string) int {
	ret, _ := b.Get(key)
	return cast.ToInt(ret)
}

func (b Bindings) GetString(key string) string {
	ret, _ := b.Get(key)
	return cast.ToString(ret)
}

func (b Bindings) GetFloat(key string) float64 {
	ret, _ := b.Get(key)
	return cast.ToFloat64(ret)
}

func (b Bindings) GetBool(key string) bool {
	ret, _ := b.Get(key)
	return cast.ToBool(ret)
}

func (b Bindings) PutAll(data map[string]any) {
	if data != nil {
		for k, v := range data {
			b[k] = v
		}
	}
}
