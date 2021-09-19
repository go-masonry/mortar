/*
Package confkeys defines all the configuration key alias that configure all the bundled features of mortar.
Mortar expects that a `config.yaml` or similar will be loaded as Configuration and exposed via `interfaces/cfg/Config` interface.
Mortar will access embedded keys using a `.` (dot) notation.

Example:
	Config.Get("mortar.name").String()

*******************************************************************************

Everything related to mortar itself must be under the ROOT key == "mortar"

*******************************************************************************

Expected configuration structure is a map, below is a complete example in the YAML format:

	# Root key of everything related to mortar configuration
	mortar:
		# Application/Project name
		# Type: string
		name: "Application Name"
		# Web server related configuration
		server:
			# Host is the host on which the webserver will serve APIs
			# Type: string
			host: localhost
			grpc:
				# gRPC API External port
				# Type: int
				port: 5380
			rest:
				# RESTful API External port
				# Type: int
				external:
					port: 5381
				# RESTful API Internal port
				# Type: int
				internal:
					port: 5382
		# Default Logger related configuration
		logger:
			# Set the default log level for mortar logger
			# Possible values:
			#		trace, debug, info, warn, error
			# Type: string
			level: debug
			static:
				# enables/disables adding a git commit SHA in every log entry
				# Type: bool
				git: true
				# enables/disables adding a hostname in every log entry
				# Type: bool
				host: false
				# enables/disables adding an application/project name in every log entry
				# Type: bool
				name: false

			# Log service start and stop events, custom by log level
			# Possible values:
			#		trace, debug, info, warn, error
			# Type: string
			startStop: info
		# Metrics/Monitoring related configuration
		monitor:
			# sets the namespace/prefix of every metric. Depends on the Metrics implementation
			# Type: string
			prefix: "awesome"
			# allows to include static labels/tags to every published metric
			# Type: map[string]string
			tags:
				tag1: value1
				tag2: value2
				tag3: value3
		# Bundled handlers configuration
		handlers:
			config:
				# defines a list of keywords that once contained within the configuration key will obfuscate the value
				# Type: []string
				obfuscate:
					- "pass"
					- "auth"
					- "secret"
					- "login"
					- "user"
		# Interceptors/Extractors configuration
		middleware:
			# set the default log level of all the bundled middleware that writes to log
			# Possible values:
			# 	trace, debug, info, warn, error
			# Type: string
			logLevel: "trace"
			# set the default log level of all the bundled middleware that writes to log and has an error
			# Possible values:
			# 	trace, debug, info, warn, error
			# Type: string
			logErrorLevel: "warn"
			# add Incoming gRPC request to every log entry
			# Type: bool
			logRequest: true
			# add gRPC response to every log entry
			# Type: bool
			logResponse: true
			# list of headers to be extracted from Incoming gRPC and added to every log entry
			# Type: []string
			logHeaders:
				- "x-forwarded-for"
				- "special-header"
			trace:
				http:
					client:
						# include HTTP client request to trace info ?
						# Type: bool
						request: true
						# include HTTP client response to trace info ?
						# Type: bool
						response: true
				grpc:
					client:
						# include gRPC client request to trace info ?
						# Type: bool
						request: true
						# include gRPC client response to trace info ?
						# Type: bool
						response: true
					server:
						# include incoming gRPC request to trace info ?
						# Type: bool
						request: true
						# include a gRPC response of incoming request to trace info ?
						response: true
			map:
				# List of HTTP header prefixes to map from HTTP headers to gRPC context (Incoming Metadata).
				# Useful when you have some form of correlation IDs that is passed using HTTP headers and you can access via gRPC Incomming Metadata.
				#	['requestID', 'X-Company-']
				#
				# Type: []string
				httpHeaders:
					- "X-Special-"
			copy:
				# list of header prefixes to copy/forward from Incoming gRPC context to outgoing Request context/headers
				# Type: []string
				headers:
					- "authorization"
*/
package confkeys
