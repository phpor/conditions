package tests

var ValidateTestDataGtLt = []TestData{

	// 关于集合的测试
	// 空集合解析失败
	{`1 in []`, map[string]interface{}{}, false, true},

	// IN
	{"[foo] in [foobar]", map[string]interface{}{"foo": "findme", "foobar": []string{"notme", "may", "findme", "lol"}}, true, false},

	// NOT IN
	{"[foo] not in [foobar]", map[string]interface{}{"foo": "dontfindme", "foobar": []string{"notme", "may", "findme", "lol"}}, true, false},

	// IN with array of string
	{`[foo] in ["bonjour", "le monde", "oui"]`, map[string]interface{}{"foo": "le monde"}, true, false},
	{`[foo] in ["bonjour", "le monde", "oui"]`, map[string]interface{}{"foo": "world"}, false, false},

	// NOT IN with array of string
	{`[foo] not in ["bonjour", "le monde", "oui"]`, map[string]interface{}{"foo": "le monde"}, false, false},
	{`[foo] not in ["bonjour", "le monde", "oui"]`, map[string]interface{}{"foo": "world"}, true, false},

	// IN with array of numbers
	{`[foo] in [2,3,4]`, map[string]interface{}{"foo": 4}, true, false},
	{`[foo] in [2,3,4]`, map[string]interface{}{"foo": 5}, false, false},

	// NOT IN with array of numbers
	{`[foo] not in [2,3,4]`, map[string]interface{}{"foo": 4}, false, false},
	{`[foo] not in [2,3,4]`, map[string]interface{}{"foo": 5}, true, false},

	// NOT IN with array of variable slice numbers
	{`[foo] not in [bar]`, map[string]interface{}{"foo": 4, "bar": []float64{2, 3, 4}}, false, false},
	{`[foo] not in [bar]`, map[string]interface{}{"foo": 5, "bar": []float64{2, 3, 4}}, true, false},

	{`[2,3,4] contains [foo] `, map[string]interface{}{"foo": 4}, true, false},
	{`[2,3,4] contains [foo] `, map[string]interface{}{"foo": 5}, false, false},

	// 允许不小心把字符串和数字写入一个silce，但是slice的类型取决于第一个元素的类型，后续元素强制转换成第一个元素的类型
	{`0 in [1,"b"]`, nil, true, false}, // 这里的"b" 会被转成数字 0，所以，结果为真
	{`1 in [1,"b"]`, nil, true, false}, // 这里的"b" 会被转成数字 0
	{`2 in [1,"2"]`, nil, true, false}, // 这里的"2" 会被转成数字 2，所以结果为真
	{`"b" in [1,"b"]`, nil, false, true},
	{`"b" in ["b",1]`, nil, true, false}, // 这里的1会被转换为字符串 1

	// 关于花括号语法的集合
	{`1 in {1,2}`, nil, true, false},
	{`2 in {1,2}`, nil, true, false},
	{`3 in {1,2}`, nil, false, false},
}
