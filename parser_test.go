package conditions

import (
	"fmt"
	"github.com/phpor/conditions/tests"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validTestDataAll []tests.TestData

func init() {
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataAndOr...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataCollector...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataFunction...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataGtLt...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataNil...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataNumeric...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataPriority...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataRegex...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataStruct...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataSyntax...)
	validTestDataAll = append(validTestDataAll, tests.ValidateTestDataVar...)
}

func init() {
	validTestDataAll2 = append(validTestDataAll2, func() []tests.TestData {
		type st struct {
			Bar bool
		}
		return []tests.TestData{
			{"[foo][unknow] == true and [bar] == true", struct{ Foo st }{Foo: st{Bar: true}}, false, false},
			// todo : 关于结构体的属性的嵌套的支持
			{"[Foo][Bar] == true ", struct{ Foo st }{Foo: st{Bar: true}}, true, false},
			{"[Foo][Bar]", struct{ Foo st }{Foo: st{Bar: true}}, true, false},
		}
	}()...)
}

var validTestDataAll2 = []tests.TestData{

	//{"[foo][unknow] == true and [bar] == true", map[string]interface{}{"foo.dfs": true, "bar": true}, false, false},

}

func TestValid(t *testing.T) {

	var (
		expr Expr
		err  error
		r    bool
	)

	for _, td := range validTestDataAll2 {
		t.Log("--------")
		t.Logf("Parsing: %s", td.Cond)

		p := NewParser(strings.NewReader(td.Cond))
		expr, err = p.Parse()
		t.Logf("Expression: %v", expr)
		if err != nil {
			if td.IsErr {
				continue
			}
			t.Errorf("Unexpected error parsing expression: %s", td.Cond)
			t.Error(err.Error())
			break
		}

		t.Logf("Evaluating with: %#v", td.Args)
		t.Logf("Expect: {result: %v , isErr: %v}", td.Result, td.IsErr)
		r, err = Evaluate(expr, td.Args)
		if err != nil {
			if td.IsErr {
				continue
			}
			t.Errorf("Unexpected error evaluating: %s", expr)
			t.Error(err.Error())
			break
		}
		if td.IsErr {
			t.Error("Expect fail but not")
			break
		}
		if r != td.Result {
			t.Errorf("Expected %v, received: %v", td.Result, r)
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
		cond := "($i < $j) or ($j < $k) and ($k < $l) and ($l < $m) and ($m < $n)"
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
