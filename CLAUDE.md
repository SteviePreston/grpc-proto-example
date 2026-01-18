# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

This project uses mise for task management. Key commands:

```bash
mise run generate      # Generate Go code from proto files (runs install-tools first)
mise run build         # Build server and client binaries to bin/
mise run test          # Run all tests
mise run lint-proto    # Lint proto files with buf
mise run fmt-proto     # Format proto files with buf
mise run server        # Run the gRPC server (port 50051)
mise run client        # Run the gRPC client
```

Single test: `go test -v -run TestName ./server/...`

## Architecture

This is a gRPC example project with a UserService implementing CRUD operations and streaming.

**Proto → Generated Code → Server/Client flow:**
- Proto definitions: `proto/users/v1/user.proto` (package `users.v1`)
- Generated Go code: `proto-buffers/users/v1/` (package `usersv1`)
- Server implementation: `server/main.go`
- Client: `client/main.go`

**Code generation:**
- `buf.gen.yaml` configures buf to generate Go code with `paths=source_relative`
- Proto directory structure must match package name for buf lint compliance (e.g., `proto/users/v1/` → `package users.v1`)
- `go_package` option must match the output path: `proto-buffers/users/v1;usersv1`

**Server structure:**
- `userServer` struct embeds `UnimplementedUserServiceServer` and holds an in-memory user map
- RPC methods: GetUser, CreateUser, ListUsers, StreamUsers (server-side streaming)
- Uses gRPC status codes for errors (InvalidArgument, NotFound)
