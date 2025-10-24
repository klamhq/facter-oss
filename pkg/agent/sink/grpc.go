package sink

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/klamhq/facter-oss/pkg/options"
	schema "github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func callInventory(client schema.FactGrpcServiceClient, message *schema.InventoryRequest, logger *logrus.Logger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	b, _ := protojson.MarshalOptions{Indent: "  "}.Marshal(message)
	logger.Debugf("sending proto: %s", string(b))
	resp, err := client.Inventory(ctx, message)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			logger.Errorf("Inventory RPC failed: code=%s msg=%q", st.Code(), st.Message())
			return err
		} else {
			logger.Errorf("client.FactGrpcService(Inventory) = _, %v: ", err)
			return err
		}
	}
	logger.Infof("FactGrpcService: %s", resp.Message)
	return nil
}

func sendOverGrpc(cfg *options.FacterServerOptions, inventory *schema.InventoryRequest, logger *logrus.Logger) error {
	cert, err := tls.LoadX509KeyPair(cfg.CertificatePath, cfg.CertificateKeyPath)
	if err != nil {
		logger.Errorf("failed to load client cert: %v", err)
		return err
	}

	ca := x509.NewCertPool()
	caBytes, err := os.ReadFile(cfg.CaPath)
	if err != nil {
		logger.Errorf("failed to read ca cert %q: %v", cfg.CaPath, err)
		return err
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		logger.Errorf("failed to parse %q", cfg.CaPath)
		return err
	}

	tlsConfig := &tls.Config{
		ServerName:   cfg.SSLHostname,
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		logger.Errorf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	if err = callInventory(schema.NewFactGrpcServiceClient(conn), inventory, logger); err != nil {
		switch c := status.Code(err); c {
		case codes.PermissionDenied:
			logger.Infof("User is unauthorized to call this ressources, check your certificate SPIFFE ID: %v", err)
			return err
		default:
			logger.Errorf("Unary RPC failed unexpectedly: %v, %v", c, err)
			return err
		}
	}

	return nil
}
