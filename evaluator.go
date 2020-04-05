package conditions

import (
	"fmt"
	"reflect"
)

var (
	falseExpr = &BooleanLiteral{Val: false}
	nilExpr   = &NilLiteral{Val: false}
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
	case *NilLiteral:  // 将nil视为bool来计算，其它表达式计算中也已经大量的将nil视为bool值了
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
			// todo: 这里的lv里面是什么类型并不确定，对于非bool类型，要么强制转成bool类型来求值，要么报类型错误
			// 1. 或者添加一个Bool()接口，如果实现了该接口，则通过该接口获取bool值，如果没有实现该接口，则报类型错误
			// 2. 或者在getBool()中添加对相关数据类型的支持
			b, _ := getBoolean(lv)
			if !b {
				return &BooleanLiteral{Val: false}, nil
			}
		}
		if n.Op == OR {
			b, _ := getBoolean(lv)
			if b {
				return &BooleanLiteral{Val: true}, nil
			}
		}
		rv, err = evaluateSubtree(n.RHS, args)
		if err != nil {
			return falseExpr, err
		}
		return applyOperator(n.Op, lv, rv)
	case *VarRef:
		name := n.Val
		typ := reflect.TypeOf(args)
		if typ == nil { // 这里 nil 可以报错，允许上下文数据为空，使得所有数据都为nil
			return nilExpr, nil
		}
		argsKind := typ.Kind()
		var val interface{}

		// todo: 下面对于值找不到或者为nil的情况，需要慎重处理，否则会带来很多的混乱,是否可以有一个nil Literal，允许参与到多数类型的比较和运算
		// 下面的逻辑和上面的逻辑结合来看，如果报错就不会继续进行后续的求值； 如果希望当做nil来求值，则不应该返回错误；
		// 至于nil能否和其他表达式来求值，应该是对nil的细节处理来决定的
		switch argsKind {
		case reflect.Map:
			argsMap, ok := args.(map[string]interface{})
			if !ok {
				return falseExpr, fmt.Errorf("Args: `%v` convert to map[string] fail", args)
			}
			if _, ok := argsMap[name]; !ok {
				return nilExpr, nil
			}
			val, _ = argsMap[name]
		case reflect.Struct:
			ps := reflect.ValueOf(args)
			fval := ps.FieldByName(name) //todo: 这里有可能panic，需要注意
			if !fval.IsValid() {         // 未找到对应字段，返回nil，尽量避免报错
				return nilExpr, nil
			}
			val = fval.Interface()
		case reflect.Ptr:
			v := reflect.ValueOf(args)
			v = v.Elem()
			if v.Kind() != reflect.Struct {
				return falseExpr, fmt.Errorf("args: `%v` is not map or struct or *struct", args)
			}
			fval := v.FieldByName(name) //todo: 这里有可能panic，需要注意
			if !fval.IsValid() {         // 未找到对应字段，返回nil，尽量避免报错
				return nilExpr, nil
			}
			val = fval.Interface()
		default:
			return falseExpr, fmt.Errorf("Args: `%v` is not map or struct", args)
		}

		typ = reflect.TypeOf(val)
		if typ == nil { // 这里对于找不到的数据如何处理非常关键
			return nilExpr, nil
		}
		kind := typ.Kind()
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
			return evalFunc(val, name)
		}
		return falseExpr, fmt.Errorf("unsupported argument %s type: %s", n.Val, kind)
	}

	return expr, nil
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
