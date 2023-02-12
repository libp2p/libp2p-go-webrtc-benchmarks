# WebRTC Transport Benchmarks

This directory contains a benchmarking tool and instructions how to use it,
to measure the performance of the WebRTC transport.

- [1. Instructions](#1-instructions)
  - [1.1. Listener](#11-listener)
    - [1.1.1. Metrics](#111-metrics)
  - [1.2. Client](#12-client)
  - [1.3. Profile](#13-profile)
- [2. Benchmarks](#2-benchmarks)
    - [2.1. Scenario 1](#21-scenario-1)
        - [2.1.1. Results](#211-results)
    - [2.2. Scenario 2](#22-scenario-2)
        - [2.2.1. Results](#221-results)

## 1. Instructions

In this section we'll show you how to run this benchmarking tool on your local (development) machine.

1. Run a listener
2. Run a client

What you do next to this depends on what you're after.

- Are you using it to get metrics from a standard and well defined cloud run?
- Are you using it to get metrics from your local machine?
- Are you using it to (Go) profile one or multiple things?

With that in mind, we'll show you how to do all of the above.

### 1.1. Listener

Run:

```
go run ./benchmark/transports/webrtc -metrics csv listen
```

This should output a multiaddr which can be used by the client to connect.
Other transport values supported instead of `webrtc` are: `tcp`, `quic`, `websocket` and `webtransport`.

The listener will continue to run until you kill it.

#### 1.1.1. Metrics

The metrics can be summarized using the `report` command:

```
go run ./benchmark/transports/webrtc report -s 16 metrics_listen_webrtc_c2_s8_e1_p0.csv
```

Which will print the result to the stdout of your terminal.
Or you can visualize them using the bundled python script:

```
./benchmark/transports/webrtc/scripts/visualise/visualise.py metrics_listen_webrtc_c2_s8_e1_p0.csv -s 16
```

Which will open a new window with your graph in it.

More useful is however to save it to a file so we can share it. For the WebRTC results of Scenario 1
we might for example use the following command:

```
 ./benchmark/transports/webrtc/scripts/visualise/visualise.py \
    -s 10000 \
    -o ./benchmark/transports/webrtc/images/s1_webrtc.png \
    ./benchmark/transports/webrtc/results/metrics_dial_webrtc_c10_s100_p0.csv \
    ./benchmark/transports/webrtc/results/metrics_listen_webrtc_e1_p0.csv
```

### 1.2. Client

Run:

```
go run ./benchmark/transports/webrtc -c 2 -s 8 dial <multiaddr>
```

You can configure the number of streams and connections opened by the dialer using opt-in flags.

The client will continue to run until you kill it.

> Tip:
> 
> similar to the `listen` command you can also use the `-metrics <path>.csv` flag to output the metrics to a file.

### 1.3. Profile

Profiling the benchmark tool is supported using the Golang std pprof tool.

E.g. you can start your listener (or client) with the `-profile 6060` flag to enable profiling over http.

With your listener/client running you can then profile using te std golang tool, e.g.:

```
# get cpu profile
go tool pprof http://localhost:6060/debug/pprof/profile

# get memory (heap) profile
go tool pprof http://localhost:6060/debug/pprof/heap

# check contended mutexes
go tool pprof http://localhost:6060/debug/pprof/mutex

# check why threads block
go tool pprof http://localhost:6060/debug/pprof/block

# check the amount of created goroutines
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

It will open an interactive window allowing you to inspect the heap/cpu profile, e.g. to see te top offenders
of your own code by focussing on the relevant module (e.g. `top github.com/libp2p/go-libp2p/p2p/transport/webrtc`).

And of course you can also use the `-pdf` flag to output it to a file instead that you can view in your browser or
any other capable pdf viewer.

## 2. Benchmarks

The goal of this tooling was to be able to benchmark how the WebRTC transport performs on its own
as well as compared to other transports such as QUIC and WebTransport. Not all scenarios which are benchmarked
are compatible with the different transports, but WebRTC is tested on all benchmarked scenarios.

The scenarios described below and the results you'll find at the end are ran on / come from two c5 large EC2 instances.
Each instance has 8 vCPUs and 16GB RAM. More information can be found at:
https://aws.amazon.com/ec2/instance-types/c5/

Dream goal for WebRTC in terms of performance is to consume 2x or less resources compared to quic. For [Scenario 2](#22-scenario-2) the results are currently as follows when comparing WebRTC to quic:

![Scenario 2 — WebRTC and Quic — CPU](./images/s2_webrtc_quic_cpu.png)

![Scenario 2 — WebRTC and Quic — Memory](./images/s2_webrtc_quic_mem.png)


**Scenario 1:**

1. Server, on EC2 instance A, listens on a generated multi address.
2. Client, on EC2 instance B, dials 10 connections, with 1000 streams per connection to the server.

**Scenario 2:**

1. Server, on EC2 instance A, listens on a generated multi address.
2. Client, on EC2 instance B, dials 100 connections, with 100 streams per connection to the server.

For both scenarios the following holds true:

- Connections are ramped up at the rate of 1 connection/sec. 
- Streams are created at the rate of 10 streams/sec.
- This is done to ensure the webrtc transport's inflight request limiting does not start rejecting connections.
- The client opens streams to the server and runs the echo protocol writing 2KiB/s per stream (1 KiB every 500ms).
- We let the tests run for about 5 minute each.

The instances are running each scenario variation one by one, as such there at any given moment only one benchmark script running.

### 2.1. Scenario 1

Server:

```
go run ./benchmark/transports/webrtc/scripts/multirunner listen
```

Client:

```
go run ./benchmark/transports/webrtc/scripts/multirunner dial
```

#### 2.1.1. Results

**All transports in function of CPU and Memory**

![Scenario 1 — All CPU](./images/s1_all_cpu.png)

![Scenario 1 — All Memory](./images/s1_all_mem.png)

**TCP**

![Scenario 1 — TCP](./images/s1_tcp.png)

|                          | s1_tcp_dial.csv | s1_tcp_listen.csv |
|----------------------|-----------------|-------------------|
|              **CPU (%)** |                 |                   |
|                      min |                1|                  0|
|                      max |                6|                  5|
|                      avg |                  1|                1|
|    **Memory Heap (MiB)** |                 |                   |
|                      min |           22.621|             23.225|
|                      max |           97.679|            143.525|
|                      avg |           67.627|             91.435|
|     **Bytes Read (KiB)** |                 |                   |
|                      min |           11.000|              0.000|
|                      max |         2611.000|           2591.000|
|                      avg |           2484.010|         2495.871|
|  **Bytes Written (KiB)** |                 |                   |
|                      min |           11.000|              0.000|
|                      max |         2690.000|           2591.000|
|                      avg |         2525.182|           2484.010|

**WebSocket (WS)**

![Scenario 1 — WebSocket](./images/s1_websocket.png)

|                          | s1_websocket_listen.csv | s1_websocket_dial.csv |
|----------------------|-----------------------|-------------------------|
|              **CPU (%)** |                       |                         |
|                      min |                      1|                        0|
|                      max |                      5|                        4|
|                      avg |                      1|                        1|
|    **Memory Heap (MiB)** |                         |                       |
|                      min |                 26.053|                   23.293|
|                      max |                  143.693|                 98.811|
|                      avg |                   91.029|                 67.741|
|     **Bytes Read (KiB)** |                       |                         |
|                      min |                 11.000|                    0.000|
|                      max |               2598.000|                 2590.000|
|                      avg |               2493.006|                 2485.318|
|  **Bytes Written (KiB)** |                       |                         |
|                      min |                 11.000|                    0.000|
|                      max |               2690.000|                 2590.000|
|                      avg |               2522.309|                 2485.318|


**WebRTC**

![Scenario 1 — WebRTC](./images/s1_webrtc.png)

|                          | s1_webrtc_listen.csv | s1_webrtc_dial.csv |
|----------------------|--------------------|----------------------|
|              **CPU (%)** |                    |                      |
|                      min |                   0|                     0|
|                      max |                  12|                    11|
|                      avg |                   5|                     6|
|    **Memory Heap (MiB)** |                    |                      |
|                      min |              24.665|                22.917|
|                      max |             578.218|               538.788|
|                      avg |             359.889|               329.959|
|     **Bytes Read (KiB)** |                    |                      |
|                      min |               0.000|                 0.000|
|                      max |            2480.000|              2480.000|
|                      avg |            2161.364|              2378.812|
|  **Bytes Written (KiB)** |                    |                      |
|                      min |               0.000|                 0.000|
|                      max |            2591.000|              2480.000|
|                      avg |              2378.812|            2191.760|

### 2.2. Scenario 2

Server:

```
go run ./benchmark/transports/webrtc/scripts/multirunner listen
```

Client:

```
go run ./benchmark/transports/webrtc/scripts/multirunner -s 1 dial
```

#### 2.2.1. Results

**All transports in function of CPU and Memory**

![Scenario 2 — All CPU](./images/s2_all_cpu.png)

![Scenario 2 — All Memory](./images/s2_all_mem.png)

**TCP**

![Scenario 2 — TCP](./images/s2_tcp.png)

|                          | s2_tcp_dial.csv | s2_tcp_listen.csv |
|----------------------|-----------------|-------------------|
|              **CPU (%)** |                 |                   |
|                      min |                1|                  0|
|                      max |                7|                  5|
|                      avg |                  1|                2|
|    **Memory Heap (MiB)** |                 |                   |
|                      min |           23.491|             23.212|
|                      max |          117.364|            104.624|
|                      avg |             72.506|           79.907|
|     **Bytes Read (KiB)** |                 |                   |
|                      min |              0.000|            9.000|
|                      max |         2590.000|           2590.000|
|                      avg |         2492.597|           2484.886|
|  **Bytes Written (KiB)** |                 |                   |
|                      min |           10.000|              0.000|
|                      max |         2692.000|           2590.000|
|                      avg |         2505.835|           2484.886|

**WebSocket (WS)**

![Scenario 2 — WebSocket](./images/s2_websocket.png)

|                          | s2_websocket_dial.csv | s2_websocket_listen.csv |
|----------------------|-----------------------|-------------------------|
|              **CPU (%)** |                       |                         |
|                      min |                      1|                        0|
|                      max |                      7|                        8|
|                      avg |                      2|                        1|
|    **Memory Heap (MiB)** |                       |                         |
|                      min |                 26.143|                   23.008|
|                      max |                129.231|                  148.597|
|                      avg |                 87.724|                   92.364|
|     **Bytes Read (KiB)** |                         |                       |
|                      min |                 13.000|                    0.000|
|                      max |               2609.000|                 2580.000|
|                      avg |               2486.373|                 2474.365|
|  **Bytes Written (KiB)** |                       |                         |
|                      min |                 13.000|                    0.000|
|                      max |               2700.000|                 2580.000|
|                      avg |               2515.700|                 2474.365|

**WebRTC**

![Scenario 2 — WebRTC](./images/s2_webrtc.png)

|                          | s2_webrtc_dial.csv | s2_webrtc_listen.csv |
|----------------------|--------------------|----------------------|
|              **CPU (%)** |                    |                      |
|                      min |                   0|                     0|
|                      max |                  10|                     7|
|                      avg |                   5|                     4|
|    **Memory Heap (MiB)** |                    |                      |
|                      min |              23.876|                26.001|
|                      max |               328.152|             387.359|
|                      avg |             165.073|               211.501|
|     **Bytes Read (KiB)** |                      |                    |
|                      min |               0.000|                 0.000|
|                      max |            2404.000|              2411.000|
|                      avg |            2099.325|              2311.262|
|  **Bytes Written (KiB)** |                    |                      |
|                      min |               0.000|                 0.000|
|                      max |            2506.000|              2411.000|
|                      avg |            2111.179|              2311.262|

**QUIC**

![Scenario 2 — QUIC](./images/s2_quic.png)

|                          | s2_quic_dial.csv | s2_quic_listen.csv |
|----------------------|------------------|--------------------|
|              **CPU (%)** |                  |                    |
|                      min |                 0|                   0|
|                      max |                12|                   7|
|                      avg |                   1|                 2|
|    **Memory Heap (MiB)** |                  |                    |
|                      min |            24.562|              23.292|
|                      max |           184.649|             174.270|
|                      avg |            96.161|              80.405|
|     **Bytes Read (KiB)** |                  |                    |
|                      min |             0.000|               0.000|
|                      max |          2590.000|            2592.000|
|                      avg |          2141.184|            2486.244|
|  **Bytes Written (KiB)** |                    |                  |
|                      min |             0.000|               0.000|
|                      max |          2700.000|            2592.000|
|                      avg |          2154.440|            2486.244|

**WebTransport**

![Scenario 2 — WebTransport](./images/s2_webtransport.png)

|                          | s2_webtransport_dial.csv | s2_webtransport_listen.csv |
|----------------------|--------------------------|----------------------------|
|              **CPU (%)** |                          |                            |
|                      min |                         0|                           0|
|                      max |                         8|                           4|
|                      avg |                         3|                           1|
|    **Memory Heap (MiB)** |                          |                            |
|                      min |                    22.801|                      24.597|
|                      max |                   156.746|                     162.736|
|                      avg |                   103.668|                     102.261|
|     **Bytes Read (KiB)** |                          |                            |
|                      min |                     0.000|                       0.000|
|                      max |                  2590.000|                    2598.000|
|                      avg |                    2486.258|                  2452.049|
|  **Bytes Written (KiB)** |                          |                            |
|                      min |                     0.000|                       0.000|
|                      max |                    2598.000|                  2693.000|
|                      avg |                  2485.053|                    2486.258|
