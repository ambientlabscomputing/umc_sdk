# UMC SDK

Shared SDK for Underleaf Agent Managed Components (UMCs).

## Overview

This SDK provides common patterns and utilities for building UA-Managed Components. It extracts patterns proven in the Mycelium Mesh Agent (MMA) and provides a consistent foundation for all edge components.

## Architecture

UMCs are userland subsystems supervised by the UA Kernel (UA-K). They:
- Run as separate processes
- Communicate with UA-K via gRPC over Unix domain sockets
- Use syscalls to request privileged operations
- Subscribe to events from the kernel
- Follow a standard lifecycle pattern

## Packages

### `transport/`
Unix domain socket helpers, `SO_PEERCRED` authentication, and mTLS credential management.

### `lifecycle/`
Standard Launcher + Runtime pattern for ordered initialization, teardown, and signal handling.

### `config/`
Thread-safe configuration store with atomic read/write operations.

### `events/`
gRPC event consumer with automatic reconnection and sequence-based resumption.

### `telemetry/`
Buffered telemetry collection and periodic flushing to UA-K.

### `logging/`
Structured logging with slog and automatic log rotation.

### `version/`
Build-time version information injection.

### `types/`
Common types, error codes, and enums shared across UMCs.

### `syscall/`
Generated gRPC client stubs for the UA-K syscall API.

## Usage

```go
import (
    "github.com/ambientlabscomputing/umc_sdk/lifecycle"
    "github.com/ambientlabscomputing/umc_sdk/syscall"
)

func main() {
    launcher := lifecycle.NewLauncher("my-umc")
    
    // Initialize your components
    launcher.Add("server", myServer)
    
    // Run with signal handling
    runtime := lifecycle.NewRuntime(launcher)
    runtime.Run()
}
```

## Building UMCs

See the [UMC Developer Guide](../underleaf_client/agent_docs/UMC_DEVELOPER_GUIDE.md) for detailed instructions on creating new UA-Managed Components.
