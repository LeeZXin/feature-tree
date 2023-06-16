package tree

import (
	"errors"
	"github.com/LeeZXin/feature-tree/cast"
	"github.com/LeeZXin/feature-tree/operator"
	"github.com/LeeZXin/feature-tree/util"
	"github.com/shopspring/decimal"
	"regexp"
	"sync"
)

// 表达式比较器
var (
	regCache = sync.Map{}
)

type compileRegFunc func() *regexp.Regexp

// compileReg 编译正则并缓存
func compileReg(expr string) compileRegFunc {
	var (
		wg sync.WaitGroup
		f  compileRegFunc
	)
	wg.Add(1)
	i, loaded := regCache.LoadOrStore(expr, compileRegFunc(func() *regexp.Regexp {
		wg.Wait()
		return f()
	}))
	if loaded {
		return i.(compileRegFunc)
	}
	compile, err := regexp.Compile(expr)
	if err == nil {
		f = func() *regexp.Regexp {
			return compile
		}
		regCache.Store(expr, f)
	} else {
		f = func() *regexp.Regexp {
			return nil
		}
		regCache.Delete(expr)
	}
	wg.Done()
	return f
}

// StringFeatureHandler 字符串处理器
type StringFeatureHandler struct {
	opMap map[*operator.Operator]Comparator[string]
}

// GetSupportedOperators 获取支持的操作符
func (m *StringFeatureHandler) GetSupportedOperators() []*operator.Operator {
	ret := make([]*operator.Operator, 0, len(m.opMap))
	for key := range m.opMap {
		ret = append(ret, key)
	}
	return ret
}

// GetDataType 支持的数据类型
func (m *StringFeatureHandler) GetDataType() string {
	return "string"
}

// Handle 实际处理逻辑
func (m *StringFeatureHandler) Handle(value *util.StringValue, operator *operator.Operator, userValue any, ctx *FeatureAnalyseContext) (bool, error) {
	actual := cast.ToString(userValue)
	targets := operator.ValueSplitter.SplitValue(value.Value)
	return m.opMap[operator](actual, targets), nil
}

func NewStringFeatureHandler() FeatureHandler {
	return &StringFeatureHandler{
		opMap: map[*operator.Operator]Comparator[string]{
			operator.Eq: func(actual string, targets []string) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual == targets[0]
			},
			operator.In: func(actual string, targets []string) bool {
				if targets == nil {
					return false
				}
				for _, target := range targets {
					if target == actual {
						return true
					}
				}
				return false
			},
			operator.Neq: func(actual string, targets []string) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual != targets[0]
			},
			operator.Blank: func(actual string, targets []string) bool {
				return actual == ""
			},
			operator.NotBlank: func(actual string, targets []string) bool {
				return actual != ""
			},
			operator.RegMatch: func(actual string, targets []string) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				reg := compileReg(targets[0])()
				if reg == nil {
					return false
				}
				ret := reg.MatchString(actual)
				return ret
			},
		},
	}
}

// NumberFeatureHandler 数字处理器
type NumberFeatureHandler struct {
	opMap map[*operator.Operator]Comparator[decimal.Decimal]
}

// GetSupportedOperators 获取支持的操作符
func (m *NumberFeatureHandler) GetSupportedOperators() []*operator.Operator {
	ret := make([]*operator.Operator, 0, len(m.opMap))
	for key := range m.opMap {
		ret = append(ret, key)
	}
	return ret
}

// GetDataType 支持的数据类型
func (m *NumberFeatureHandler) GetDataType() string {
	return "number"
}

// Handle 实际处理逻辑
func (m *NumberFeatureHandler) Handle(value *util.StringValue, operator *operator.Operator, userValue any, ctx *FeatureAnalyseContext) (bool, error) {
	actual := cast.ToString(userValue)
	targets := operator.ValueSplitter.SplitValue(value.Value)
	actualDecimal, err := decimal.NewFromString(actual)
	if err != nil {
		return false, nil
	}
	targetsDecimal := make([]decimal.Decimal, 0, len(targets))
	for _, target := range targets {
		targetDecimal, err := decimal.NewFromString(target)
		if err != nil {
			return false, nil
		}
		targetsDecimal = append(targetsDecimal, targetDecimal)
	}
	return m.opMap[operator](actualDecimal, targetsDecimal), nil
}

func NewNumberFeatureHandler() FeatureHandler {
	return &NumberFeatureHandler{
		opMap: map[*operator.Operator]Comparator[decimal.Decimal]{
			operator.Eq: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual.Equal(targets[0])
			},
			operator.In: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil {
					return false
				}
				for _, target := range targets {
					if target.Equal(actual) {
						return true
					}
				}
				return false
			},
			operator.Neq: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return !actual.Equal(targets[0])
			},
			operator.Gt: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual.GreaterThan(targets[0])
			},
			operator.Gte: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual.GreaterThanOrEqual(targets[0])
			},
			operator.Lt: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual.LessThan(targets[0])
			},
			operator.Lte: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) == 0 {
					return false
				}
				return actual.LessThanOrEqual(targets[0])
			},
			operator.Between: func(actual decimal.Decimal, targets []decimal.Decimal) bool {
				if targets == nil || len(targets) < 2 {
					return false
				}
				return actual.GreaterThanOrEqual(targets[0]) && actual.LessThanOrEqual(targets[1])
			},
		},
	}
}

// ScriptFeatureHandler 脚本处理器
type ScriptFeatureHandler struct {
}

// GetSupportedOperators 获取支持的操作符
func (m *ScriptFeatureHandler) GetSupportedOperators() []*operator.Operator {
	return []*operator.Operator{
		operator.Script,
	}
}

// GetDataType 支持的数据类型
func (m *ScriptFeatureHandler) GetDataType() string {
	return "script"
}

// Handle 实际处理逻辑
func (m *ScriptFeatureHandler) Handle(value *util.StringValue, operator *operator.Operator, userValue any, ctx *FeatureAnalyseContext) (bool, error) {
	if value == nil || value.CachedScript == nil {
		return false, errors.New("empty script config string value")
	}
	bindings := util.Copy2Bindings(ctx.OriginMessage)
	bindings.Set("userValue", userValue)
	eval, err := value.CachedScript.Eval(bindings)
	if err != nil {
		return false, err
	}
	return cast.ToBool(eval), nil
}

func NewScriptFeatureHandler() FeatureHandler {
	return &ScriptFeatureHandler{}
}

func init() {
	// 注册字符串处理器
	RegisterHandler(NewStringFeatureHandler())
	// 注册数字处理器
	RegisterHandler(NewNumberFeatureHandler())
	// 注册脚本处理器
	RegisterHandler(NewScriptFeatureHandler())
}
