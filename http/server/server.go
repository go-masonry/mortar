package server

import (
	"context"
	"fmt"
	"github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"sync"
)

type listenerMuxPair struct {
	m mux
	l net.Listener
}

type webService struct {
	sync.Mutex
	serviceConfig   *webServiceConfig
	grpcServer      *grpc.Server
	grpcAddr        string
	muxAndListeners []*listenerMuxPair
	close           bool
}

func newWebService(cfg *webServiceConfig) (instance server.WebService, err error) {
	ws := &webService{
		serviceConfig: cfg,
	}
	if err = ws.setupGRPC(ws.serviceConfig.grpc); err == nil {
		err = ws.setupREST(ws.serviceConfig.rest)
	}
	if err != nil { // make sure to clean, since there might still be open listeners
		for _, pair := range ws.muxAndListeners {
			if err := pair.l.Close(); err != nil {
				ws.serviceConfig.logger(context.Background(), "error closing [%s] listener, %v", pair.l.Addr().String(), err)
			}
		}
	} else {
		instance = ws
	}
	return
}

func (ws *webService) Run(context.Context) (err error) {
	ws.Lock()
	if !ws.close {
		shutdownChannel := make(chan error)
		for _, srvPair := range ws.muxAndListeners {
			go func(pair *listenerMuxPair) {
				shutdownChannel <- pair.m.Serve(pair.l)
				ws.Stop(context.Background())
			}(srvPair)
		}
		ws.Unlock()
		return <-shutdownChannel
	}
	ws.Unlock()
	return fmt.Errorf("server already closed")
}

func (ws *webService) Stop(ctx context.Context) error {
	ws.Lock()
	defer ws.Unlock()
	if ws.close {
		return nil // already closed
	}
	ws.serviceConfig.logger(ctx, "Shutting down...")
	ws.close = true
	var wg sync.WaitGroup
	for _, listenerAndMux := range ws.muxAndListeners {
		switch s := listenerAndMux.m.(type) {
		case grpcServerStopper:
			wg.Add(1)
			go func(stopper grpcServerStopper, listener net.Listener) {
				defer wg.Done()
				stopper.GracefulStop()
				listener.Close()
			}(s, listenerAndMux.l)
		case restServerShutdown:
			wg.Add(1)
			go func(stopper restServerShutdown, listener net.Listener) {
				defer wg.Done()
				stopper.Shutdown(ctx)
				listener.Close()
			}(s, listenerAndMux.l)
		}
	}
	var allClosed = make(chan error)
	go func() {
		wg.Wait()
		allClosed <- nil
	}()
	select {
	case <-ctx.Done():
		ws.serviceConfig.logger(ctx, "Graceful shutdown wasn't finished, %v", ctx.Err())
		return ctx.Err()
	case err := <-allClosed:
		return err
	}
}

func (ws *webService) Ports() (list []server.ListenInfo) {
	grpcPort := extractPort(ws.grpcAddr)
	list = append(list, server.ListenInfo{
		Address: ws.grpcAddr,
		Port:    grpcPort,
		Type:    server.GRPCServer,
	})
	for _, pair := range ws.muxAndListeners {
		port := extractPort(pair.l.Addr().String())
		if port != grpcPort {
			list = append(list, server.ListenInfo{
				Address: pair.l.Addr().String(),
				Port:    port,
				Type:    server.RESTServer,
			})
		}
	}
	return
}

func (ws *webService) setupGRPC(cfg *grpcConfig) (err error) {
	if cfg.registerApi == nil {
		err = fmt.Errorf("no GRPC APIs registered, make sure to call 'RegisterGRPCAPIs' when building")
	} else {
		// Listener
		grpcListener := cfg.listener
		if grpcListener == nil {
			if grpcListener, err = createListener("tcp", cfg.addr); err != nil {
				return err
			}
		}
		// Server
		ws.grpcServer = cfg.server
		if ws.grpcServer == nil {
			ws.grpcServer = grpc.NewServer(cfg.options...)
		}
		cfg.registerApi(ws.grpcServer)
		ws.registerHealthService(ws.grpcServer) // register internal health service
		// save, since this should run first we have no problem with previous values
		ws.muxAndListeners = append(ws.muxAndListeners, &listenerMuxPair{l: grpcListener, m: ws.grpcServer})
		ws.grpcAddr = grpcListener.Addr().String() // we need this later for grpc gateway
	}
	return
}

func (ws *webService) setupREST(restConfigs []*restConfig) (err error) {
	for _, cfg := range restConfigs {
		var emptyListener = true // indicate that we have some kind of handler here, grpcgateway or custom handler/handlerfunc
		webSrv := cfg.server
		// Listener
		restListener := cfg.listener
		if restListener == nil {
			if webSrv != nil && len(webSrv.Addr) > 0 {
				restListener, err = createListener("tcp", webSrv.Addr)
			} else {
				restListener, err = createListener("tcp", cfg.addr)
			}
			if err != nil {
				return err
			}
		}
		// Server
		if webSrv == nil {
			webSrv = &http.Server{Addr: restListener.Addr().String()}
		}
		if webSrv.Handler == nil {
			webSrv.Handler = http.NewServeMux()
		}
		// register handlers
		for pattern, handler := range cfg.handlers {
			if muxHandler, ok := webSrv.Handler.(muxHandler); ok {
				emptyListener = false
				muxHandler.Handle(pattern, handler)
			} else {
				return fmt.Errorf("[%s] handler can't be registered since the provided *http.Server can't handle them", pattern)
			}
		}
		// register handler functions
		for pattern, handlerFunc := range cfg.handlerFuncs {
			if muxHandler, ok := webSrv.Handler.(muxHandler); ok {
				emptyListener = false
				muxHandler.HandleFunc(pattern, handlerFunc)
			} else {
				return fmt.Errorf("[%s] handler function can't be registered since the provided *http.Server can't handle them", pattern)
			}
		}
		// GRPC Gateway
		if len(cfg.grpcGatewayHandlers) > 0 {
			var rootTaken bool // Check if the root '/' pattern is taken
			if _, rootTaken = cfg.handlers["/"]; !rootTaken {
				_, rootTaken = cfg.handlerFuncs["/"]
			}
			if rootTaken {
				return fmt.Errorf("if you want to use GRPC Gateway, you can't take over the root '/' pattern with one of the provided handlers")
			}
			gwMux := cfg.grpcGatewayMux
			if gwMux == nil {
				gwMux = runtime.NewServeMux(cfg.grpcGatewayOptions...)
			}
			// register grpc gateway handlers
			for _, gwHandler := range cfg.grpcGatewayHandlers {
				if err = gwHandler(gwMux, ws.grpcAddr); err != nil {
					return err
				}
				emptyListener = false
			}
			if muxHandler, ok := webSrv.Handler.(muxHandler); ok {
				muxHandler.Handle("/", gwMux)
			} else {
				return fmt.Errorf("grpc Gateway handlers can't be registered, since the provided *http.Server can't handle them")
			}
		}
		// check if we have configured anything
		if emptyListener {
			return fmt.Errorf("nothing to handle for this address: %s", restListener.Addr())
		}
		// Save
		ws.muxAndListeners = append(ws.muxAndListeners, &listenerMuxPair{l: restListener, m: webSrv})
	}
	return
}

func (ws *webService) registerHealthService(server *grpc.Server) {
	ws.serviceConfig.logger(context.Background(), "Registering internal health service")
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthService)
}

// Sanity

var _ server.WebService = (*webService)(nil)
