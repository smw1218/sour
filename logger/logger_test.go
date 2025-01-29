package logger

import "testing"

func TestHandlerName(t *testing.T) {
	cases := map[string]string{
		"github.com/smw1218/sour/cmd/test-service/te.(*TEHandlers).Test-fm":                        "te.Test",
		"github.com/smw1218/sour/cmd/test-service/te.TopLevel":                                     "te.TopLevel",
		"github.com/smw1218/sour/cmd/test-service/app.(*TestService).RegisterRoutes.Wrapped.func1": "app.RegisterRoutes.Wrapped.func1",
	}
	for input, expected := range cases {
		actual := HandlerName(input)
		if actual != expected {
			t.Fatalf("Failed %v, expected %v got %v", input, expected, actual)
		}
	}
}
