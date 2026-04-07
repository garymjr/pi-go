# AGENTS.md

- Module path is `github.com/garymjr/pi-go`; the repo targets Go 1.22.
- Public packages are `rpc` (subprocess/RPC client) and `wire` (protocol types, framing, codecs).
- `internal/testproc` is test-only. `rpc/client_test.go` builds it on the fly with `go build ./internal/testproc`, so unit tests do not need a real `pi` binary.
- Use `go test ./...` for full verification.
- Use `go test ./rpc -run TestName` or `go test ./wire -run TestName` for a focused test.
- Examples live in `examples/basic` and `examples/extensionui`; run them with `go run ./examples/basic` or `go run ./examples/extensionui`.
- If README text and code disagree, trust the code and tests.
- `wire` preserves unknown frames and unknown UI requests as raw JSON for forward compatibility; keep that behavior.
- The decoder accepts LF, CRLF, and blank lines even though the README says LF-only; do not tighten framing unless you mean to change the protocol.
- `rpc.Client` is exercised concurrently in tests; preserve request/response correlation and notification dispatch safety when changing start/read/write paths.
- `rpc.Options` defaults: executable `pi`, response timeout 30s, shutdown timeout 2s, startup delay 100ms, stderr buffer 64KiB.
