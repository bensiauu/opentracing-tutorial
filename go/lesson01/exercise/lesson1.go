package main

import (
	"fmt"
	"io"
	"os"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one arguement!")
	}

	//service name is hello-world
	tracer, closer := initJaeger("hello-world")
	//operation name is hello
	span := tracer.StartSpan("hello")
	defer closer.Close()
	helloTo := os.Args[1]
	helloStr := fmt.Sprintf("Hello, %s!", helloTo)
	//the timestamp is relative to when the span started. so always put it after the operation you want to log.
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)
	span.SetTag("hello", helloTo)

	println(helloStr)

	span.LogKV("event", "println")

	span.Finish()
}

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.New(service, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
