package test

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"strings"
	"testing"
)

func Test_Jaeger(t *testing.T) {
	cfg := jaegercfg.Configuration{
		// 定义取样器，即要收集的信息
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// 信息发送的对象，这里为Jaeger的服务器对象
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			// 服务器地址
			//🔥 我们当前的请求，是基于API的方式来进行接收的，所以这里的类型是api/traces
			CollectorEndpoint: fmt.Sprintf("http://%s/api/traces", "118.178.120.11:14268"),
		},
	}

	// 创建jaeger的客户端
	// @Param 服务名，日志格式
	Jaeger, err := cfg.InitGlobalTracer("client test", jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		t.Log(err)
		return
	}
	defer Jaeger.Close()

	// 执行任务

	// 通过opentracing获取tracer
	tracer := opentracing.GlobalTracer()

	// 任务节点定义span
	parentSpan := tracer.StartSpan("A")
	defer parentSpan.Finish() // 刷新到服务器上

	B(tracer, parentSpan)
}

// 执行任务
func B(tracer opentracing.Tracer, parentSpan opentracing.Span) {
	// 创建子级span
	childSpan := tracer.StartSpan("B", opentracing.ChildOf(parentSpan.Context()))
	// 刷新到服务器上
	defer childSpan.Finish()

}

func Test_reverseWords(t *testing.T) {
	fmt.Println([]byte("blue  is sky the"))
	tests := []struct {
		input  string
		output string
	}{
		{"the sky is blue", "blue is sky the"}, {" hello world ", "world hello"}, {"a good example", "example good a"}, {"", ""}, {" a ", "a"},
	}
	for _, tt := range tests {
		result := reverseWords(tt.input)
		if result != tt.output {
			t.Errorf("expected:%s\t,result:%s\n", tt.output, result)
		} else {
			fmt.Printf("expected:%s\t,result:%s\n", tt.output, result)
		}
	}
}
func reverseWords(s string) string {
	s = strings.TrimSpace(s)
	sb := []byte(s)
	left, right := 0, 0
	for right < len(s) {
		for 0 < right && right < len(s) && s[right] == ' ' && s[right-1] == ' ' {
			right++
			continue
		}
		sb[left] = sb[right]
		right++
		left++
	}
	sb = sb[0:left]
	s = string(sb)
	ans := make([]string, 0, len(sb))
	j := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) {
			ans = append(ans, s[j:i])
			break
		}
		if sb[i] == ' ' {
			ans = append(ans, s[j:i])
			j = i + 1
		}
	}
	a := strings.Builder{}
	for i := len(ans) - 1; i >= 0; i-- {
		a.WriteString(ans[i])
		if i != 0 {
			a.WriteString(" ")
		}
	}
	return a.String()
}
