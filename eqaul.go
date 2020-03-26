package conditions

import "fmt"

// applyEQ applies == operation to l/r operands
func applyEQ(l, r Expr) (*BooleanLiteral, error) {
	var (
		as, bs string
		an, bn float64
		ab, bb bool
		err    error
	)
	as, err = getString(l)
	if err == nil {
		bs, err = getString(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare string with non-string")
		}
		return &BooleanLiteral{Val: (as == bs)}, nil
	}
	an, err = getNumber(l)
	if err == nil {
		bn, err = getNumber(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare number with non-number")
		}
		return &BooleanLiteral{Val: (an == bn)}, nil
	}
	ab, err = getBoolean(l)
	if err == nil {
		bb, err = getBoolean(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare boolean with non-boolean")
		}
		return &BooleanLiteral{Val: (ab == bb)}, nil
	}
	return falseExpr, nil
}

// applyNQ applies != operation to l/r operands
func applyNQ(l, r Expr) (*BooleanLiteral, error) {
	var (
		as, bs string
		an, bn float64
		ab, bb bool
		err    error
	)
	as, err = getString(l)
	if err == nil {
		bs, err = getString(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare string with non-string")
		}
		return &BooleanLiteral{Val: (as != bs)}, nil
	}
	an, err = getNumber(l)
	if err == nil {
		bn, err = getNumber(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare number with non-number")
		}
		return &BooleanLiteral{Val: (an != bn)}, nil
	}
	ab, err = getBoolean(l)
	if err == nil {
		bb, err = getBoolean(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare boolean with non-boolean")
		}
		return &BooleanLiteral{Val: (ab != bb)}, nil
	}
	return falseExpr, nil
}
