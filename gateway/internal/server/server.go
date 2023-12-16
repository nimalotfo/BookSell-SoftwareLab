package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	ServerConf Config
}

func NewServer(conf Config) *Server {
	router := gin.Default()
	//AddCors(router)
	// logger := logrus.New()
	// router.Use(ErrorHandler(logger))
	router.Use(CORSMiddleware())
	return &Server{
		router,
		conf,
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Control-Allow-Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc{
// 	return func(c *gin.Context) {
//         c.Next()

//         for _, ginErr := range c.Errors {
//             logger.Errorf("%v\n", ginErr.Err)
//         }

// 		c.JSON(-1, )
//     }
// }

// func AddCors(router *gin.Engine) {
// 	router.Use(cors.New(cors.Config{
// 		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
// 		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "access-control-allow-origin"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		//AllowAllOrigins:  true,
// 		AllowOriginFunc: func(origin string) bool {
// 			return true
// 		},
// 		MaxAge: 86400,
// 	}))
// }

func (s *Server) Serve() chan error {
	addr := fmt.Sprintf("%s:%s", s.ServerConf.Host, s.ServerConf.Port)
	errChan := make(chan error)
	go func() {
		errChan <- s.Router.Run(addr)
	}()
	return errChan
}
