package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

// RED Pattern
// Request Rate
// Error rate
// Duration

var (
	dbRequestsDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "Db_request_duration_seconds",
		Help: "The duration of the requests to the DB service.",
	})

	dbRequestsCurrent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Db_requests_current",
		Help: "The current number of requests to the DB service.",
	})

	dbClientErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "Db_errors",
		Help: "The total number of DB client errors",
	})
)

func init() {
	prometheus.MustRegister(dbRequestsDuration, dbRequestsCurrent, dbClientErrors)
}

func metricsHandler(w http.ResponseWriter, r *http.Request, s string) {
	// TODO expose metrics here
	promhttp.Handler().ServeHTTP(w, r)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			http.NotFound(w, r)
			return
		}
		fn(w, r, "")
	}
}

func main() {
	password := os.Getenv("PGPASSWORD")
	host := os.Getenv("PGHOST")

	// setup DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to DB!")

	// continously send db queries
	go func(db *sql.DB) {
		for {
			now := time.Now()
			dbRequestsCurrent.Inc()
			log.Println("Sending sql query")
			sqlStatement := "select pg_sleep(floor(random() * 10 + 1)::int);"
			_, err = db.Exec(sqlStatement)
			dbRequestsDuration.Observe(time.Since(now).Seconds())
			if err != nil {
				dbClientErrors.Inc()
				panic(err)
			}
		}
	}(db)

	// setup handler
	http.HandleFunc("/metrics/", makeHandler(metricsHandler))
	log.Println("Port :8000 is ready for metrics collection")

	log.Fatal(http.ListenAndServe(":8000", nil))
}
