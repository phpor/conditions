package tests

var ValidateTestDataStruct = []TestData{
	// 普通结构体属性的测试
	{`"str" == [Foo]`, struct{ Foo string }{Foo: "str"}, true, false},
}
