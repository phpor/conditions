package conditions

import "regexp"

// applyEREG applies EREG operation to l/r operands
func applyNEREG(l, r Expr) (*BooleanLiteral, error) {
	result, err := applyEREG(l, r)
	if result != nil {
		result.Val = !result.Val
	}
	return result, err
}

// applyEREG applies EREG operation to l/r operands
func applyEREG(l, r Expr) (*BooleanLiteral, error) {
	var (
		a     string
		b     string
		err   error
		match bool
	)
	a, err = getString(l)
	if err != nil {
		return nil, err
	}

	b, err = getString(r)
	if err != nil {
		return nil, err
	}
	match = false
	match, err = regexp.MatchString(b, a)

	// pp.Print(a, b, match)
	return &BooleanLiteral{Val: match}, err
}
