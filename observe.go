package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var (
	fakeHttpReqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fake_http_request_count",
			Help: "How many HTTP requests processed, partitioned by status code.",
		},
		[]string{"code"},
	)

	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "observe_cockroachdb",
			Name:      "total_rows",
			Help:      "This is my gauge",
		})

	insertHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "observe_cockroachdb",
			Name:      "insert_latency",
			Help:      "This is my histogram",
		})

	selectHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "observe_cockroachdb",
			Name:      "select_latency",
			Help:      "This is my histogram",
		})

	tracer opentracing.Tracer
)

func init() {
	prometheus.MustRegister(fakeHttpReqs)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(insertHistogram)
	prometheus.MustRegister(selectHistogram)
	initTracer()
	rand.Seed(time.Now().UnixNano())
}

func initTracer() {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	traceHost, exists := os.LookupEnv("TRACE_HOST")
	if !exists {
		traceHost = "localhost"
	}
	cfg := jaegercfg.Configuration{
		ServiceName: "observe_cockroachdb",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: fmt.Sprintf("http://%s:14268/api/traces", traceHost),
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	gTracer, _, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(gTracer)
	tracer = opentracing.GlobalTracer()
}

func getEnvDefault(env string, envDefault string) string {
	envVal, exists := os.LookupEnv(env)
	if !exists {
		envVal = envDefault
	}

	return envVal
}

func connect() *pgx.Conn {

	pgUser := getEnvDefault("PG_USER", "root")
	pgHost := getEnvDefault("PG_HOST", "localhost")
	pgDb := getEnvDefault("PG_DB", "defaultdb")

	connString := fmt.Sprintf("postgresql://%s@%s:26257/%s?sslmode=disable", pgUser, pgHost, pgDb)
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	// Connect to the "bank" database.
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	// Create the "accounts" table.
	if _, err := conn.Exec(context.Background(),
		"CREATE TABLE IF NOT EXISTS accounts (id INT PRIMARY KEY, balance INT);CREATE SEQUENCE IF NOT EXISTS accountIds;"); err != nil {
		log.Fatal(err)
	}

	return conn
}

func randomMilliWait() {
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(3000)))
}

func insertAccount(conn *pgx.Conn, parentSpan opentracing.Span) {
	childSpan := tracer.StartSpan(
		"insert",
		opentracing.ChildOf(parentSpan.Context()),
	)

	defer childSpan.Finish()
	start := time.Now()

	insert := fmt.Sprintf("INSERT INTO accounts (id, balance) VALUES (nextval('accountIds'), %d)", rand.Intn(3000))
	// Insert two rows into the "accounts" table.
	if _, err := conn.Exec(context.Background(), insert); err != nil {
		log.Println(err)
		return
	}

	randomMilliWait()

	insertHistogram.Observe(time.Since(start).Seconds())
}

func selectAccounts(conn *pgx.Conn, parentSpan opentracing.Span) int {
	childSpan := tracer.StartSpan(
		"select",
		opentracing.ChildOf(parentSpan.Context()),
	)
	defer childSpan.Finish()

	start := time.Now()
	rows, err := conn.Query(context.Background(), "SELECT id, balance FROM accounts")
	if err != nil {
		log.Println(err)
		return 0
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
		var id, balance int
		if err := rows.Scan(&id, &balance); err != nil {
			log.Fatal(err)
		}
	}

	selectHistogram.Observe(time.Since(start).Seconds())
	gauge.Set(float64(count))

	return count
}

// designed to show how a counter would work with http status codes
// has a 1/4 chance of returing a failing status code
func makeFakeHttpReq() {
	httpStats := []string{"200", "200", "200", "500"}
	randomIndex := rand.Intn(len(httpStats))
	code := httpStats[randomIndex]
	fakeHttpReqs.WithLabelValues(code).Add(1)
}

func main() {
	conn := connect()
	defer conn.Close(context.Background())

	go func() {
		for {
			makeFakeHttpReq()
			span := tracer.StartSpan("insert-select-cockroachdb")
			insertAccount(conn, span)
			count := selectAccounts(conn, span)
			span.Finish()

			log.Println("Row added. Total rows:", count)

			time.Sleep(time.Second * 5)

		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("Starting server on 0.0.0.0:8090")
	http.ListenAndServe(":8090", nil)
}
