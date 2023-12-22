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
	flag.StringVar(&listen, "listen", "127.0.0.1:9001", "listen port")
	flag.IntVar(&metrics, "metrics", 10000, "metric number")
	flag.IntVar(&length, "length", 25, "metric name length")
	flag.IntVar(&labels, "labels", 3, "label number")
	flag.IntVar(&key, "key", 15, "label key length")
	flag.IntVar(&value, "value", 25, "label value length")
	flag.IntVar(&max, "max", 1e8, "metric max value")
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
	for i := 0; i < metrics; i++ {
		g := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: generate("metric", length, i),
			},
			generateLabels("key", key, i),
		).WithLabelValues(generateLabels("value", value, i)...)
		go func(g prometheus.Gauge) {
			prometheus.MustRegister(g)
			for {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				g.Set(float64(r.Intn(value)))
				time.Sleep(time.Minute)
			}
		}(g)
	}
}

func main() {
	log.Printf("listen on: \"%s\"", listen)
	log.Printf("metric number is: %d", metrics)
	log.Printf("metirc length is: %d", length)
	log.Printf("labbel number is: %d", labels)
	log.Printf("labbel key length is: %d", key)
	log.Printf("labbel value length is: %d", value)
	log.Printf("metric mac value is: %d", max)
	log.Printf("please check endpoint \"http://%s/metrics\"", listen)

	go func() {
		log.Println(http.ListenAndServe(":9000", nil))
	}()

	generateMetric()
	http.Handle("/metrics", promhttp.Handler())
	log.Println(http.ListenAndServe(listen, nil))
}
