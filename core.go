package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // nolint
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/projecteru2/core/auth"
	"github.com/projecteru2/core/cluster/calcium"
	"github.com/projecteru2/core/engine/factory"
	"github.com/projecteru2/core/log"
	"github.com/projecteru2/core/metrics"
	"github.com/projecteru2/core/rpc"
	pb "github.com/projecteru2/core/rpc/gen"
	"github.com/projecteru2/core/selfmon"
	"github.com/projecteru2/core/utils"
	"github.com/projecteru2/core/version"

	"github.com/getsentry/sentry-go"
	cli "github.com/urfave/cli/v2"
	_ "go.uber.org/automaxprocs"
	"google.golang.org/grpc"
)

var (
	configPath      string
	embeddedStorage bool
)

func serve(c *cli.Context) error {
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("[main] %v", err)
	}

	if err := log.SetupLog(config.LogLevel); err != nil {
		log.Fatalf("[main] %v", err)
	}

	if config.SentryDSN != "" {
		defer utils.SentryDefer()
		log.Infof(c.Context, "[main] sentry %v", config.SentryDSN)
		_ = sentry.Init(sentry.ClientOptions{Dsn: config.SentryDSN})
	}

	// init engine cache and start engine cache checker
	factory.InitEngineCache(c.Context, config)

	var t *testing.T
	if embeddedStorage {
		t = &testing.T{}
	}
	cluster, err := calcium.New(c.Context, config, t)
	if err != nil {
		log.Errorf(c.Context, "[main] %v", err) //nolint
		return err
	}
	defer cluster.Finalizer()
	cluster.DisasterRecover(c.Context)

	stop := make(chan struct{}, 1)
	vibranium := rpc.New(cluster, config, stop)
	s, err := net.Listen("tcp", config.Bind)
	if err != nil {
		log.Errorf(c.Context, "[main] %v", err) //nolint
		return err
	}

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(uint32(config.GRPCConfig.MaxConcurrentStreams)),
		grpc.MaxRecvMsgSize(config.GRPCConfig.MaxRecvMsgSize),
	}

	if config.Auth.Username != "" {
		log.Info("[main] Cluster auth enable.")
		auth := auth.NewAuth(config.Auth)
		opts = append(opts, grpc.StreamInterceptor(auth.StreamInterceptor))
		opts = append(opts, grpc.UnaryInterceptor(auth.UnaryInterceptor))
		log.Infof(c.Context, "[main] Username %s Password %s", config.Auth.Username, config.Auth.Password) //nolint
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCoreRPCServer(grpcServer, vibranium)
	utils.SentryGo(func() {
		if err := grpcServer.Serve(s); err != nil {
			log.Errorf(c.Context, "[main] start grpc failed %v", err) //nolint
		}
	})

	if config.Profile != "" {
		http.Handle("/metrics", metrics.Client.ResourceMiddleware(cluster)(promhttp.Handler()))
		utils.SentryGo(func() {
			if err := http.ListenAndServe(config.Profile, nil); err != nil {
				log.Errorf(c.Context, "[main] start http failed %v", err) //nolint
			}
		})
	}

	unregisterService, err := cluster.RegisterService(c.Context)
	if err != nil {
		log.Errorf(c.Context, "[main] failed to register service: %v", err) //nolint
		return err
	}
	log.Info("[main] Cluster started successfully.")

	// wait for unix signals and try to GracefulStop
	ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	// start node status checker
	go selfmon.RunNodeStatusWatcher(ctx, config, cluster, t)

	<-ctx.Done()

	log.Info("[main] Interrupt by signal")
	close(stop)
	unregisterService()
	grpcServer.GracefulStop()
	log.Info("[main] gRPC server gracefully stopped.")

	log.Info("[main] Check if cluster still have running tasks.")
	vibranium.Wait()
	log.Info("[main] cluster gracefully stopped.")
	return nil
}

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Print(version.String())
	}

	app := cli.NewApp()
	app.Name = version.NAME
	app.Usage = "Run eru core"
	app.Version = version.VERSION
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       "/etc/eru/core.yaml",
			Usage:       "config file path for core, in yaml",
			Destination: &configPath,
			EnvVars:     []string{"ERU_CONFIG_PATH"},
		},
		&cli.BoolFlag{
			Name:        "embedded-storage",
			Usage:       "active embedded storage",
			Destination: &embeddedStorage,
		},
	}
	app.Action = serve
	_ = app.Run(os.Args)
}
