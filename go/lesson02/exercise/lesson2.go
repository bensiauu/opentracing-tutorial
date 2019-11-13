package main

import (
	"context"
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/opentracing/opentracing-go/log"
	"github.com/yurishkuro/opentracing-tutorial/go/lib/tracing"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one arguement!")
	}

	//service name is hello-world
	tracer, closer := tracing.Init("testing-service")
	//operation name is hello
	span := tracer.StartSpan("Greeting")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()
	defer closer.Close()
	str := formatString(ctx, os.Args[1])
	printGreeting(ctx, str)

}

func formatString(ctx context.Context, name string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "formatString")
	defer span.Finish()

	out := fmt.Sprintf("Hello, %s!", name)

	span.LogFields(
		log.String("event", "formatting string"),
		log.String("value", name),
	)

	return out
}
func printGreeting(ctx context.Context, out string) {
	//make this span a child of the main span, the second arguement that we are currently ignoring is the new context with
	//this span. If we are calling any other functions within then we can store it and pass it down.
	span, _ := opentracing.StartSpanFromContext(ctx, "print Greeting")
	defer span.Finish()

	println(out)

	span.LogKV("event", "prining of string")
}
