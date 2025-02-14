# Docs

Hi, this is `general guidance` for run the `auth` service:

👉 Watch step by step for `run the project`,

## 📖 General Guidance

[![License](https://img.shields.io/badge/License-Dinoco%20-red.svg)]()
[![Go](https://img.shields.io/badge/go-1.23-green.svg)](https://golang.org/)

### Table of Contents
- [Getting started](#-getting-started)
- [Requirements](#requirements)
- [Docker setup](#docker-setup)
- [Project setup](#project-setup)
- [Testing](#testing)
- [Explore](#explore)

## 🏃 Getting Started
Hello, thank you for reading the documentation. Don't forget to sit tightly, take a deep breath, and drink a coffee 🍺.

### Requirements
The minimum requirements **you must have** is:
1. `Linux` or `Mac OS` machine 
2. Already have IDE like goland or something else 
3. Already Database Manager, in my case i am using dbeaver
4. Already install postman
5. `4 GB of RAM`. Higher is better, 
6. The program must be already installed on your machine:
    - Go
    - PostgreSQL
    - Redis
    - Docker, docker-compose

### Docker Setup
Open the terminal and inside the project run docker compose using this command: 
```
docker-compose up -d
```
Make sure the container postgres and redis is up, you can see the status container using this:
```
docker ps -a
```

### Project Setup
Open your database manager to create connection and fill this value `localhost`:`54320` and username `admin` and password `secret`.
After connection is success, you can right-click and create new database with name `dsi_auth` and `dsi_auth_test`
Refresh the connection and make sure the database is already exist.

Open the golang project and exec this command:
```
go run main.go
```
or
```
make serve
```

After the project running, it will run the auto migration to create table `users`.

### Testing
On this project already have unit testing. Make sure your already have `dsi_auth_test` and container postgres and redis is up.
You can exec this command:
```
make test-coverage
```

### Explore
You can explore the API from my postman collection. Thank you for read me