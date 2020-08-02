package mortar

/*
Keys in this file assume the following configuration structure (example is in YAML, but can be TOML/JSON)
```
mortar:
	name: "awesome project"
	server:
		grpc:
			port: 5380
		rest:
			external:
				port: 5381
			internal:
				port: 5382
	logger:
		level: info
		console: false
	monitor:
		address: "host:port"
		prefix: "awesome.project"
		tags:
			tag1: value
			tag2: value
			tag3: value
	middleware:
		rest:
			client:
				trace:
					request: true
					response: true
		grpc:
			server:
				log:
					request: true
					response: false
					level: debug
				trace:
					request: true
					response: true
				headers:
					- authorization
			client:
				trace:
					request: true
					response: true
		logger:
			headers:
				- x-forwarded-for
				- special-header
	handlers:
		health:
			timeout: 10s
		self:
			obfuscate:
				- "pass"
				- "auth"
				- "secret"
				- "login"
				- "user"
```
*/
const (
	Name = "mortar.name"
	// Server
	ServerGRPCPort         = "mortar.server.grpc.port"
	ServerRESTInternalPort = "mortar.server.rest.internal.port"
	ServerRESTExternalPort = "mortar.server.rest.external.port"
	ServerGRPCLogLevel     = "mortar.server.grpc.log.level"
	// Logger
	LoggerLevelKey      = "mortar.logger.level"
	LoggerWriterConsole = "mortar.logger.console"

	// Monitoring
	MonitorAddressKey = "mortar.monitor.address"
	MonitorTagsKey    = "mortar.monitor.tags"
	MonitorPrefixKey  = "mortar.monitor.prefix"

	// Middleware
	//// Server
	MiddlewareServerGRPCLogIncludeRequest    = "mortar.middleware.grpc.server.log.request"
	MiddlewareServerGRPCLogIncludeResponse   = "mortar.middleware.grpc.server.log.response"
	MiddlewareServerGRPCTraceIncludeResponse = "mortar.middleware.grpc.server.trace.response"
	MiddlewareServerGRPCTraceIncludeRequest  = "mortar.middleware.grpc.server.trace.request"
	MiddlewareServerGRPCCopyHeadersPrefixes  = "mortar.middleware.grpc.server.headers"
	//// Client
	MiddlewareClientGRPCTraceIncludeRequest  = "mortar.middleware.grpc.client.trace.request"
	MiddlewareClientGRPCTraceIncludeResponse = "mortar.middleware.grpc.client.trace.response"
	MiddlewareClientRESTTraceIncludeRequest  = "mortar.middleware.rest.client.trace.request"
	MiddlewareClientRESTTraceIncludeResponse = "mortar.middleware.rest.client.trace.response"
	///// Logger
	MiddlewareLoggerHeaders = "mortar.middleware.logger.headers"
	// Handlers
	HandlersSelfObfuscateConfigKeys = "mortar.handlers.self.obfuscate"
	HandlersHealthTimeout           = "mortar.handlers.health.timeout"
)
