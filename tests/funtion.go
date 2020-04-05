package tests

var ValidateTestDataFunction = []TestData{

	// 返回值为字符串的函数
	{`[funcString] == "test"`, map[string]interface{}{"funcString": func() string { return "test" }}, true, false},
	{`"test" == [funcString]`, map[string]interface{}{"funcString": func() string { return "test" }}, true, false},

	// 返回值为float64的函数
	{`[funcFloat64] == 3`, map[string]interface{}{"funcFloat64": func() float64 { return 3 }}, true, false},
	{`3 == [funcFloat64]`, map[string]interface{}{"funcFloat64": func() float64 { return 3 }}, true, false},

	// 返回值为bool值的函数
	{`[funcBool] == true`, map[string]interface{}{"funcBool": func() bool { return true }}, true, false},
	{`false == [funcBool]`, map[string]interface{}{"funcBool": func() bool { return true }}, false, false},

	// 返回值为[]string 的函数
	{`"test" in [funcSliceString]`, map[string]interface{}{"funcSliceString": func() []string { return []string{"test"} }}, true, false},
	{`[funcSliceString] contains "test"`, map[string]interface{}{"funcSliceString": func() []string { return []string{"test"} }}, true, false},

	// 返回值为[]float64 的函数
	{`5 in [funcSliceFloat64]`, map[string]interface{}{"funcSliceFloat64": func() []float64 { return []float64{5} }}, true, false},
	// 返回值为空slice的函数
	{`5 in [funcSliceFloat64]`, map[string]interface{}{"funcSliceFloat64": func() []float64 { return []float64{} }}, false, false},
	// 返回值为nil的slice的函数
	{`5 in [funcSliceFloat64]`, map[string]interface{}{"funcSliceFloat64": func() []float64 { return nil }}, false, false},
	{`[funcSliceFloat64] contains 5`, map[string]interface{}{"funcSliceFloat64": func() []float64 { return []float64{5} }}, true, false},

	// 函数为nil时的错误处理，不视为错误
	{`[funcString] == "test"`, map[string]interface{}{"funcString": nil}, false, false},

	// 接口提相关的测试
	// 函数作为结构体的属性的测试, 注意： 函数必须是public的
	{`5 in [FuncSliceFloat64]`, struct{ FuncSliceFloat64 func() []float64 }{FuncSliceFloat64: func() []float64 { return []float64{5} }}, true, false},

	// 函数定义为nil的测试

	{"$Func", struct{ Func func() bool }{}, false, true},
	{"$Func", struct{ Func func() string }{}, false, true},
	{"$Func", struct{ Func func() float64 }{}, false, true},
	{"$Func", struct{ Func func() []string }{}, false, true},
	{"$Func", struct{ Func func() []float64 }{}, false, true},
}
