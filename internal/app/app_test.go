package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	pb "github.com/mmuoDev/core-proto/gen/wallet"
	"github.com/mmuoDev/transactions/internal"
	"github.com/mmuoDev/transactions/internal/app"
	pg "github.com/mmuoDev/transactions/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

//postgresProvider mocks postgres
func postgresProvider() *pg.Connector {
	return &pg.Connector{}
}

//mockGrpcClient mocks grpc client
func mockGrpcClient() grpc.ClientConnInterface {
	return nil
}

type mockWalletClient struct {
	createWallet   func(context.Context, *pb.CreateWalletRequest, ...grpc.CallOption) (*pb.CreateWalletResponse, error)
	updateWallet   func(context.Context, *pb.UpdateWalletRequest, ...grpc.CallOption) (*emptypb.Empty, error)
	retrieveWallet func(context.Context, *pb.RetrieveWalletRequest, ...grpc.CallOption) (*pb.RetrieveWalletResponse, error)
}

func (m *mockWalletClient) CreateWallet(ctx context.Context, in *pb.CreateWalletRequest, opts ...grpc.CallOption) (*pb.CreateWalletResponse, error) {
	if m.createWallet != nil {
		return m.createWallet(ctx, in, opts...)
	}
	return nil, errors.New("client not set up")
}

func (m *mockWalletClient) UpdateWallet(ctx context.Context, in *pb.UpdateWalletRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	if m.updateWallet != nil {
		return m.updateWallet(ctx, in, opts...)
	}
	return nil, errors.New("client not set up")
}

func (m *mockWalletClient) RetrieveWallet(ctx context.Context, in *pb.RetrieveWalletRequest, opts ...grpc.CallOption) (*pb.RetrieveWalletResponse, error) {
	if m.retrieveWallet != nil {
		return m.retrieveWallet(ctx, in, opts...)
	}
	return nil, errors.New("client not set up")
}

func TestInsertTransactionWorksAsExpected(t *testing.T) {
	expectedAccount := int32(926592)
	expectedAmount := int32(900)
	expectedCategory := 1
	isDBInvoked := false
	isWalletRetrieveInvoked := false
	isWalletUpdateInvoked := false
	isWalletCreateInvoked := false

	mockDbInsert := func(o *app.OptionalArgs) {
		o.InsertTransaction = func(req internal.TransactionDBRequest) (int64, error) {
			isDBInvoked = true
			t.Run("DB data is as expected", func(t *testing.T) {
				assert.Equal(t, expectedAccount, req.AccountID)
				assert.Equal(t, expectedAmount, req.Amount)
				assert.Equal(t, expectedCategory, req.Category)
			})
			return 1, nil
		}
	}
	mockWalletClient := func(o *app.OptionalArgs) {
		wc := &mockWalletClient{}
		o.WalletClient = wc

		wc.retrieveWallet = func(c context.Context, rwr *pb.RetrieveWalletRequest, co ...grpc.CallOption) (*pb.RetrieveWalletResponse, error) {
			isWalletRetrieveInvoked = true
			var expectedResponse pb.RetrieveWalletResponse
			fileToStruct(filepath.Join("testdata", "retrieve_wallet_response.json"), &expectedResponse)
			return &expectedResponse, nil
		}
		wc.updateWallet = func(c context.Context, uwr *pb.UpdateWalletRequest, co ...grpc.CallOption) (*emptypb.Empty, error) {
			isWalletUpdateInvoked = true
			expectedNewBalance := int32(1400)
			t.Run("New balance is as expected", func(t *testing.T) {
				assert.Equal(t, uwr.CurrentBalance, expectedNewBalance)
			})
			return nil, nil
		}
		wc.createWallet = func(c context.Context, cwr *pb.CreateWalletRequest, co ...grpc.CallOption) (*pb.CreateWalletResponse, error) {
			isWalletCreateInvoked = true
			return nil, nil
		}

	}

	opts := []app.Options{
		mockDbInsert,
		mockWalletClient,
	}
	ap := app.New(postgresProvider(), mockGrpcClient(), opts...)
	serverURL, cleanUpServer := newTestServer(ap.Handler())
	defer cleanUpServer()

	reqPayload, _ := os.Open(filepath.Join("testdata", "add_transaction_request.json"))
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/transactions", serverURL), reqPayload)

	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Insert to transactions DB is invoked", func(t *testing.T) {
		assert.True(t, isDBInvoked)
	})
	t.Run("Wallet retrieve is invoked", func(t *testing.T) {
		assert.True(t, isWalletRetrieveInvoked)
	})
	t.Run("Wallet update is invoked", func(t *testing.T) {
		assert.True(t, isWalletUpdateInvoked)
	})
	t.Run("Wallet create is not invoked", func(t *testing.T) {
		assert.False(t, isWalletCreateInvoked)
	})
	t.Run("Http status code is 201", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

//newTestServer returns a test server
func newTestServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() { ts.Close() }
}

// fileToStruct reads a json file to a struct
func fileToStruct(filepath string, s interface{}) io.Reader {
	bb, _ := ioutil.ReadFile(filepath)
	json.Unmarshal(bb, s)
	return bytes.NewReader(bb)
}
