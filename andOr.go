package conditions

// applyNAND applies NAND operation to l/r operands
func applyNAND(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b bool
		err  error
	)
	a, err = getBoolean(l)
	if err != nil {
		return nil, err
	}
	b, err = getBoolean(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (!(a && b))}, nil
}

// applyAND applies && operation to l/r operands
func applyAND(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b bool
		err  error
	)
	a, err = getBoolean(l)
	if err != nil {
		return nil, err
	}
	b, err = getBoolean(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a && b)}, nil
}

// applyOR applies || operation to l/r operands
func applyOR(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b bool
		err  error
	)
	a, err = getBoolean(l)
	if err != nil {
		return nil, err
	}
	b, err = getBoolean(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a || b)}, nil
}

// applyXOR applies || operation to l/r operands
func applyXOR(l, r Expr) (*BooleanLiteral, error) {
	var (
		a, b bool
		err  error
	)
	a, err = getBoolean(l)
	if err != nil {
		return nil, err
	}
	b, err = getBoolean(r)
	if err != nil {
		return nil, err
	}
	return &BooleanLiteral{Val: (a != b)}, nil
}
