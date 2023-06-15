package util

import (
	lua "github.com/yuin/gopher-lua"
	"sync"
)

var (
	defaultExecutor *ScriptExecutor
	mu              = sync.RWMutex{}
)

func init() {
	defaultExecutor, _ = NewScriptExecutor(1000, 10, map[string]lua.LGFunction{})
}

func RegisterScriptExecutor(e *ScriptExecutor) {
	if e == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	defaultExecutor = e
}

func getScriptExecutor() *ScriptExecutor {
	mu.RLock()
	defer mu.RUnlock()
	ret := defaultExecutor
	return ret
}

// CachedScript 缓存脚本
type CachedScript struct {
	//表达式字符串
	Value string
	//表达式
	expr *lua.FunctionProto
	//锁
	mutex sync.RWMutex
}

// Eval 执行脚本
func (c *CachedScript) Eval(bindings map[string]any) (any, error) {
	executor := getScriptExecutor()
	gb := Copy2Bindings(bindings)
	c.mutex.RLock()
	if c.expr != nil {
		c.mutex.RUnlock()
		return executor.Execute(c.expr, gb)
	}
	c.mutex.RUnlock()
	c.mutex.Lock()
	if c.expr == nil {
		expression, err := executor.CompileLua(c.Value)
		if err != nil {
			c.mutex.Unlock()
			return nil, err
		}
		c.expr = expression
	}
	c.mutex.Unlock()
	return executor.Execute(c.expr, gb)
}
