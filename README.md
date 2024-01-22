A sample api for managing contacts

### Dependencies
- [Chi Router](https://github.com/go-chi/chi)
- [pgx - PostgreSQL Driver and Toolkit](https://github.com/jackc/pgx) 
- [GoDotEnv](https://github.com/joho/godotenv) 
- [jwt-go](https://github.com/golang-jwt/jwt) 

# Installation
```bash
git clone https://github.com/Noblefel/Rest-Api-Managemen-Kontak
```  

# Usage
### Setup
Navigate inside the directory and download all the dependencies
```bash
go mod download
``` 

### ENV
Configure the environment variables
```sh
DB_HOST=localhost
DB_NAME=managemen_kontak
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
```

### Start the server
Simply run:
```sh
go run main.go
``` 
(Make sure to run the migrations)

### Note!
Set your request header to <strong>"application/x-www-form-urlencoded"</strong>
