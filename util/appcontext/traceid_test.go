package appcontext_test

import (
	"net/http/httptest"
	"testing"

	"github.com/AI1411/golang-admin-api/util/appcontext"

	"github.com/gin-gonic/gin"
)

func newContext() *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	return c
}

func TestSetTraceIDIntoContext(t *testing.T) {
	t.Parallel()

	testTraceID := "test traceid dayo"

	tests := []struct {
		name    string
		context *gin.Context
		traceID string
	}{
		{
			name:    "TraceIDがContextに設定されること",
			context: newContext(),
			traceID: testTraceID,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			appcontext.SetTraceIDIntoContext(tt.context, tt.traceID)

			traceID, exist := tt.context.Get("api-trace-id")
			if !exist {
				t.Errorf("want = %v, got = %v", true, exist)
			} else if testTraceID != traceID.(string) {
				t.Errorf("want = %v, got = %v", testTraceID, traceID.(string))
			}
		})
	}
}

func TestGetTraceIDFromContext(t *testing.T) {
	t.Parallel()
	t.Run("ContextにTraceIDがある場合に取得できること", func(t *testing.T) {
		t.Parallel()
		traceID := "qwertyuiopasdfhjklzxcvbnm"
		con := newContext()
		con.Set("api-trace-id", traceID)
		want := traceID
		if got := appcontext.GetTraceID(con); got != want {
			t.Errorf("want= %v, got = %v", want, got)
		}
	})

	t.Run("ContextにTraceIDがない場合にランダムの文字列が取得できること", func(t *testing.T) {
		t.Parallel()
		con := newContext()
		if got := appcontext.GetTraceID(con); got == "" {
			t.Errorf("want= <random string>, got = \"\"")
		}
	})
}
