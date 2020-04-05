package conditions

import "fmt"

// support (and only support) no argument function which return bool
func evalFunc(val interface{}, name string) (Expr, error) {

	fun, ok := val.(func() bool)
	if ok {
		return evalFuncBool(fun, name)
	}
	funString, ok := val.(func() string)
	if ok {
		return evalFuncString(funString, name)
	}
	funFloat64, ok := val.(func() float64)
	if ok {
		return evalFuncFloat64(funFloat64, name)
	}
	funcSliceString, ok := val.(func() []string)
	if ok {
		return evalFuncSliceString(funcSliceString, name)
	}
	funcSliceFloat64, ok := val.(func() []float64)
	if ok {
		return evalFuncSliceFloat64(funcSliceFloat64, name)
	}

	return falseExpr, fmt.Errorf("func %T only can be 'func() bool'", name)
}

func evalFuncBool(fun func() bool, name string) (Expr, error) {
	if fun == nil {
		return falseExpr, fmt.Errorf("func %s:(%T) defined is nil", name, fun)
	}
	var err error
	result := func() bool {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("call func $%s error (return as false): %s", name, r)
			}
		}()
		return fun()
	}()
	if err != nil {
		return falseExpr, err
	}
	return &BooleanLiteral{Val: result}, nil
}

func evalFuncFloat64(fun func() float64, name string) (Expr, error) {
	expr := &NumberLiteral{Val: 0}
	if fun == nil {
		return expr, fmt.Errorf("func %s:(%T) defined is nil", name, fun)
	}
	var err error
	result := func() float64 {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("call func $%s error (return as false): %s", name, r)
			}
		}()
		return fun()
	}()
	if err != nil {
		return expr, err
	}
	expr.Val = result
	return expr, nil
}

func evalFuncString(fun func() string, name string) (Expr, error) {
	expr := &StringLiteral{Val: ""}
	if fun == nil {
		return expr, fmt.Errorf("func %s:(%T) defined is nil", name, fun)
	}
	var err error
	result := func() string {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("call func $%s error (return as false): %s", name, r)
			}
		}()
		return fun()
	}()
	if err != nil {
		return expr, err
	}
	expr.Val = result
	return expr, nil
}

func evalFuncSliceFloat64(fun func() []float64, name string) (Expr, error) {
	expr := &SliceNumberLiteral{Val: []float64{}}
	if fun == nil {
		return expr, fmt.Errorf("func %s:(%T) defined is nil", name, fun)
	}
	var err error
	result := func() []float64 {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("call func $%s error (return as false): %s", name, r)
			}
		}()
		return fun()
	}()
	if err != nil {
		return falseExpr, err
	}
	expr.Val = result
	return expr, nil
}

func evalFuncSliceString(fun func() []string, name string) (Expr, error) {
	expr := &SliceStringLiteral{Val: []string{}}
	if fun == nil {
		return expr, fmt.Errorf("func %s:(%T) defined is nil", name, fun)
	}
	var err error
	result := func() []string {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("call func $%s error (return as false): %s", name, r)
			}
		}()
		return fun()
	}()
	if err != nil {
		return falseExpr, err
	}
	expr.Val = result
	return expr, nil
}
