# Manifold Project Development Guidelines

This document provides essential information for developers working on the Manifold project.

## Build/Configuration Instructions

### Prerequisites

1. **Go Version**: The project requires Go 1.24.2 or later. The toolchain is configured for Go 1.24.4.

2. **External Dependencies**: The project depends on the libsecp256k1 library for cryptographic operations. To install it on Ubuntu:

   ```bash
   # Run the provided script
   ./ubuntu_install_libsecp256k1.sh
   ```

   This script:
   - Installs build dependencies (build-essential, autoconf, libtool)
   - Clones the secp256k1 repository from Bitcoin Core
   - Configures and builds the library with specific modules enabled:
     - schnorrsig
     - ecdh
     - ellswift

### Building the Project

The project uses standard Go build tools:

```bash
# Get dependencies
go mod tidy

# Install to $GOBIN
go install ./path/to/main/package
```

## Testing Information

### Running Tests

Tests are written using the standard Go testing package. To run tests:

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./package/path

# Run a specific test
go test ./package/path -run TestName

# Use `-v` to print more of the output of the test
go test -v ./path/to/package 
```

### Test Structure

Tests follow standard Go testing conventions:

1. Test files are named with the `_test.go` suffix
2. Test functions are prefixed with `Test`
3. Benchmark functions are prefixed with `Benchmark`

### Error Handling in Tests

The project uses a custom error checking mechanism through the `chk` package:

```go
import "manifold.mleku.dev/chk"

// Check if an error occurred and log it
if chk.E(err) {
    t.Fatal("An error occurred")
}
```

The `chk` package provides different log levels:, Fatal, Error, Warning, Info, Debug, Trace, using F, E, W, I, D, T.

The `log` package provides a way to print things anywhere in the code, and uses the same second field: `log.E.` as `chk`
but then a further set of options exists: `Ln` for `fmt.Println` style, `F` for `fmt.Printf` style, `C` to place a 
stringer function as parameter (`func() string`), and `S`, which uses `spew.Sdump` to construct a line-separated list of 
values to print, usually using reflection to dig inside the unexported fields.


### Example Test

Here's an example of a simple test:

```go
package relay

import (
    "testing"
    "manifold.mleku.dev/chk"
)

func TestNewMessage(t *testing.T) {
    // Test creating a new message
    msg := NewMessage("sender", "content")
    if msg.Sender != "sender" {
        t.Errorf("Expected sender to be 'sender', got '%s'", msg.Sender)
    }
    if msg.Content != "content" {
        t.Errorf("Expected content to be 'content', got '%s'", msg.Content)
    }
}
```

## Additional Development Information

### Project Structure

The project is organized into multiple packages, each with a specific purpose:

- `chk`, `log`, `errorf`: Convenience shortcuts for logging
- `lol`: Custom logging library with source location tracking
- `ints`: Optimized encoder for decimal numbers in ASCII format
- `ec`: Elliptic curve cryptography implementations
...and so on

### Logging System

The project uses a custom logging library (`lol` - log of location) that:

1. Prints high-precision timestamps and source locations
2. Supports multiple log levels (Fatal, Error, Warn, Info, Debug, Trace)
3. Uses colored output for better readability
4. Provides convenient error checking mechanisms

To use the logging system:

```go
import "manifold.mleku.dev/lol"

// Set the log level
lol.SetLogLevel("debug")

// Log at different levels
log.I.F("This is an info message with format: %s", "value")
log.D.Ln("This is a debug message")
log.E.S(complexObject) // Dumps the object using spew
chk.E(err) // prints an error and returns true if there is an error
errorf.E("something with %d printf formatting", somenumber) // which creates an error variable as well as printing it
```

### Code Style

The project follows these coding conventions:

1. **Documentation**: All exported functions, types, and variables should have proper documentation comments
2. **Error Handling**: Use the `chk` package for error checking in tests
3. **Logging**: Use the appropriate log level from the `lol` package
4. **Testing**: Write comprehensive tests for all functionality
5. **Use named return variables, and naked returns**: Don't make it hard to find what is a return value. DRY.
6. **Don't write long functions**: if it's more than 300 lines long, you probably should split it into parts.

### Cryptographic Operations

The project uses the libsecp256k1 library for cryptographic operations, particularly for:
- Schnorr signatures
- ECDH (Elliptic Curve Diffie-Hellman)
- EllSwift (todo: not yet implemented. will replace ECDH from the btcec package—which is repeated and merged with the decred cryptography package in /ec)

Ensure the library is properly installed before working with cryptographic functions. There is a script for building it 
on ubuntu in the root of the repository. Aside from the ways of installing the dependencies, the remainder should be the 
same for any POSIX compliant operating system.