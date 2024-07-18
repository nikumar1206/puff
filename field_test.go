package puff

import "testing"

type test1 struct {
	Name string
}

type test3 struct {
	ID   []bool `kind:"path"`
	Name string `kind:"query"`
}
type test4 struct {
	ID   int    `kind:"path"`
	Name string `kind:"query" required:"wef"`
}
type test5 struct {
	ID   int    `kind:"path"`
	Name string `kind:"query" required:"false"`
}

func should_be_nil(got any) bool {
	return got == nil
}
func TestHandleInputSchema(t *testing.T) {
	result_test1 := handleInputSchema(3)          //should not pass
	result_test2 := handleInputSchema(new(test1)) // should not pass
	result_test3 := handleInputSchema(new(test3)) //should not pass
	result_test4 := handleInputSchema(new(test4)) //should not pass
	result_test5 := handleInputSchema(new(test5)) // should passed
	if result_test1 == nil || result_test2 == nil || result_test3 == nil || result_test4 == nil || result_test5 != nil {
		t.Fail()
	}
	t.Logf(
		"Result 1: %t\nResult 2: %t\nResult 3: %t\nResult 4: %t\nResult 5: %t",
		!should_be_nil(result_test1),
		!should_be_nil(result_test2),
		!should_be_nil(result_test3),
		!should_be_nil(result_test4),
		should_be_nil(result_test5),
	)
}
