package conditions

import (
	"fmt"
	"reflect"
)

var (
	falseExpr = &BooleanLiteral{Val: false}
)

// Evaluate takes an expr and evaluates it using given args
func Evaluate(expr Expr, args interface{}) (bool, error) {
	if expr == nil {
		return false, fmt.Errorf("Provided expression is nil")
	}

	result, err := evaluateSubtree(expr, args)
	if err != nil {
		return false, err
	}
	switch n := result.(type) {
	case *BooleanLiteral:
		return n.Val, nil
	}
	return false, fmt.Errorf("Unexpected result of the root expression: %#v", result)
}

// evaluateSubtree performs given expr evaluation recursively
func evaluateSubtree(expr Expr, args interface{}) (Expr, error) {
	if expr == nil {
		return falseExpr, fmt.Errorf("Provided expression is nil")
	}

	var (
		err    error
		lv, rv Expr
	)

	switch n := expr.(type) {
	case *ParenExpr:
		return evaluateSubtree(n.Expr, args)
	case *BinaryExpr:
		lv, err = evaluateSubtree(n.LHS, args)
		if err != nil {
			return falseExpr, err
		}
		if n.Op == AND {
			b, _ := getBoolean(lv) // 这里不可能错，因为上面已经正确求值了
			if !b {
				return &BooleanLiteral{Val: false}, nil
			}
		}
		rv, err = evaluateSubtree(n.RHS, args)
		if err != nil {
			return falseExpr, err
		}
		return applyOperator(n.Op, lv, rv)
	case *VarRef:
		index := n.Val
		argsKind := reflect.TypeOf(args).Kind()
		var val interface{}

		switch argsKind {
		case reflect.Map:
			argsMap, ok := args.(map[string]interface{})
			if !ok {
				return falseExpr, fmt.Errorf("Args: `%v` convert to map not ok", args)
			}
			if _, ok := argsMap[index]; !ok {
				return falseExpr, fmt.Errorf("Argument: `%v` not found", index)
			}
			val, _ = argsMap[index]
		case reflect.Struct:
			ps := reflect.ValueOf(args)
			fval := ps.FieldByName(index)
			if !fval.IsValid() {
				return falseExpr, fmt.Errorf("Argument: `%v` not found in args `%v`", index, args)
			}
			val = fval.Interface()
		default:
			return falseExpr, fmt.Errorf("Args: `%v` is not map or struct", args)
		}

		kind := reflect.TypeOf(val).Kind()
		switch kind {
		case reflect.Int:
			return &NumberLiteral{Val: float64(val.(int))}, nil
		case reflect.Int32:
			return &NumberLiteral{Val: float64(val.(int32))}, nil
		case reflect.Int64:
			return &NumberLiteral{Val: float64(val.(int64))}, nil
		case reflect.Float32:
			return &NumberLiteral{Val: float64(val.(float32))}, nil
		case reflect.Float64:
			return &NumberLiteral{Val: float64(val.(float64))}, nil
		case reflect.String:
			return &StringLiteral{Val: val.(string)}, nil
		case reflect.Bool:
			return &BooleanLiteral{Val: val.(bool)}, nil
		case reflect.Slice:
			if v, ok := val.([]string); ok {
				return &SliceStringLiteral{Val: v}, nil
			}
			if v, ok := val.([]float64); ok {
				return &SliceNumberLiteral{Val: v}, nil
			}
			return falseExpr, fmt.Errorf("unsupported argument %s type: %s", n.Val, kind)
		case reflect.Func:
			return evalFunc(val, index)
		}
		return falseExpr, fmt.Errorf("unsupported argument %s type: %s", n.Val, kind)
	}

	return expr, nil
}

// support (and only support) no argument function which return bool
func evalFunc(val interface{}, name string) (Expr, error) {
	fun, ok := val.(func() bool)
	if !ok {
		return falseExpr, fmt.Errorf("func %s only can be 'func() bool'", name)
	}
	if fun == nil {
		return falseExpr, fmt.Errorf("func %s defined is nil", name)
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

// applyOperator is a dispatcher of the evaluation according to operator
func applyOperator(op Token, l, r Expr) (*BooleanLiteral, error) {
	switch op {
	case AND:
		return applyAND(l, r)
	case OR:
		return applyOR(l, r)
	case EQ:
		return applyEQ(l, r)
	case NEQ:
		return applyNQ(l, r)
	case GT:
		return applyGT(l, r)
	case GTE:
		return applyGTE(l, r)
	case LT:
		return applyLT(l, r)
	case LTE:
		return applyLTE(l, r)
	case XOR:
		return applyXOR(l, r)
	case NAND:
		return applyNAND(l, r)
	case IN:
		return applyIN(l, r)
	case CONTAINS:
		return applyContains(l, r)
	case NOTIN:
		return applyNOTIN(l, r)
	case EREG:
		return applyEREG(l, r)
	case NEREG:
		return applyNEREG(l, r)
	}
	return &BooleanLiteral{Val: false}, fmt.Errorf("Unsupported operator: %s", op)
}
