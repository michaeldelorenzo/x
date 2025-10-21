# Service Runner
The config package contains common interfaces to configure service parameters for the Go programming language.

## Usage

### Run a Service
The service runner packages exposes an API for defining a new service runner and starting it.

```go
package main

import (
	"github.com/Blooming-Health/my-service/config"
	"github.com/Blooming-Health/my-service/internal/api/transport"
	"github.com/michaeldelorenzo/x/pkg/servicerunner"
	"github.com/urfave/cli"
)

func cliAction(_ *cli.Context, instr instrumenting.Instrumentor) {
	svc, endpoints := initializeService(instr)
	shutdown := transport.RunGRPCServer(svc, endpoints, instr, config.AppConfig.Server)

	runner := servicerunner.NewServiceRunner(func() {
		instr.Log("serve", "Running RPC server!")
	}).SetCleanup(func() {
		instr.Log("shutdown", "gRPC server shutting down gracefully")
		shutdown()
	})

	runner.Start()
}
```