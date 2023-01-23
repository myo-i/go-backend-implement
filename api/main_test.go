package api

import (
	"testing"

	"github.com/gin-gonic/gin"
)

// ginをデバッグモードではなく、テストモードで実行することでログを見やすくしている
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	m.Run()
}
