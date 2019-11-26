# fluent-logger-io-writer [![GoDoc](https://godoc.org/github.com/awslabs/fluent-golang-io-writer?status.svg)](https://godoc.org/github.com/awslabs/fluent-golang-io-writer)

*This library is covered in an AWS Open Source blog post: [Splitting an application’s logs into multiple streams: a Fluent tutorial](https://aws.amazon.com/blogs/opensource/splitting-application-logs-multiple-streams-fluent/)*

This library was created to demonstrate a somewhat experimental idea. If you end up using it (or write your own similar code), [please plus one this issue to let us know](https://github.com/awslabs/fluent-golang-io-writer/issues/1). Thoughts/comments/feedback also welcome.

### What is it?

Go Code that wraps the fluent-logger-golang in a struct that implements [io.Writer](https://golang.org/pkg/io/).

This means it can be used as the underlying io stream for many loggers. See [main.go](main.go) for full example usage.

Simple example with the popular `sirupsen/logrus` logger:

```
import (
	"fmt"

	"github.com/awslabs/fluent-golang-io-writer/logger"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/sirupsen/logrus"
)

func main() {
	// configure logrus to output as JSON
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	// create a FluentWriter instance
	fluentLogger, err := logger.NewFluentWriter(fluent.Config{}, "app", []string{"level"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set logrus to use it
	logrus.SetOutput(fluentLogger)

	// Now use logrus as normal!
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}
```

```
logger.NewFluentWriter(fluent.Config{}, "app", []string{"level"})
```

Above, we see the call to the constructor for the library. The first argument is the Fluent Logger Golang config object which configures the connection to Fluentd/Bit. The second argument is a prefix for the tags given to logs emitted by this writer. The third argument is a list of keys in the log messages whose values will be the suffix of the tag. This library relies on the fact that the logs produced by Logrus will be JSON formatted. Thus, it will find the `level` key in each log message and append this to the tag prefix to construct the final tag. In practice, logs will be emitted with tags as follows:

- `app.debug` for the debug logs
- `app.info` for the info logs
- `app.warning` for the warn logs
- `app.error` for the error logs
- `app.fatal` for the fatal logs

Logs which do not have a level field or which can not be parsed as JSON will simply be given the tag `app`.

## License

This project is licensed under the Apache-2.0 License.
