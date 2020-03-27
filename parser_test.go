package conditions

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var invalidTestData = []string{
	"",
	// "[] AND true",
	"A",
	"[var0] == DEMO",
	"[var0] == 'DEMO'",
	"![var0]",
	"[var0] <> `DEMO`",
}

func TestInvalid(t *testing.T) {

	var (
		expr Expr
		err  error
	)

	for _, cond := range invalidTestData {
		t.Log("--------")
		t.Logf("Parsing: %s", cond)

		p := NewParser(strings.NewReader(cond))
		expr, err = p.Parse()
		if err == nil {
			t.Error("Should receive error")
			break
		}
		if expr != nil {
			t.Error("Expression should nil")
			break
		}
	}
}

var validTestData = []struct {
	cond   string
	args   interface{}
	result bool
	isErr  bool
}{
	{"true", nil, true, false},
	{"false", nil, false, false},
	{"false OR true OR false OR false OR true", nil, true, false},
	{"((false OR true) AND false) OR (false OR true)", nil, true, false},
	{"[var0]", map[string]interface{}{"var0": true}, true, false},
	{"[var0]", map[string]interface{}{"var0": false}, false, false},
	{"[var0] > true", nil, false, true},
	{"[var0] > true", map[string]interface{}{"var0": 43}, false, true},
	{"[var0] > true", map[string]interface{}{"var0": false}, false, true},
	{"[var0] and [var1]", map[string]interface{}{"var0": true, "var1": true}, true, false},
	{"[var0] AND [var1]", map[string]interface{}{"var0": true, "var1": false}, false, false},
	{"[var0] AND [var1]", map[string]interface{}{"var0": false, "var1": true}, false, false},
	{"[var0] AND [var1]", map[string]interface{}{"var0": false, "var1": false}, false, false},
	{"[var0] AND false", map[string]interface{}{"var0": true}, false, false},
	{"56.43", nil, false, true},
	{"[var5]", nil, false, true},
	{"[var0] > -100 AND [var0] < -50", map[string]interface{}{"var0": -75.4}, true, false},
	{"[var0]", map[string]interface{}{"var0": true}, true, false},
	{"[var0]", map[string]interface{}{"var0": false}, false, false},
	{"\"OFF\"", nil, false, true},
	{"\"ON\"", nil, false, true},
	{"[var0] == \"OFF\"", map[string]interface{}{"var0": "OFF"}, true, false},
	{"[var0] > 10 AND [var1] == \"OFF\"", map[string]interface{}{"var0": 14, "var1": "OFF"}, true, false},
	{"([var0] > 10) AND ([var1] == \"OFF\")", map[string]interface{}{"var0": 14, "var1": "OFF"}, true, false},
	{"([var0] > 10) AND ([var1] == \"OFF\") OR true", map[string]interface{}{"var0": 1, "var1": "ON"}, true, false},
	{"[foo][dfs] == true and [bar] == true", map[string]interface{}{"foo.dfs": true, "bar": true}, true, false},
	{"[foo][dfs][a] == true and [bar] == true", map[string]interface{}{"foo.dfs.a": true, "bar": true}, true, false},
	{"[@foo][a] == true and [bar] == true", map[string]interface{}{"@foo.a": true, "bar": true}, true, false},

	//原本这里是需要报错的
	//{"[foo][unknow] == true and [bar] == true", map[string]interface{}{"foo.dfs": true, "bar": true}, false, true},

	//XOR
	{"false XOR false", nil, false, false},
	{"false xor true", nil, true, false},
	{"true XOR false", nil, true, false},
	{"true xor true", nil, false, false},

	//NAND
	{"false NAND false", nil, true, false},
	{"false nand true", nil, true, false},
	{"true nand false", nil, true, false},
	{"true NAND true", nil, false, false},

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

	// =~
	{"[status] =~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, true, false},
	{"[status] =~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, false, false},

	// !~
	{"[status] !~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, false, false},
	{"[status] !~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, true, false},

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
	{`[funcSliceFloat64] contains 5`, map[string]interface{}{"funcSliceFloat64": func() []float64 { return []float64{5} }}, true, false},

	// 函数为nil时的错误处理，不视为错误
	{`[funcString] == "test"`, map[string]interface{}{"funcString": nil}, false, false},

	// 接口提相关的测试
	// 函数作为结构体的属性的测试, 注意： 函数必须是public的
	{`5 in [FuncSliceFloat64]`, struct{ FuncSliceFloat64 func() []float64 }{FuncSliceFloat64: func() []float64 { return []float64{5} }}, true, false},
	// 普通结构体属性的测试
	{`"str" == [Foo]`, struct{ Foo string }{Foo: "str"}, true, false},

	// 关于nil值的处理
	{"[nil] or true", map[string]interface{}{}, true, false},
	{`[nil] == ""`, map[string]interface{}{}, true, false},
	//{`[nil] == 0`, map[string]interface{}{}, true, false},

	// nil 不属于任何集合
	{`[nil] in [1,2]`, map[string]interface{}{}, false, false},

	//{`[nil] in []`, map[string]interface{}{}, true, false},
	// 允许不小心把字符串和数字写入一个silce，但是slice的类型取决于第一个元素的类型，后续元素强制转换成第一个元素的类型
	{`0 in [1,"b"]`, map[string]interface{}{}, true, false}, // 这里的"b" 会被转成数字 0，所以，结果为真
	{`1 in [1,"b"]`, map[string]interface{}{}, true, false}, // 这里的"b" 会被转成数字 0
	{`2 in [1,"2"]`, map[string]interface{}{}, true, false}, // 这里的"2" 会被转成数字 2，所以结果为真
	{`"b" in [1,"b"]`, map[string]interface{}{}, false, true},
	{`"b" in ["b",1]`, map[string]interface{}{}, true, false}, // 这里的1会被转换为字符串 1

	//测试不小心把 下面形式的slice当做参数的bug
	{`"a" in ["a"]`, map[string]interface{}{}, true, false},
	//避免表达式不完整导致死循环的bug
	{`"a" in ["a"`, map[string]interface{}{}, false, true},

	// 暂时对优先级支持是错误的，需要通过加括号搞定
	{`([foo] > [bar]) and ([foo] > [bar]) and ([foo] > [bar])`, map[string]interface{}{"foo": "222", "bar": "111"}, true, false},
}

var validTestData2 = []struct {
	cond   string
	args   interface{}
	result bool
	isErr  bool
}{
	// 这里有一个bug，无限循环了，不过这个涉及parse的内容了，稍后研究
	{`"a" in ["a"]`, map[string]interface{}{}, true, false},
}

func TestValid(t *testing.T) {

	var (
		expr Expr
		err  error
		r    bool
	)

	for _, td := range validTestData {
		t.Log("--------")
		t.Logf("Parsing: %s", td.cond)

		p := NewParser(strings.NewReader(td.cond))
		expr, err = p.Parse()
		t.Logf("Expression: %s", expr)
		if err != nil {
			if td.isErr {
				continue
			}
			t.Errorf("Unexpected error parsing expression: %s", td.cond)
			t.Error(err.Error())
			break
		}

		t.Logf("Evaluating with: %#v", td.args)
		r, err = Evaluate(expr, td.args)
		if err != nil {
			if td.isErr {
				continue
			}
			t.Errorf("Unexpected error evaluating: %s", expr)
			t.Error(err.Error())
			break
		}
		if td.isErr {
			t.Error("Expect fail but not")
			break
		}
		if r != td.result {
			t.Errorf("Expected %v, received: %v", td.result, r)
			break
		}
	}
}

func TestExpressionsVariableNames(t *testing.T) {
	cond := "[@foo][a] == true and [bar] == true or [var9] > 10"
	p := NewParser(strings.NewReader(cond))
	expr, err := p.Parse()
	assert.Nil(t, err)

	args := Variables(expr)
	assert.Contains(t, args, "@foo.a", "...")
	assert.Contains(t, args, "bar", "...")
	assert.Contains(t, args, "var9", "...")
	assert.NotContains(t, args, "foo", "...")
	assert.NotContains(t, args, "@foo", "...")
}

//字符串比较本身是比较耗费时间的；下面的测试中，包含所有比较，则26ns/op； 如果只有 i < j的比较7ns/op
//也就是说短路求值效率会高不少
func BenchmarkStringCompare(b *testing.B) {
	expFunc := func() func() bool {
		i, j, k, l, m, n := "string1", "string2", "string3", "string4", "string5", "string6"
		return func() bool {
			return i < j && j < k && k < l && l < m && m < n
		}
	}()
	for n := 0; n < b.N; n++ {
		expFunc()
	}
}

//字符串比较本身是比较耗费时间的；下面的测试中，包含所有比较，如果第一个条件触发短路，则161ns/op；
//如果条件表达式不支持短路操作的话，需要耗时 1100ns/op; 所以，短路求值支持还是非常必要的
func BenchmarkCondition(b *testing.B) {
	expFunc := func() func() bool {
		cond := "($i < $j) and ($j < $k) and ($k < $l) and ($l < $m) and ($m < $n)"
		p := NewParser(strings.NewReader(cond))
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("%v\n", err)
			return nil
		}
		i, j, k, l, m, n := "string1", "string2", "string3", "string4", "string5", "string6"
		data := map[string]interface{}{"i": i, "j": j, "k": k, "l": l, "m": m, "n": n}
		return func() bool {
			ok, err := Evaluate(expr, data)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			return ok
		}
	}()
	if expFunc == nil {
		return
	}
	for n := 0; n < b.N; n++ {
		expFunc()
		//fmt.Printf("%v\n", ok)
	}
}
