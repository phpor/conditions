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
	var (
		err error
		in  bool
	)
	switch t := r.(type) {
	case *StringLiteral:
		var a string
		var b []string
		a, err = getString(r)
		if err != nil {
			return nil, err
		}

		b, err = getSliceString(l)

		if err != nil {
			return nil, err
		}

		in = false
		for _, e := range b {
			if a == e {
				in = true
			}
		}
	case *NumberLiteral:
		var a float64
		var b []float64
		a, err = getNumber(r)
		if err != nil {
			return nil, err
		}

		b, err = getSliceNumber(l)

		if err != nil {
			return nil, err
		}

		in = false
		for _, e := range b {
			if a == e {
				in = true
			}
		}
	default:
		return nil, fmt.Errorf("Can not evaluate Literal of unknow type %s %T", t, t)
	}

	return &BooleanLiteral{Val: in}, nil
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
	default:
		return nil, fmt.Errorf("Can not evaluate Literal of unknow type %s %T", t, t)
	}

	return &BooleanLiteral{Val: found}, nil
}
