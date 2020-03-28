package tests

var ValidateTestDataRegex = []TestData{

	// =~
	{"[status] =~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, true, false},
	{"[status] =~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, false, false},

	// !~
	{"[status] !~ /^5\\d\\d/", map[string]interface{}{"status": "500"}, false, false},
	{"[status] !~ /^4\\d\\d/", map[string]interface{}{"status": "500"}, true, false},
}
