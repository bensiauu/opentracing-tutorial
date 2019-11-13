package main

import (
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

func main() {
	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		tracer, closer := tracing.Init("formatter")
		defer closer.Close()

		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		w.Write([]byte(helloStr))

		span.LogFields(
			otlog.String("event", "string formatting"),
			otlog.String("value", helloStr),
		)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
