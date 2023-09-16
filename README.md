# GRPC Rest Example

This repository is an example of how to use GRPC and REST together in a single service. This gives the benefits of using a protocol buffer definition for your service, but also allows you to use REST for clients that don't support GRPC.

The backbone of this architecture is the protocol buffer which defines the service, the data interface, and the methods supported by the service. This architecture provides a language agnostic implementation of static typing data transfer. The best use cases for this type of architecture are internal services communication, aka service mesh.

The service is written in Go, but the protocol buffer definition can be used to generate clients in any language. The example client is written in Python.

## Getting Started

Install Dependencies

```bash
make install-deps
```

Produce the Protobuf Objects

```bash
make gen-code
```

Start the Server

```bash
make start-services
```
