package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"

	v1 "hd/api/helloworld/v1"
	"hd/internal/conf"
	"hd/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService,
	summary *service.StuSummaryService, student *service.StudentService, logger log.Logger) *http.Server {

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.Filter(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)))

	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)

	srv.HandleFunc("/task", student.GetTask)
	srv.HandleFunc("/submit-task", student.SubmitTask)

	srv.HandleFunc("/login", student.Login)

	srv.HandleFunc("/list-signs", summary.ListSigns)
	srv.HandleFunc("/list-tasks", summary.ListTaskSummary)

	router := srv.Route("/")
	router.GET("/sign", student.SignNew)
	router.GET("/pre/{task-id}/{name}", student.Pre)
	router.GET("/index", student.Index)

	return srv
}
