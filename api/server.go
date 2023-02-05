package api

import (
	"fmt"
	sqlc "go-backend/db/sqlc"
	"go-backend/token"
	"go-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      sqlc.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// create a new server and setup routing
func NewServer(config util.Config, store sqlc.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// ginが使用している現在のバリデーターエンジンを取得
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

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
}

// run HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
