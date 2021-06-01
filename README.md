# Houyi

Houyi (Chinese: 后羿), based on [Jaeger](https://github.com/jaegertracing/jaeger),
is a distributed tracing tool implemented by Golang 1.15.

The major difference between Houyi and other tracing tools is that Houyi implements
dynamic sampling to make biased sampling decision for different executions (normal or error) in a new way.

## Components

### Client

[Client](https://github.com/houyi-tracing/houyi-client-go) is the instrumentation code which is embedded in the microservice code to record the process and details of the request processing by the microservice. Details related to single request is called **Span**.

A span contains:

- Trace Context:
 - Trace ID
 - Span ID
 - Parent Span ID

- Caller operation name (The name of remote operation that call this local operation)
- Local service name
- Duration
- Tags
- Logs
- ...

### Agent

The Agent is responsible for forwarding spans from the Client to the Collector and forwarding the request from Client to pull sampling strategy to Central Server.

### Collector

- Process spans from agent
- Evaluate tags of span to update sampling strategy
- Store spans to databases (Apache Cassandra, etc.)

### Central Server

- Store sampling strategies
- Process periodical request to pull sampling strategy from Agent (actually from Client)
- Process request to update sampling strategy from Collector.

## How to build

Chang directory to `build`, run:

```shell
./build-all.sh
```

All binaries of components can be found in the directory `$HOME/houyi/`, which contains subdirectories named as components.

## How to run

### 1. Docker

All components in this project has been build as Docker images and all images has been upload to Docker Hub, so that we can locally run all components by pulling images from Docker Hub.

#### Run Agent

```shell
docker run houyitracing/agent
```

#### Run Collector

```
docker run houyitracing/collector
```

#### Run Central Server

```shell
docker run houyitracing/cs
```

### 2. Kubernetes + Istio

We provided `yaml` files for deploying all components in Kubernetes cluster which has deploy service mesh [Istio](https://istio.io/latest/).

#### Deploy collector

```
kubectl apply -f ./build/collector/kube.yaml
kubectl apply -f ./build/collector/istio.yaml
```

#### Deploy Central Server

```
kubectl apply -f ./build/cs/kube.yaml
kubectl apply -f ./build/cs/istio.yaml
```

***Agent** is bundled with actual microservices which contains Client so that it should not be deployed independently.