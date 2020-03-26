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

var validTestData = []struct {
	cond   string
	args   map[string]interface{}
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
	{"[foo][unknow] == true and [bar] == true", map[string]interface{}{"foo.dfs": true, "bar": true}, false, true},
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

	// =~
	{"[status] =~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, true, false},
	{"[status] =~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, false, false},

	// !~
	{"[status] !~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, false, false},
	{"[status] !~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, true, false},

	{`foo > bar and foo > bar and foo > bar`, map[string]interface{}{"foo": "222", "bar": "111"}, true, false},
}

var validTestData2 = []struct {
	cond   string
	args   map[string]interface{}
	result bool
	isErr  bool
}{
	{`$foo1 > $bar1 and ($foo2 > $bar2 and $foo3) > $bar3`, map[string]interface{}{"foo": "222", "bar": "111"}, true, false},
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

func TestValid(t *testing.T) {

	var (
		expr Expr
		err  error
		r    bool
	)

	for _, td := range validTestData2 {
		t.Log("--------")
		t.Logf("Parsing: %s", td.cond)

		p := NewParser(strings.NewReader(td.cond))
		expr, err = p.Parse()
		t.Logf("Expression: %s", expr)
		if err != nil {
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

	// Valid
	//s := "$1 > 3 OR (\"OFF\" == $0)"
	// s:= "true"
	// s:= "false"
	// s:= "$1 > true"
	// s:= "3 == true"
	// s:= "$1 > $2"
	//s := "$0 > 3 AND (78 > $0) AND ($0 >= -3 AND $1 < 20.3) OR ($2 > 10) AND ($1 != 44) AND ($2 <= 900) AND ($3 == \"ACTIVE\" OR $3 == \"IDLE\") OR $3 == $1 OR $3 == false"
	//s := "(P0 == -3) AND -100 >= P1"
	//s := "78 > P0 AND (P0 >= -3 AND P1 < 20.3) OR (P2 > 10) AND (P1 != 44) AND (P3 <= 900) AND (P3 == \"ACTIVE\" OR P3 == \"IDLE\") OR P3 == P1 OR P3 == false"

	// Invalid
	//s := "($1 >= -3 AND $1 < 20.3) OR ($2 >= 10s) AND ($3 == \"ACTIVE\" OR $3 == \"IDLE\") OR $3 == $1 OR $3 == false OR $5"

	/*
		p := NewParser(strings.NewReader(s))
		t.Log("Parsing...")
		t.Log(s)
		expr, err := p.Parse()
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}

		t.Log("Evaluating...")
		t.Log(expr)
		//TODO: test case with empty args: matched, err := Evaluate(expr)
		matched, err := Evaluate(expr, "OFF", 56)
		t.Log("Analyzing...")
		if err != nil {
			t.Error(err.Error())
			t.FailNow()
		}
		if !matched {
			t.Errorf("Expected matched=true, but got %v", matched)
			t.Fail()
		}
	*/
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
//如果条件表达式不支持短路操作的话，需要耗时 2200ns/op; 所以，短路求值支持还是非常必要的
func BenchmarkCondition(b *testing.B) {
	expFunc := func() func() bool {
		cond := "($i < $j) and (\"$j\" < \"$k\") " //and $k < $l and $l < $m and $m < $n"
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
	for n := 0; n < 1; n++ {
		expFunc()
	}
}
