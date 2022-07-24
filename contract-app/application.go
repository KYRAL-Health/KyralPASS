package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type NetworkHandler struct {
	method *client.Contract
}

var (
	mspID         = os.Getenv("MSP_ID")
	certPath      = filepath.Join(os.Getenv("CRYPTO_PATH"), "signcerts/cert.pem")
	keyPath       = filepath.Join(os.Getenv("CRYPTO_PATH"), "keystore")
	tlsCertPath   = os.Getenv("TLS_CERT_PATH")
	peerEndpoint  = os.Getenv("PEER_ENDPOINT")
	gatewayPeer   = os.Getenv("GATEWAY_PEER")
	channelName   = os.Getenv("CHANNEL_NAME")
	chaincodeName = os.Getenv("CHAINCODE_NAME")
	eventUpdateApiKey = os.Getenv("EVENT_UPDATE_API_KEY")
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := mux.NewRouter()

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Panicln(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)
	/*
		log.Println("--> Submit: InitLedger")
		results, err := contract.SubmitTransaction("InitLedger")
		if err != nil {
			log.Fatalf("Failed to Submit transaction: %v", err)
		}
		log.Println(string(results))
	*/
	log.Println("Started")
	ContractEnable := NetworkHandler{method: contract}

	r.HandleFunc("/getAllOrgVisit", ContractEnable.getAllOrgVisit)

	r.HandleFunc("/checkAsset/{id}", ContractEnable.CheckAsset)

	r.HandleFunc("/readPrivateProperties/{id}", ContractEnable.readPrivateProperties)
	r.HandleFunc("/readPrivatePublicAsset/{id}", ContractEnable.readPrivatePublicProperties)

	r.HandleFunc("/verifyAssetProperties", ContractEnable.verifyAssetProperties).Methods("POST")

	r.HandleFunc("/createPrivate", ContractEnable.createPrivate).Methods("POST")

	r.HandleFunc("/changeDescription", ContractEnable.changeDescription).Methods("PUT")

	r.HandleFunc("/createUser", ContractEnable.createUser).Methods("POST")

	r.HandleFunc("/readUser", ContractEnable.readUser).Methods("POST")

	r.HandleFunc("/checkUser/{id}", ContractEnable.CheckUser).Methods("GET")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startChaincodeEventListening(ctx, network)

	log.Fatal(http.ListenAndServe(":1234", r))

}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		log.Panicln(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		log.Panicln(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		log.Panicln(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		log.Panicln(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := ioutil.ReadDir(keyPath)
	if err != nil {
		log.Panicln(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := ioutil.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		log.Panicln(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		log.Panicln(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		log.Panicln(err)
	}

	return sign
}

func startChaincodeEventListening(ctx context.Context, network *client.Network) {
	fmt.Println("\n*** Start chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeName)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, event.Payload)
			if event.EventName == "create_user" {
				httpClient := resty.New()
				resp, err := httpClient.R().Get(fmt.Sprintf("%s/eventUpdate?kyralID=%s&apiKey=%s", "http://192.168.86.26:8080", string(event.Payload), eventUpdateApiKey))
				if err != nil {
					log.Println(fmt.Sprintf("Error sending event update: %v", err))
				}
				if resp.StatusCode() != 200 {
					log.Println(fmt.Sprintf("API call not successful. %v", string(resp.Body())))
					return
				}
				log.Println(fmt.Sprintf("API Update call. StatusCode: %v. Body: %v", resp.StatusCode(), string(resp.Body())))
			}
		}
	}()
}
