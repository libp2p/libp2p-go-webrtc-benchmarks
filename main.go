package main

import (
	"context"
	"flag"
	"time"

	"github.com/little-bear-labs/libp2p-go-webrtc-benchmarks/benchrunner"
)

func main() {
	var cfg benchrunner.RunnerConfig

	// flags used only for listen cmd
	flag.IntVar(&cfg.ListenPort, "l", 9999, "port to listen to, used for listen cmd")
	flag.BoolVar(&cfg.ListenInsecure, "insecure", false, "use an unencrypted connection, used for listen cmd")
	flag.Int64Var(&cfg.ListenSeed, "seed", 0, "set random seed for id generation, used for listen cmd")

	// flags used for both cmds
	flag.StringVar(&cfg.Transport, "t", "webrtc", "use quic instead of webrtc")
	flag.IntVar(&cfg.ProfilePort, "profile", 0, "enable Golang pprof over http on the given port (disabled by default)")
	flag.DurationVar(&cfg.MetricInterval, "interval", time.Second, "interval at which to track/trace a metric point")
	flag.StringVar(&cfg.MetricOutput, "metrics", "", "wrote metrics to CSV or use 'stdout' for stdout")

	// used for dial cmd only
	flag.IntVar(&cfg.DialConnections, "c", 1, "total connections to open")

	// used for dial and report cmd only
	flag.IntVar(&cfg.DialStreams, "s", 1, "set number of streams")

	// parse all flags
	flag.Parse()
	if err := benchrunner.Run(context.Background(), cfg, flag.Args()...); err != nil {
		panic(err)
	}
}
