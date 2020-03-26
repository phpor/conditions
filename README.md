# 简单规则引擎 conditions

规则引擎对 Golang 的 `struct` 或者 `map` 进行布尔判断求值，返回 `true|false` 。


This package offers a parser of a simple conditions specification language (reduced set of arithmetic/logical operations). The package is mainly created for Flow-Based Programming components that require configuration to perform some operations on the data received from multiple input ports. But it can be used whereever you need externally define some logical conditions on the internal variables.

Additional credits for this package go to [Handwritten Parsers & Lexers in Go](http://blog.gopheracademy.com/advent-2014/parsers-lexers/) by Ben Johnson on [Gopher Academy blog](http://blog.gopheracademy.com) and [InfluxML package from InfluxDB repository](https://github.com/influxdb/influxdb/tree/master/influxql).

## 用法示例 Usage example 
```
package main

import (
	"fmt"
	"strings"

	"github.com/yowenter/conditions"
)

func main() {
	// Our condition to check
	type people struct {
		Name   string
		Height int32
		Male   bool
	}

	s := ` $Name == "test" AND $Height > 100 AND $Male == false`

	// Parse the condition language and get expression
	p := conditions.NewParser(strings.NewReader(s))
	expr, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
		// ...
	}

	// Evaluate expression passing data for $vars
	p1 := map[string]interface{}{"Name": "test", "Height": 180, "Male": false}
	r, err := conditions.Evaluate(expr, p1)
	if err != nil {
		fmt.Println(err)
		// ...
	}
	fmt.Printf("Condition: `%v`, Val: `%v`, Result: `%v`\n", s, p1, r)

	// use struct
	var p2 = people{Name: "test", Height: 200, Male: false}
	r, err = conditions.Evaluate(expr, p2)
	if err != nil {
		fmt.Println(err)
		// ...
	}
	fmt.Printf("Condition: `%v`, Val: `%v`, Result: `%v`\n", s, p2, r)

	// test invalid args . not map or struct.
	r, err = conditions.Evaluate(expr, "")
	if err != nil {
		fmt.Println(err)
		// ...
	}
	fmt.Printf("Condition: `%v`, Val: `%v`, Result: `%v`\n", s, "invalid", r)
}

```

## Where do we use it?

Here is a diagram for a sample FBP flow (created using [FlowMaker](https://github.com/cascades-fbp/flowmaker)). You can see how we configure the ContextA process with a condition via IIP packet.

![](https://raw.githubusercontent.com/yowenter/conditions/master/Example.png)

##todo
1. 目前字符串只支持==操作； 完善字符串的 < > <=  >= 操作
1. 发现算符优先级bug，目前通过()来解决
1. evaluateSubtree 时，对于值找不到或者为nil的情况，需要慎重处理，否则会带来很多的混乱,是否可以有一个nil Literal，允许参与到多数类型的比较和运算
    1. nil 的参与计算稍微复杂一些，因为nil可能是左值，也可能是右值； 原来是根据左值约束右值的，有了nil之后，nil为左值时，就需要参考右值的值。
        每个apply* 函数里面需要处理
    1. 无限循环的bug

