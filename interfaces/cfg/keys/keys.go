package confkeys

// Root
const mortar = "mortar"

// Root level keys
const (
	// ApplicationName is the name of this Application/Project
	//
	// Type: string
	ApplicationName string = mortar + ".name"

	// Webserver specific configurations
	server = mortar + ".server"
	// Mortar Logger configuration
	logger = mortar + ".logger"
	// Mortar monitoring configuration
	monitor = mortar + ".monitor"
	// Mortar middleware, interceptors
	middleware = mortar + ".middleware"
	// Mortar bundled handlers
	handlers = mortar + ".handlers"
)

// Webserver related keys
const (
	// Webserver -> gRPC related configuration
	gRPC = server + ".grpc"
	// Webserver -> RESTful related configuration
	rest = server + ".rest"

	// Host is the host on which the webserver will serve APIs
	//
	// Type: string
	Host = server + ".host"

	// ExternalGRPCPort is the Port on which the webserver will serve gRPC API
	//
	// Type: int
	ExternalGRPCPort string = gRPC + ".port"

	// ExternalRESTPort is the Port on which the webserver will serve it's external/public RESTful API
	//
	// Type: int
	ExternalRESTPort string = rest + ".external.port"

	// InternalRESTPort is the Port on which the webserver will serve it's internal/private RESTful API
	//
	// Type: int
	InternalRESTPort string = rest + ".internal.port"
)

// Logger related keys
const (
	// LogLevel set the default log level for mortar logger
	// Possible values:
	//		trace, debug, info, warn, error
	//
	// Type: string
	LogLevel string = logger + ".level"

	// Log service start and stop events, custom by log level
	// Possible values:
	//		trace, debug, info, warn, error
	// Type: string
	LogStartStop string = logger + ".startStop"
	// LogIncludeGitSHA enables/disables adding a git commit SHA in every log entry
	//
	// Type: bool
	LogIncludeGitSHA string = logger + ".static.git"

	// LogIncludeHost enables/disables adding a hostname in every log entry
	//
	// Type: bool
	LogIncludeHost string = logger + ".static.host"

	// LogIncludeName enables/disables adding an application/project name in every log entry
	//
	// Type: bool
	LogIncludeName string = logger + ".static.name"
)

// Monitoring related keys
const (
	// MonitorPrefix sets the namespace/prefix of every metric. Depends on the Metrics implementation
	//
	// Type: string
	MonitorPrefix string = monitor + ".prefix"

	// MonitorTags allows to include static labels/tags to every published metric
	//
	// Example:
	//		tags:
	//			tag1: value1
	//			tag2: value2
	//			tag3: value3
	//
	// Type: map[string]string
	MonitorTags string = monitor + ".tags"
)

// Bundled Handlers
const (
	// ConfigHandlerObfuscateKeys defines a list of keywords that once contained within the configuration key will obfuscate the value
	//
	// Type: []string
	ConfigHandlerObfuscateKeys = handlers + ".config.obfuscate"
)

// Middleware
const (
	// MiddlewareLogLevel set the default log level of all the bundled middleware that writes to log
	// Possible values:
	//		trace, debug, info, warn, error
	//
	// Type: string
	MiddlewareLogLevel = middleware + ".logLevel"

	// MiddlewareOnErrorLogLevel set the default log level of all the bundled middleware that writes to log and has an error
	// Possible values:
	//		trace, debug, info, warn, error
	//
	// Type: string
	MiddlewareOnErrorLogLevel = middleware + ".onErrorLogLevel"

	// MiddlewareLogIncludeRequest add incoming gRPC request to log
	//
	// Type: bool
	MiddlewareLogIncludeRequest = middleware + ".logRequest"

	// MiddlewareLogIncludeResponse add gRPC response of incoming request to log
	//
	// Type: bool
	MiddlewareLogIncludeResponse = middleware + ".logResponse"

	// trace related keys with http context
	traceHTTP = middleware + ".trace.http"

	// HTTPClientTraceIncludeRequest add HTTP client request to trace info ?
	//
	// Type: bool
	HTTPClientTraceIncludeRequest = traceHTTP + ".client.request"

	// HTTPClientTraceIncludeResponse add HTTP client response to trace info ?
	//
	// Type: bool
	HTTPClientTraceIncludeResponse = traceHTTP + ".client.response"

	// trace related keys with grpc context
	traceGRPC = middleware + ".trace.grpc"

	// GRPCClientTraceIncludeRequest add gRPC client request to trace info
	//
	// Type: bool
	GRPCClientTraceIncludeRequest = traceGRPC + ".client.request"

	// GRPCClientTraceIncludeResponse add gRPC client response to trace info
	//
	// Type: bool
	GRPCClientTraceIncludeResponse = traceGRPC + ".client.response"

	// GRPCServerTraceIncludeRequest add incoming gRPC request to trace info
	//
	// Type: bool
	GRPCServerTraceIncludeRequest = traceGRPC + ".server.request"

	// GRPCServerTraceIncludeResponse add a gRPC response of incoming request to trace info
	//
	// Type: bool
	GRPCServerTraceIncludeResponse = traceGRPC + ".server.response"

	// ForwardIncomingGRPCMetadataHeadersList is a list of header prefixes to copy from Incoming gRPC context to outgoing Request context/headers
	//
	// Type: []string
	ForwardIncomingGRPCMetadataHeadersList = middleware + ".copy.headers"

	// List of HTTP header prefixes to map from HTTP headers to gRPC context (Incoming Metadata).
	// Useful when you have some form of correlation IDs that is passed using HTTP headers and you can access via gRPC Incomming Metadata.
	//	['requestID', 'X-Company-']
	//
	// Type: []string
	MapHTTPRequestHeadersToGRPCMetadata = middleware + ".map.httpHeaders"

	// LoggerIncomingGRPCMetadataHeadersExtractor is a list of headers to be extracted from Incoming gRPC and added to every log entry
	//
	// Type: []string
	LoggerIncomingGRPCMetadataHeadersExtractor = middleware + ".logHeaders"
)
