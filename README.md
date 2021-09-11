# Mortar

![Go](https://github.com/go-masonry/mortar/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/go-masonry/mortar/branch/master/graph/badge.svg)](https://codecov.io/gh/go-masonry/mortar)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/go-masonry/mortar)](https://pkg.go.dev/mod/github.com/go-masonry/mortar)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-masonry/mortar)](https://goreportcard.com/report/github.com/go-masonry/mortar)

<table>
        <tr>
            <th><p align="left"><img src=wiki/logo.svg align="center" height=256></p></th>
            <th>
                <p align="left">Mortar is a GO framework/library for building gRPC (and REST) web services. Mortar has out-of-the-box support for configuration, application metrics, logging, tracing, profiling, dependency injection and more. While it comes with predefined defaults, Mortar gives you total control to fully customize it.            
             </p>
            </th>
        </tr>
</table>

## Demo

Clone this [demo](http://github.com/go-masonry/mortar-demo) repository to better understand some of Mortar capabilities.

When you done, read the [documentation](https://go-masonry.github.io) or create your own service using pre-made GitHub [template](#service-template) repository.

## Service Template

To help you bootstrap your services with Mortar [here](https://github.com/go-masonry/mortar-template) you can find a template. Read its README first.

## Features

- Bundled [Grpc-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) (REST Reverse-Proxy).
- Dependency Injection using [Uber-FX](https://github.com/uber-go/fx).
- Pimped `*http.Client` with interceptors support.
- Abstract support for Logging, Configuration, Tracing and Monitoring libraries. Use provided wrappers or your own.
  - [Jaeger wrapper](https://github.com/go-masonry/bjaeger) client for tracing.
  - [Prometheus wrapper](https://github.com/go-masonry/bprometheus) client for monitoring/metrics.
  - [Zerolog wrapper](https://github.com/go-masonry/bzerolog) for logging.
  - [Viper wrapper](https://github.com/go-masonry/bviper) for configuration.
- Internal HTTP [Handlers](providers/handlers.go)
  - _Profiling_ `http://.../debug/pprof`
  - _Debug_ `http://.../debug/*`
  - _Loaded Configuration_ `http://.../self/config`
  - _Build Information_ `http://.../self/build`
  - _Health_ `http://.../health`
- [Server/Client](providers) Interceptors both for gRPC and HTTP, you can choose which to use and/or add your own. 
    
    Some examples
    - HTTP Headers can be forwarded to next hop, defined by list.
    - HTTP Headers can be included in logs, defined by list.
    - Made available in `ctx.Context` via gRPC incoming Metadata.
    - Automatic monitoring and tracing (if enabled) for every RPC defined by the API.

...and more.

### Telemetry (Everything connected)

* Logs have Tracing Information `traceId=6ff7e7e38d1e86f` **across services**
    ![logs](wiki/logs.png)

* Also visible in Jaeger `traceId=6ff7e7e38d1e86f` if it's sampled.
    ![jaeger](wiki/jaeger.png)

### Support for `*http.Client` Interceptors, so you can

* Add request and response info to Trace

    <!-- ![jaeger_http](wiki/jaeger_http.png) -->

* Log/Dump requests and/or responses when http request fails.

    ```golang
    return func(req *http.Request, handler client.HTTPHandler) (resp *http.Response, err error) {
        var reqBytes, respBytes []byte
        // If the response is Bad Request, log both Request and Response
        reqBytes, _ = httputil.DumpRequestOut(req, true) // it can be nil and it's ok
        if resp, err = handler(req); err == nil && resp.StatusCode >= http.StatusBadRequest {
            respBytes, _ = httputil.DumpResponse(resp, true) // it can be nil
            logger.WithError(fmt.Errorf("http request failed")).
            WithField("status",resp.StatusCode).
            Warn(req.Context(), "\nRequest:\n%s\n\nResponse:\n%s\n", reqBytes, respBytes)
        }
        return
    }
    ```

    ![http_client](wiki/http_client_dump.png)

* Alter requests and/or responses (useful in [Tests](https://github.com/go-masonry/mortar-demo/blob/master/workshop/app/controllers/workshop_test.go#L162))

    ```golang
    func(*http.Request, clientInt.HTTPHandler) (*http.Response, error) {
        // special case, don't go anywhere just return the response
        return &http.Response{
            Status:        "200 OK",
            StatusCode:    200,
            Proto:         "HTTP/1.1",
            ProtoMajor:    1,
            ProtoMinor:    1,
            ContentLength: 11,
            Body:          ioutil.NopCloser(strings.NewReader("car painted")),
        }, nil
    }
    ```

### Monitoring/Metrics support

Export to either Prometheus/Datadog/statsd/etc, it's your choice. Mortar only provides the Interface and also **caches** the metrics so you don't have to.

```golang
counter := w.deps.Metrics.WithTags(monitor.Tags{
 "color":   request.GetDesiredColor(),
 "success": fmt.Sprintf("%t", err == nil),
}).Counter("paint_desired_color", "New paint color for car")

counter.Inc()
```

> `counter` is actually a *singleton*, uniqueness calculated [here](monitoring/registry.go#L87)

![grafana](wiki/grafana.png)

For more information about Mortar Monitoring read [here](https://go-masonry.github.io/middleware/telemetry/monitoring/).

### Additional Features

* `/debug/pprof` and other useful [handlers](handlers)
* Use `config_test.yml` during [tests](https://github.com/go-masonry/mortar-demo/blob/master/workshop/app/controllers/workshop_test.go#L151) to **override** values in `config.yml`, it saves time.

There are some features not listed here, please check the [Documentation](#documentation) for more.

## Documentation

Mortar is not a drop-in replacement.

It's important to read its [Documentation](https://go-masonry.github.io) first.
