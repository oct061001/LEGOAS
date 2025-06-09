# LEGOAS
repo for LEGOAS test
# Employee System Backend (LEGOAS)

This project is the backend service for the Employee Account Management System using **gRPC** with **Go** and **MongoDB**.

---

## Tech Stack

- Language: Go (Golang)  
- Communication: gRPC with Protobuf  
- Database: MongoDB  
- Dependencies:  
  - `google.golang.org/grpc`  
  - `go.mongodb.org/mongo-driver`  

---

## Features

- Account registration with user data and roles  
- Role and permission management  
- MongoDB integration for persistent storage  
- gRPC service for efficient and strongly typed communication  

---

## Getting Started

### Prerequisites

- Go 1.24+  
- MongoDB instance running locally or remotely  
- `protoc` and `protoc-gen-go` installed for protobuf generation  

### Running Locally

1. Clone the repository:

```bash
git clone https://github.com/oct061001/LEGOAS.git
cd LEGOAS
