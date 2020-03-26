package main

import (
	"fmt"
	"strings"

	"github.com/phpor/conditions"
)

func main() {
	priority()
}

func priority() {
	s := `$foo1 > $bar1 and $foo2 > $bar2 and $foo3 > $bar3`

	// Parse the condition language and get expression
	p := conditions.NewParser(strings.NewReader(s))
	expr, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
		// ...
	}
	fmt.Printf("\nExpr: %v\n", expr)

	// Evaluate expression passing data for $vars
	p1 := map[string]interface{}{"foo1": "test", "foo2": "test", "foo3": "test", "bar1": "test", "bar2": "test", "bar3": "test"}
	r, err := conditions.Evaluate(expr, p1)
	if err != nil {
		fmt.Println("%v\n", err)
		// ...
	}
	fmt.Printf("\nCondition: `%v`, Val: `%v`, Result: `%v`\n", s, p1, r)

}

func basic() {
	// Our condition to check
	type people struct {
		Name    string
		Height  int32
		Male    bool
		Goods   []string
		IsTrue  func() bool // correct func
		IsFalse func()      // incorrect func
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
	//var p2 = people{Name: "test", Height: 200, Male: false, Goods: []string{"A", "B"}, IsTrue: func() bool {
	//	panic("aaa")
	//	return true
	//}}
	var p2 = people{Name: "test", Height: 200, Male: false, Goods: []string{"A", "B"}, IsFalse: func() {}}
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

	//contains := ` ($Goods CONTAINS "A") AND $Name == "test" AND $IsTrue`
	contains := ` ($Goods CONTAINS "A") AND $Name == "test" AND $IsFalse`

	// Parse the condition language and get expression
	containsP := conditions.NewParser(strings.NewReader(contains))
	containsExpr, err := containsP.Parse()
	if err != nil {
		fmt.Println(err)
		return
		// ...
	}
	r, err = conditions.Evaluate(containsExpr, p2)
	if err != nil {
		fmt.Println(err)
		// ...
	}
	fmt.Printf("Condition: `%v`, Val: `%v`, Result: `%v`\n", contains, p2, r)
}
