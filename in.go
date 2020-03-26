package conditions

import "fmt"

// applyNOTIN applies NOT IN operation to l/r operands
func applyNOTIN(l, r Expr) (*BooleanLiteral, error) {
	result, err := applyIN(l, r)
	if result != nil {
		result.Val = !result.Val
	}
	return result, err
}

// applyContains applies CONTAINS to l/r operations
func applyContains(l, r Expr) (*BooleanLiteral, error) {
	return applyIN(r, l)
}

// applyIN applies IN operation to l/r operands
func applyIN(l, r Expr) (*BooleanLiteral, error) {
	var (
		err   error
		found bool
	)
	// pp.Print(l)
	switch t := l.(type) {
	case *StringLiteral:
		var a string
		var b []string
		a, err = getString(l)
		if err != nil {
			return nil, err
		}

		b, err = getSliceString(r)

		if err != nil {
			return nil, err
		}

		found = false
		for _, e := range b {
			if a == e {
				found = true
			}
		}
	case *NumberLiteral:
		var a float64
		var b []float64
		a, err = getNumber(l)
		if err != nil {
			return nil, err
		}

		b, err = getSliceNumber(r)

		if err != nil {
			return nil, err
		}

		found = false
		for _, e := range b {
			if a == e {
				found = true
			}
		}
	case *NilLiteral:
		found = false
	default:
		return nil, fmt.Errorf("Can not evaluate Literal of unknow type %s %T", t, t)
	}

	return &BooleanLiteral{Val: found}, nil
}
