package tests

var ValidateTestDataPriority = []TestData{

	// 算法优先级测试
	{`([foo] > [bar]) and ([foo] > [bar]) and ([foo] > [bar])`, map[string]interface{}{"foo": "222", "bar": "111"}, true, false},
	{`[foo1] > [bar1] or [foo2] > [bar2] and [foo3] > [bar3]`, map[string]interface{}{"foo1": "222", "bar1": "111", "foo2": "222", "bar2": "111", "foo3": "222", "bar3": "111"}, true, false},
	{`[foo1] > [bar1] and [foo2] > [bar2] or [foo3] > [bar3]`, map[string]interface{}{"foo1": "222", "bar1": "111", "foo2": "222", "bar2": "111", "foo3": "222", "bar3": "111"}, true, false},
	{`[foo1] >= [bar1] and [foo2] > [bar2] or [foo3] < [bar3]`, map[string]interface{}{"foo1": "222", "bar1": "111", "foo2": "222", "bar2": "111", "foo3": "222", "bar3": "111"}, true, false},
}
