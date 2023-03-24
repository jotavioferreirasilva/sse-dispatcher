package dispatcher_client

import (
	"backend/src/config"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetConnectionToDispatcherServer() (*grpc.ClientConn, error) {
	dispatcherServerAddress := fmt.Sprintf("%s:%d", config.BackendConfiguration.DispatcherServerHost, config.BackendConfiguration.DispatcherServerPort)
	conn, err := grpc.Dial(dispatcherServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("Failed to connect to dispatcher server: %s", err)
		return nil, err
	}

	return conn, nil
}
