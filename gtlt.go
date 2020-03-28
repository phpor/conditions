package conditions

import "fmt"

// applyGT applies > operation to l/r operands
func applyGT(l, r Expr) (*BooleanLiteral, error) {
	return applyGtLt(l, r, GT)
}

// applyGT applies >= operation to l/r operands
func applyGTE(l, r Expr) (*BooleanLiteral, error) {
	return applyGtLt(l, r, GTE)
}

// applyGT applies < operation to l/r operands
func applyLT(l, r Expr) (*BooleanLiteral, error) {
	return applyGtLt(l, r, LT)
}

// applyGT applies <= operation to l/r operands
func applyLTE(l, r Expr) (*BooleanLiteral, error) {
	return applyGtLt(l, r, LTE)
}

// applyGtLt applies > < >= <= for string or number
func applyGtLt(l, r Expr, tok Token) (*BooleanLiteral, error) {
	var (
		as, bs string
		an, bn float64
		err    error
	)
	// 如果第一个值是nil的话，类型取决于第二个参数
	if _, ok := l.(*NilLiteral); ok {
		// 交换左值和右值， 同时使用相反的操作符
		l, r = r, l
		switch tok {
		case GT:
			tok = LT
		case LT:
			tok = GT
		case GTE:
			tok = LTE
		case LTE:
			tok = GTE
		}
	}
	as, err = getString(l)
	if err == nil {
		bs, err = getString(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare string with non-string")
		}
		switch tok {
		case GT:
			return &BooleanLiteral{Val: as > bs}, nil
		case GTE:
			return &BooleanLiteral{Val: as >= bs}, nil
		case LT:
			return &BooleanLiteral{Val: as < bs}, nil
		case LTE:
			return &BooleanLiteral{Val: as <= bs}, nil
		}
		return falseExpr, fmt.Errorf("Unsupport operate [ %s ] for string", tokens[tok])
	}
	an, err = getNumber(l)
	if err == nil {
		bn, err = getNumber(r)
		if err != nil {
			return falseExpr, fmt.Errorf("Cannot compare number with non-number")
		}
		switch tok {
		case GT:
			return &BooleanLiteral{Val: an > bn}, nil
		case GTE:
			return &BooleanLiteral{Val: an >= bn}, nil
		case LT:
			return &BooleanLiteral{Val: an < bn}, nil
		case LTE:
			return &BooleanLiteral{Val: an <= bn}, nil
		}
		return falseExpr, fmt.Errorf("Unsupport operate [ %s ] for number", tokens[tok])
	}
	return falseExpr, fmt.Errorf("Unsupport type for operate [ %s ]", tokens[tok])
}
