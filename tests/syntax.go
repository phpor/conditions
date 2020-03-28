package tests

var ValidateTestDataSyntax = []TestData{
	{`""`, nil, false, true},
	{`"A"`, nil, false, true},
	{`A`, nil, false, true},
	{"[var0] == DEMO", nil, false, true},
	// 不支持单引号引用字符串
	{"[var0] == 'DEMO'", nil, false, true},
	// 不支持反引号引用字符串
	{"[var0] == `DEMO`", nil, false, true},
	{"[var0] <> DEMO", nil, false, true},

	{"true", nil, true, false},
	{"false", nil, false, false},
	{"false OR true OR false OR false OR true", nil, true, false},
	{"((false OR true) AND false) OR (false OR true)", nil, true, false},
	{"[var0]", map[string]interface{}{"var0": true}, true, false},
	{"[var0]", map[string]interface{}{"var0": false}, false, false},

	//测试不小心把 下面形式的slice当做参数的bug
	{`"a" in ["a"]`, nil, true, false},
	//避免表达式不完整导致死循环的bug
	{`"a" in ["a"`, nil, false, true},

	// 识别错误表达式,代码初期，这个被识别为true，且没有错误的，现在修复该bug，作为解析错误处理
	{`true "bb"`, nil, false, true},
	{`true > >`, nil, false, true},

	{"56.43", nil, false, true},
	{"[var5]", nil, false, true},
	{"[var0] > -100 AND [var0] < -50", map[string]interface{}{"var0": -75.4}, true, false},
	{"[var0]", map[string]interface{}{"var0": true}, true, false},
	{"[var0]", map[string]interface{}{"var0": false}, false, false},
	{"\"OFF\"", nil, false, true},
	{"\"ON\"", nil, false, true},
}
