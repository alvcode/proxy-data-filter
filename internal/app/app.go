package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"proxy-data-filter/internal/config"
	"proxy-data-filter/internal/controller"
	"proxy-data-filter/internal/logging"
	"strings"
)

type App struct {
	cfg        *config.Config
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(cfg *config.Config) (App, error) {
	router := httprouter.New()

	return App{
		cfg:    cfg,
		router: router,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx2 := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(ctx2)
	})

	logging.GetLogger(ctx).Info("Application initialized and started")
	return grp.Wait()
}

func (a *App) startHTTP(ctx context.Context) error {
	controllerInit := controller.New(a.cfg, a.router)
	errRoute := controllerInit.SetRoutes(ctx)
	if errRoute != nil {
		logging.GetLogger(ctx).WithError(errRoute).Fatal("failed to init routes")
	}

	//a.router.HandlerFunc(http.MethodGet, "/online/companies/:companyHash", NewMainHandler(5).Handle)

	logging.GetLogger(ctx).Printf("IP: %s, Port: %d", a.cfg.HTTP.Host, a.cfg.HTTP.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.Host, a.cfg.HTTP.Port))
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatal("failed to create http listener")
	}

	logging.GetLogger(ctx).Printf("CORS: %+v", a.cfg.Cors)

	c := cors.New(cors.Options{
		AllowedMethods:     strings.Split(a.cfg.Cors.AllowedMethods, ","),
		AllowedOrigins:     strings.Split(a.cfg.Cors.AllowedOrigins, ","),
		AllowedHeaders:     strings.Split(a.cfg.Cors.AllowedHeaders, ","),
		AllowCredentials:   a.cfg.Cors.AllowCredentials,
		OptionsPassthrough: a.cfg.Cors.OptionsPassthrough,
		ExposedHeaders:     strings.Split(a.cfg.Cors.ExposedHeaders, ","),
		Debug:              a.cfg.Cors.Debug,
	})

	hdl := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      hdl,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
	}

	logging.GetLogger(ctx).Println("http server started")
	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logging.GetLogger(ctx).Warningln("Server shutdown")
		default:
			logging.GetLogger(ctx).Fatalln(err)
		}
	}
	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}
	return err
}

//type MainHandler struct {
//	ruleID int
//}
//
//func NewMainHandler(ruleID int) *MainHandler {
//	return &MainHandler{
//		ruleID: ruleID,
//	}
//}
//
//type Resp struct {
//	RuleID int `json:"rule_id"`
//}
//
//func (h *MainHandler) Handle(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(200)
//	resp := Resp{h.ruleID}
//	if err := json.NewEncoder(w).Encode(resp); err != nil {
//		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
//	}
//}
