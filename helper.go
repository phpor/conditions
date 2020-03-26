package conditions

import "fmt"

// getBoolean performs type assertion and returns boolean value or error
func getBoolean(e Expr) (bool, error) { // 尽量让其它类型也能作为bool值处理
	switch n := e.(type) {
	case *BooleanLiteral:
		return n.Val, nil
	case *NilLiteral:
		return n.Val, nil
	case *StringLiteral:
		return len(n.Val) > 0, nil
	case *NumberLiteral:
		return n.Val != 0, nil
	case *SliceStringLiteral:
		return len(n.Val) > 0, nil
	case *SliceNumberLiteral:
		return len(n.Val) > 0, nil
	default:
		return false, fmt.Errorf("Literal is not a boolean: %v", n)
	}
}

// getString performs type assertion and returns string value or error
func getString(e Expr) (string, error) {
	switch n := e.(type) {
	case *StringLiteral:
		return n.Val, nil
	case *NilLiteral:
		return "", nil
	default:
		return "", fmt.Errorf("Literal is not a string: %v", n)
	}
}

// getSliceNumber performs type assertion and returns []float64 value or error
func getSliceNumber(e Expr) ([]float64, error) {
	switch n := e.(type) {
	case *SliceNumberLiteral:
		return n.Val, nil
	case *NilLiteral:
		return []float64{}, nil
	default:
		return []float64{}, fmt.Errorf("Literal is not a slice of float64: %v", n)
	}
}

// getSliceString performs type assertion and returns []string value or error
func getSliceString(e Expr) ([]string, error) {
	switch n := e.(type) {
	case *SliceStringLiteral:
		return n.Val, nil
	case *NilLiteral:
		return []string{}, nil
	default:
		return []string{}, fmt.Errorf("Literal is not a slice of string: %v", n)
	}
}

// getNumber performs type assertion and returns float64 value or error
func getNumber(e Expr) (float64, error) {
	switch n := e.(type) {
	case *NumberLiteral:
		return n.Val, nil
	case *NilLiteral:
		return 0, nil
	default:
		return 0, fmt.Errorf("Literal is not a number: %v", n)
	}
}
