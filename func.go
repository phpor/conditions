package conditions

import "fmt"

// support (and only support) no argument function which return bool
func evalFunc(val interface{}, name string) (Expr, error) {
	if val == nil {
		return falseExpr, fmt.Errorf("func %s defined is nil", name)
	}

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
		return falseExpr, err
	}
	return &NumberLiteral{Val: result}, nil
}

func evalFuncString(fun func() string, name string) (Expr, error) {
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
		return falseExpr, err
	}
	return &StringLiteral{Val: result}, nil
}

func evalFuncSliceFloat64(fun func() []float64, name string) (Expr, error) {
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
	return &SliceNumberLiteral{Val: result}, nil
}

func evalFuncSliceString(fun func() []string, name string) (Expr, error) {
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
	return &SliceStringLiteral{Val: result}, nil
}
