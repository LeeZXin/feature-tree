package tree

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeeZXin/feature-tree/manage"
	"github.com/LeeZXin/feature-tree/util"
	lua "github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func TestCreateTree(t *testing.T) {
	manage.RefreshScriptFeatureConfigMap(map[string]*util.CachedScript{
		"user_mock_script": {Value: `return fn("aaaa")`},
	})
	e, _ := util.NewScriptExecutor(500, 10, map[string]lua.LGFunction{
		"fn": func(state *lua.LState) int {
			args := util.GetFnArgs(state)
			state.Push(args[0])
			return 1
		},
	})
	util.RegisterScriptExecutor(e)
	str := `{
				"or": [
					{
						"featureType": "script",
						"featureKey": "user_mock_script",
						"featureName": "脚本测试",      
						"dataType": "string",
						"operator": "eq",      
						"value": "hello"
					},
					{
						"featureType": "message",
						"featureKey": "user_purchase_product_code",
						"featureName": "用户成交产品",      
						"dataType": "script",
						"operator": "script",      
						"value": "params.userValue"
					},
					{
						"featureType": "message",
						"featureKey": "user_purchase_product_code",
						"featureName": "用户成交产品",      
						"dataType": "string",
						"operator": "regMatch",      
						"value": "^s\\d+$"
					},
					{
						"featureType": "message",
						"featureKey": "user_purchase_product_codeddd",
						"featureName": "用户成交产品",      
						"dataType": "number",
						"operator": "between",      
						"value": "1234,124325"
					}
				]
			}`
	var treePlainInfo PlainInfo
	err := json.Unmarshal([]byte(str), &treePlainInfo)
	if err != nil {
		panic(err)
	}
	tree, err := BuildFeatureTree("xx", &treePlainInfo)
	if err != nil {
		panic(err)
	}
	m := map[string]any{
		"user_purchase_product_code":    "ff5hh555",
		"user_purchase_product_codeddd": 1,
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancelFunc()
	for i := 0; i < 10; i++ {
		treeAnalyser := InitTreeAnalyser(BuildFeatureAnalyseContext(tree, m, timeout))
		analyseResult := treeAnalyser.Analyse()
		fmt.Println(analyseResult.GetMissResultDetailDesc())
	}

}

func BenchmarkCreateTree(t *testing.B) {
	manage.RefreshScriptFeatureConfigMap(map[string]*util.CachedScript{
		"user_mock_script": {Value: `return 'hello world'`},
	})
	str := `{
				"or": [
					{
						"featureType": "script",
						"featureKey": "user_mock_script",
						"featureName": "脚本测试",      
						"dataType": "string",
						"operator": "eq",      
						"value": "hello"
					},
					{
						"featureType": "message",
						"featureKey": "user_purchase_product_code",
						"featureName": "用户成交产品",      
						"dataType": "script",
						"operator": "script",      
						"value": "params.userValue"
					},
					{
						"featureType": "message",
						"featureKey": "user_purchase_product_code",
						"featureName": "用户成交产品",      
						"dataType": "string",
						"operator": "regMatch",      
						"value": "^s\\d+$"
					},
					{
						"featureType": "message",
						"featureKey": "user_purchase_product_codeddd",
						"featureName": "用户成交产品",      
						"dataType": "number",
						"operator": "between",      
						"value": "1234,124325"
					}
				]
			}`
	var treePlainInfo PlainInfo
	err := json.Unmarshal([]byte(str), &treePlainInfo)
	if err != nil {
		panic(err)
	}
	tree, err := BuildFeatureTree("xx", &treePlainInfo)
	if err != nil {
		panic(err)
	}
	m := map[string]any{
		"user_purchase_product_code":    "ff5hh555",
		"user_purchase_product_codeddd": 1,
	}
	for i := 0; i < t.N; i++ {
		treeAnalyser := InitTreeAnalyser(BuildFeatureAnalyseContext(tree, m, context.Background()))
		_ = treeAnalyser.Analyse()
	}

}
