# Houyi

Houyi (Chinese: 后羿), based on [Jaeger](https://github.com/jaegertracing/jaeger), 
is a distributed tracing tool implemented by Golang. 

The biggest difference between Houyi and other tracing tools is that Houyi implements
 adaptive sampling for making sampling decision in a completely new way.
 
Houyi contains two component:

- Agent: Receives spans from Houyi client and transfer the requests for pulling sampling strategies from Houyi clients to Houyi collector.
- Collector: Receives spans and requests for pulling sampling strategies from Houyi Agent .

## Features

### Span Filtering

By providing filter conditions, Houyi would increase the sampling rates of spans that meet the filter conditions.
Filtering conditions can be changed during the runtime via RESTFul APIs.

An example of filter conditions (filter span by tags):

```json
{
  "filter-tags": [
    {
      "key": "http.status_code",
      "operation": "!=",
      "value": 200
    }
  ]
}
```

### Adaptive Sampling

A trace with a higher QPS would gain a lower sample rate. 
In contrast, for a trace with a lower QPS, the sampling rate would be higher.

### Dynamical Update of Execution Path

Houyi would store the execution path of all traces, and it would dynamically update the execution path due to 
the unpredictable changes of state (deployment, uninstall, shutdown, etc) of microservices in system.
