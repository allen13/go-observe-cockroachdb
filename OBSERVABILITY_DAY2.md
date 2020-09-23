observability day 2
-------------------

Learning topics:
* Tracing
* Jaeger
* OpenTracing

deploy applications
-------------------

Install the appp

    oc login
    helm upgrade --install cockroachdb helm-charts/cockroachdb -n my-namespace
    helm upgrade --install go-observe-cockroachdb helm-charts/go-observe-cockroachdb -n my-namespace
    helm upgrade --install jaeger helm-charts/jaeger -n my-namespace

tracing big picture
-------------------

In software, observability typically refers to telemetry produced by services and is often divided into three major verticals:

* Tracing, aka distributed tracing, provides insight into the full lifecycles, aka traces, of requests to the system, allowing you to pinpoint failures and performance issues.

* Metrics provide quantitative information about processes running inside the system, including counters, gauges, and histograms.

* Logging provides insight into application-specific messages emitted by processes.
These verticals are tightly interconnected. Metrics can be used to pinpoint, for example, a subset of misbehaving traces. Logs associated with those traces could help to find the root cause of this behavior. And then new metrics can be configured, based on this discovery, to catch this issue earlier next time. Other verticals exist (continuous profiling, production debugging, etc.), however traces, metrics, and logs are the three most well adopted across the industry.

jaeger
------

Jaeger, inspired by Dapper and OpenZipkin, is a distributed tracing system released as open source by Uber Technologies. It is used for monitoring and troubleshooting microservices-based distributed systems, including:

* Distributed context propagation
* Distributed transaction monitoring
* Root cause analysis
* Service dependency analysis
* Performance / latency optimization

* [Architecture](https://www.jaegertracing.io/docs/1.19/architecture/)

jaeger hotrod example app
-------------------------

Walk through the official jaeger guide [Jaeger Hotrod](https://medium.com/opentracing/take-opentracing-for-a-hotrod-ride-f6e3141f7941)


opentracing
-----------

The OpenTracing API provides a standard, vendor neutral framework for instrumentation. This means that if a developer wants to try out a different distributed tracing system, then instead of repeating the whole instrumentation process for the new distributed tracing system, the developer can simply change the configuration of the Tracer.

* [Overview](https://opentracing.io/docs/overview/)

go-observe-cockroachdb
----------------------

Review the included `observe.go` application to see opentracing in action

cockroachdb
-----------

Attempt to enable tracing on cockroachdb. The Zipkin port has already been setup

How to log into the SQL shell

    oc rsh cockroachdb-0
    cockroach sql --insecure

* [cockroachdb tracing](https://wiki.crdb.io/wiki/spaces/CRDB/pages/73171339/Tracing+logs+with+Jaeger+and+Zipkin)





