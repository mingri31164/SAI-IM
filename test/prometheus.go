package test

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func main() {
	/**
	* 采集目标有两种：
	      1. go语言运行的相关数据指标（比如内存的使用、协程数，在下载go语言对prometheus支持的包中已内部完成）
	      2. 自定义采集信息
	**/

	// 自定义监控指标
	temp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tests_temp_gauge",
		Help: "the is test gauge",
	})
	//将采集信息注册到prometheus中
	prometheus.MustRegister(temp)

	var i int
	go func() {
		for {
			i++
			// 每执行2次就添加一次数据指标（记录）
			if i%2 == 0 {
				temp.Inc()
			}
			time.Sleep(time.Second)
		}
	}()

	// 对prometheus提供一个监听的接口，接口路径与配置中的对应
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println(http.ListenAndServe(":1234", nil))
}
