package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listen  string
	metrics int // metrics 数量
	length  int // metrics 长度
	labels  int // 标签数量
	key     int // 标签key长度
	value   int // 标签value长度
	max     int // 指标最大值
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.StringVar(&listen, "listen", ":9001", "listen port")
	flag.IntVar(&metrics, "metrics", 10000, "metric number")
	flag.IntVar(&length, "length", 25, "metric name length")
	flag.IntVar(&labels, "labels", 3, "label number")
	flag.IntVar(&key, "key", 15, "label key length")
	flag.IntVar(&value, "value", 25, "label value length")
	flag.IntVar(&max, "max", 10000, "metric max value")
	flag.Parse()
}

func generate(prefix string, length, index int) string {
	length = length - len(prefix)
	return fmt.Sprintf("%s%0*d", prefix, length, index)
}

func generateLabels(prefix string, length int, index int) (result []string) {
	for i := 0; i < labels; i++ {
		prefix := fmt.Sprintf("%s%d", prefix, i)
		result = append(result, generate(prefix, length, index))
	}
	return
}

func generateMetric() {
	gauges := make([]prometheus.Gauge, metrics)

	for i := 0; i < metrics; i++ {
		gauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: generate("metric", length, i),
			},
			generateLabels("key", key, i),
		).WithLabelValues(generateLabels("value", value, i)...)
		gauges[i] = gauge
	}

	go func(gauges []prometheus.Gauge) {
		for i := range gauges {
			prometheus.MustRegister(gauges[i])
		}

		for {
			for i := range gauges {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				gauges[i].Set(float64(r.Intn(max)))
			}
			time.Sleep(time.Minute)
		}
	}(gauges)
}

func main() {
	log.Printf("listen on: \"%s\"", listen)
	log.Printf("metric number is: %d", metrics)
	log.Printf("metirc length is: %d", length)
	log.Printf("labbel number is: %d", labels)
	log.Printf("labbel key length is: %d", key)
	log.Printf("labbel value length is: %d", value)
	log.Printf("metric mac value is: %d", max)

	go func() {
		log.Println(http.ListenAndServe(":9000", nil))
	}()

	generateMetric()
	http.Handle("/metrics", promhttp.Handler())
	log.Println(http.ListenAndServe(listen, nil))
}
