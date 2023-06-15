package operator

import (
	"github.com/LeeZXin/feature-tree/util"
)

var (
	Eq = &Operator{
		Operator:      "eq",
		Alias:         "等于",
		ValueSplitter: util.DefaultSplitter,
	}
	Neq = &Operator{
		Operator:      "neq",
		Alias:         "不等于",
		ValueSplitter: util.DefaultSplitter,
	}
	Gt = &Operator{
		Operator:      "gt",
		Alias:         "大于",
		ValueSplitter: util.DefaultSplitter,
	}
	Gte = &Operator{
		Operator:      "gte",
		Alias:         "大于等于",
		ValueSplitter: util.DefaultSplitter,
	}
	Lt = &Operator{
		Operator:      "lt",
		Alias:         "小于",
		ValueSplitter: util.DefaultSplitter,
	}
	Lte = &Operator{
		Operator:      "lte",
		Alias:         "小于等于",
		ValueSplitter: util.DefaultSplitter,
	}
	In = &Operator{
		Operator:      "in",
		Alias:         "包含",
		ValueSplitter: util.CommasSplitter,
	}
	Blank = &Operator{
		Operator:      "blank",
		Alias:         "为空",
		ValueSplitter: util.DefaultSplitter,
	}
	NotBlank = &Operator{
		Operator:      "notBlank",
		Alias:         "不为空",
		ValueSplitter: util.DefaultSplitter,
	}
	RegMatch = &Operator{
		Operator:      "regMatch",
		Alias:         "正则匹配",
		ValueSplitter: util.DefaultSplitter,
	}
	Between = &Operator{
		Operator:      "between",
		Alias:         "范围",
		ValueSplitter: util.CommasSplitter,
	}
	Script = &Operator{
		Operator:      "script",
		Alias:         "脚本",
		ValueSplitter: util.DefaultSplitter,
	}
)

// Operator 运算操作符
type Operator struct {
	Operator      string
	Alias         string
	ValueSplitter util.ValueSplitter
}
