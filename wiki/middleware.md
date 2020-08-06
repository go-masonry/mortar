# Middleware

Implicitly run different logic to update/enrich/extract values from different contexts.

Examples:

- Logger can use values stored in context.Context to add them to the Log line.
- GRPC interceptors.
- Tracing can add different values stored in the context to the Span, like traffic colouring or role.
- Monitoring can add a Tag from a context or other dependencies.

While some frameworks have this feature as a first class citizen (gRPC server/client) others need to be extended.
One of the drawbacks in GO is that it doesn't have a [storage](https://github.com/golang/go/issues/21355) per go routine.
What most of us do is propagate `context.Context` in almost every function call.

```golang
func SomeFunction(ctx context.Context, param1 int, param2 string){}
// OR
func (r *http.Request) WithContext(ctx context.Context) *http.Request
```

Since we are used to it why not leverage on it ?

Let's [pretend](../interfaces/log/interfaces.go) we have a `Logger` interface, and it's defined as following:

```golang
...
Debug(ctx context.Context, format string, args ...interface{})
Info(ctx context.Context, format string, args ...interface{})
...
```

As you can see the first parameter in each `Logger` function is `context.Context`. In that case we can define the following type

```golang
type ContextExtractor func(ctx context.Context) map[string]interface{}
```

One can implement different implementations of this `ContextExtractor` type. However all of them must be registered during the building of the `Logger` instance. In our case by the `log.Builder`.

```golang
type Builder interface {
    ...
 // Make sure that each extractor function returns fast and is "thread safe"
 AddContextExtractors(hooks ...ContextExtractor) Builder
 // Build() returns a Logger implementation, always
 Build() Logger
}
```

Since we are using [Optional Builder pattern](builder.md) we can make sure that **only** at the end of the line right before building the instance one can add specific extractors based on all the information available.
To demonstrate consider the following scenario.

We made sure that every incoming HTTP Request includes [`X-Forwarded-for`](https://en.wikipedia.org/wiki/X-Forwarded-For) header which in turn will be "mapped" to `x-forwarded-for` header within the **incoming** `context.Context` by [GRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway). Now if we want to include it's value in every log we need to introduce some kind of `ContextExtractor`. But if we are already looking at http.Headers why not include every Header if it starts with a predefined prefix.

```golang
type headerPrefixes []string // if this slice will be very large it's better to build a trie map

func (h headerPrefixes) Extract(ctx context.Context) map[string]interface{} {
 var output = make(map[string]interface{})
 if md, ok := metadata.FromIncomingContext(ctx); ok {
  for key, value := range md {
   lower := strings.ToLower(key)
   for _, prefix := range h {
    if strings.HasPrefix(lower, prefix) {
     output[lower] = strings.Join(value, ",")
    }
   }
  }
 }
 return output
}
```

`headerPrefixes` type in the above example is an alias to a list of prefixes for example [**x-for**, **auth**, **custom**]
