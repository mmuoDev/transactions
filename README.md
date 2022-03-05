# Transactions Service

The transactions service provisions an endpoint to add a transaction to a postgres database. Open API spec is in the open-api.yaml file

## Requirements

Postgres

## Usage

This service uses the generated [WalletClient](https://github.com/mmuoDev/core-proto/blob/a34d1e78d14af52b2b8915887b0d9508f758d274/gen/wallet/wallet_grpc.pb.go#L21) in making calls to the [wallet service](https://github.com/mmuoDev/wallet) via grpc.

### Starting server

To start the server, run

```bash
make run
```

### Testing

To run the tests,

```bash
make test
```
