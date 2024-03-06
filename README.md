### Dependencies
- [Chi Router](https://github.com/go-chi/chi)
- [pgx - PostgreSQL Driver and Toolkit](https://github.com/jackc/pgx) 
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

### Command Flags 
| Key | Default |
| -------- | ------- |
| host | localhost | 
| port | 5432 | 
| name | managemen_kontak | 
| u | postgres | 
| pw |  | 


### Start the server
Using default configurations, simply run:
```sh
go run cmd/main.go
``` 

With flags: 
```sh
go run cmd/main.go -host=localhost -port=5432 -name=managemen_kontak -u=postgres -pw={Your password}
```

(Make sure to run the migrations)

### Note
Set the request header to <strong>"application/x-www-form-urlencoded"</strong>
