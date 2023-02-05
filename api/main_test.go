package api

import (
	db "go-backend/db/sqlc"
	"go-backend/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

// ginをデバッグモードではなく、テストモードで実行することでログを見やすくしている
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	m.Run()
}
