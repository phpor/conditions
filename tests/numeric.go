package tests

var ValidateTestDataNumeric = []TestData{

	{"56.43 > 56.42", nil, true, false},
	{"-1 > -2", nil, true, false},
	{"-1 >= -2", nil, true, false},
	{"2 >= -3", nil, true, false},
	{"2.1 >= 0.5", nil, true, false},
	{"-3 < 2", nil, true, false},
	{"-3 < 0.4", nil, true, false},
}
