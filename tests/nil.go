package tests

var ValidateTestDataNil = []TestData{

	// 关于nil值的处理
	{"[nil] or true", nil, true, false},
	{`[nil] == ""`, nil, true, false},
	{`[nil] == 0`, nil, true, false},

	// nil 不属于任何集合
	{`[nil] in [1,2]`, nil, false, false},
	{`[nil] in {1,2}`, nil, false, false},

	//空集合解析会失败
	{`[nil] in []`, nil, false, true},

	// 关于nil的处理，右值的情况，nil尽量转成左值的类型来比较
	{`true == nil`, nil, false, false},
	{`false == nil`, nil, true, false},
	{`"" == nil`, nil, true, false},
	{`0 == nil`, nil, true, false},
	{`0 > nil`, nil, false, false},
	{`0 >= nil`, nil, true, false},
	{`0 < nil`, nil, false, false},
	{`0 <= nil`, nil, true, false},

	// nil 作为左值时的测试，nil尽量转成右值的类型来比较
	{`nil == 0`, nil, true, false},
	{`nil != 0`, nil, false, false},
	{`true == [nil]`, nil, false, false},

	{`nil >= 0`, nil, true, false},
	{`nil <= 0`, nil, true, false},
	{`nil > 0`, nil, false, false},
	{`nil < 0`, nil, false, false},
}
