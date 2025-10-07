package gogrpc

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InsecureGrpcConnection establishes a gRPC client connection to the specified target
// using insecure credentials (no TLS). This is typically used for Unix domain sockets
// or trusted internal environments.
//
// If the connection attempt fails, it logs the error as fatal using the provided logger
// and immediately terminates the application.
//
// Parameters:
//   - socket: the address of the gRPC server (e.g., "unix:///var/run/my.sock" or "localhost:50051").
//   - logger: a zap.Logger used for error reporting.
//
// Returns:
//   - A pointer to a grpc.ClientConn that can be used to create service clients.
//
// Example:
//
//	conn := InsecureGrpcConnection("unix:///var/run/service.sock", logger)
//	client := mypb.NewMyServiceClient(conn)
func InsecureGrpcConnection(
	socket string,
	logger *zap.Logger,
) *grpc.ClientConn {
	connection, err := grpc.NewClient(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

		logger.Fatal("failed to connect to gRPC socket", zap.Error(err))
	}
	return connection
}
