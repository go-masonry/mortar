# Inversion of Control

If you are unfamiliar with the concept please read about [it](https://en.wikipedia.org/wiki/Inversion_of_control).

Mortar heavily based on this principle. There are different libraries to achieve that in Go.
Mortar uses [uber-fx](https://github.com/uber-go/fx) and you're strongly encouraged to read all about it. Well actually you kinda have to.

To summarize, this is what you need to understand about IoC frameworks

- You are not the one creating the dependencies, Fx does it for you. Just tell it how.
- Every dependency is a Singleton, hence that same instance reused everywhere.
- Once your dependencies defined as an Interface it's really easy to swap it. Especially during tests.
- There is no magic, only implicits.

While Mortar interfaces allow you to create a custom instance of everything, but unless you really need to,
it has defaults you can influence.

## Groups

`uber-fx` [group](https://godoc.org/go.uber.org/fx#hdr-Value_Groups) is a feature that allows you to consume and produce
multiple values of the same type. This make it easier to influence/configure different instances.
Let us look at an [example](../constructors/logger.go).

```golang
type LoggerDeps struct {
	fx.In

	Config            cfg.Config
	LoggerBuilder     log.Builder
	ContextExtractors []log.ContextExtractor `group:"loggerContextExtractors"`
}

// DefaultLogger is a constructor that will create a logger with some default values on top of provided ones
func DefaultLogger(deps LoggerDeps) log.Logger {
	var logLevel = log.InfoLevel
	if levelValue := deps.Config.Get(mortar.LoggerLevelKey); levelValue.IsSet() {
		logLevel = log.ParseLevel(levelValue.String())
	}
	appName := deps.Config.Get(mortar.Name).String() // empty string is just fine
	return deps.LoggerBuilder.
		SetLevel(logLevel).
		AddStaticFields(selfStaticFields(appName)).
		AddContextExtractors(deps.ContextExtractors...).
		IncludeCallerAndSkipFrames(callerSkipDepth).
		Build()
}
``` 

Take a look at `ContextExtractors` field. You can see it has a **group** [tag](https://golang.org/pkg/reflect/#StructTag).
This tag tells Fx to gather all the instances of this type into a single slice. This slice is later "fed" to the Logger builder.
As previously mentioned [here](middleware.md) Mortar Logger can *implicitly* extract different values from the `context.Context`
and add them to the log line. To achieve it one must provide an `fx.Option` with the same group tag.
In this case `"loggerContextExtractors"`. If you read uber-fx documentation you understand what `fx.In` marker means.
You probably also read about it's *brother* `fx.Out`. However, there is a better option `fx.Annotated`.

```golang
// LoggerGRPCIncomingContextExtractorFxOption adds Logger Context Extractor using values within incoming grpc metadata.MD
//
// This one will be included during Logger build
func LoggerGRPCIncomingContextExtractorFxOption() fx.Option {
	return fx.Provide(fx.Annotated{
		Group:  groups.LoggerContextExtractors,
		Target: context.LoggerGRPCIncomingContextExtractor,
	})
}
```

When run by Fx this option adds the output of `Target` to the above slice. This way there is no need to change anything
in the Logger Builder, and you can still *influence* it.

[providers](../providers) directory have all the Options and Group Names exposed.
