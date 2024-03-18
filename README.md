This project is a solution for some interview question on Golang.
## 1. Description
Design and implement "Word of Wisdom" tcp server:

* TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
* The choice of the POW algorithm should be explained.
* After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
* Docker file should be provided both for the server and for the client that solves the POW challenge.

## 2. Getting started

* [Go 1.19+](https://go.dev/dl/) installed (to run tests, start server or client without Docker)
* [Docker](https://docs.docker.com/engine/install/) installed (to run docker-compose)
* Environment file `.env` (see example in [.env.example](.env.example))

### 2.1 Start server and client by docker-compose:
```
make start
```

### 2.2 Start only server:
```
make server
```

### 2.3 Start only client:
```
make client
```

### 2.4 Launch tests:
```
make test
```

## 3. Protocol definition
This protocol is based on JSON. Each message consists of two main fields: `type` and `payload`. Based on the value specified in the `type` field, we determine which structure to convert the `payload` to.

## 4. Proof of Work
The concept behind Proof of Work for DDoS protection is that a client requesting a resource from a server must first solve a challenge presented by the server. This challenge should demand significant computational work on the client side for solving it, while the server's verification of the solution should require much less computational effort.

By complicating access to server requests for the client to such an extent, abuse through mass requests will become fundamentally unprofitable.

# 4.1 Choice of an algorithm 
I have selected the Hashcash algorithm for its following advantages:
+ simplicity of implementation
+ lots of documentation and articles with description
+ simplicity of validation on server side
+ possibility to dynamically manage complexity for client by changing required leading zeros count
