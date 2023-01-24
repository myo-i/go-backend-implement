package api

import (
	sqlc "go-backend/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  sqlc.Store
	router *gin.Engine
}

// create a new server and setup routing
func NewServer(store sqlc.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// ginが使用している現在のバリデーターエンジンを取得
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	// idに合致したアカウントの情報全てを1度で取得
	router.GET("/accounts/:id", server.getAccount)

	// uriに1を含めただけで全てを取得すると重くなる可能性があるので
	// クエリパラメータで分割して取得することで軽くする狙い
	// そしてクエリ(accounts?id=1)からパラメータを取得するためパスはaccountsとなる
	// ハンドラーの名前はlistAccountでなければならない
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// run HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
