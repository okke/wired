package wtemplate_test

import (
	"testing"

	"github.com/okke/wired/wtemplate"
)

type evalContext struct {
	variables map[string]string
}

func (evalContext *evalContext) Solve(name string) string {
	return evalContext.variables[name]
}

func newEvalContext() wtemplate.Context {
	variables := make(map[string]string, 0)
	variables["pepper"] = "jalapeno"
	return &evalContext{variables: variables}
}

func testParser(t *testing.T, template string, expected string) {
	ctx := newEvalContext()

	if txt := wtemplate.Parse(ctx, template); txt != expected {
		t.Errorf("expected %s, got '%s'", expected, txt)
	}

}
func TestParser(t *testing.T) {

	testParser(t, "", "")
	testParser(t, "chipotle", "chipotle")
	testParser(t, "${pepper}", "jalapeno")
	testParser(t, "${pepper}Pepper", "jalapenoPepper")
	testParser(t, "LoveToEat${pepper}", "LoveToEatjalapeno")
	testParser(t, "$pepper", "jalapeno")
	testParser(t, "$pepper/$pepper", "jalapeno/jalapeno")
	testParser(t, "$pepper $pepper", "jalapeno jalapeno")
	testParser(t, "$pepper:$pepper", "jalapeno:jalapeno")
	testParser(t, "$pepper\t$pepper", "jalapeno\tjalapeno")
	testParser(t, "$pepper\n$pepper", "jalapeno\njalapeno")
	testParser(t, "${pepper}\\$pepper", "jalapeno\\jalapeno")
	testParser(t, "$", "")
	testParser(t, "$$", "")
	testParser(t, "${url:http://hotpeppers.com/habanero?color=red}", "http://hotpeppers.com/habanero?color=red")

}
