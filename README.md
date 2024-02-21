# Dockerizing Go application  - JWT Auth Microservice
Dibawah ini adalah *step-by-step* Dockerize aplikasi *auth-jwt-microservice* menggunakan bahasa pemrograman Golang

## 1. Pastikan aplikasi telah berjalan
```
$ cd auth-jwt-microservice
$ go run main.go

output:
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /login                    --> main.handleLogin (3 handlers)
[GIN-debug] GET    /secured                  --> main.handleSecuredEndpoint (4 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-go
```
Kemudian periksa response yang diberikan *auth-jwt-microservice*
```
$ curl -X POST http://localhost:8080/login

output:
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDg1NzQ0NzEsInJvbGUiOiJhZG1pbiIsInVzZXJuYW1lIjoiZXhhbXBsZXVzZXIifQ.QJ2HLzO4wbIWK3w8kJS_8ZZulULL2eBwEoLy2dzvt3w"} 
```
## 2. Create Dockerfile.multistage
Tambahkan sintak berikut kedalam file *Dockerfile.multistage*, dimana disini menggunakan versi golang go1.20.6

```
# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /jwt-auth-microservice

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /jwt-auth-microservice /jwt-auth-microservice

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/jwt-auth-microservice"]
```
Lebih lengkap mengenai *dockerize go application* dapat mengacu pada dokumen [Build your Go image](https://docs.docker.com/language/golang/build-images/)

## 3. Building the Docker image
```
$ docker build -t jwt-auth-microservice:multistage -f Dockerfile.multistage .
```
tunggu beberapa saat hinggap proses *build* selesai

```
$ docker image ls 
```
|REPOSITORY | TAG | IMAGE ID | CREATED | SIZE|
|----- | :---- | :---- | :---- | :---- |
|jwt-auth-microservice | multistage | 8abedb9a35e4 | About a minute ago | 32MB |

Mencoba menjalankan image di dalam sebuah *container*

```
$ docker run --publish 8080:8080 jwt-auth-microservice:multistage

output:
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /login                    --> main.handleLogin (3 handlers)
[GIN-debug] GET    /secured                  --> main.handleSecuredEndpoint (4 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8080

$ curl -X POST http://localhost:8080/login

output:
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDg1Nzc2NjUsInJvbGUiOiJhZG1pbiIsInVzZXJuYW1lIjoiZXhhbXBsZXVzZXIifQ.yyPAhMs-e-srO2b_9rhSi9B4ZMePpMPQhDENIrHAEgk"}   
```
Selengkapnya mengenai *Run your Go image as a container* bisa mengacu pada dokumen https://docs.docker.com/language/golang/run-containers/

## 4. Push image
Berikut langkah-langkap melakukan upload *push* image ke *docker registry* https://hub.docker.com/
- Login ke docker hub

    ```
    $ docker login

    output:
    Authenticating with existing credentials...
    Login Succeeded

    Logging in with your password grants your terminal complete access to your account.
    For better security, log in with a limited-privilege personal access token. Learn more at https://docs.docker.com/go/access-tokens/
    ```
- Create Tag
  ```
  $ docker tag jwt-auth-microservice:multistage anangsu13/jwt-auth-microservice:v1
  $ docker image ls
  ```
  
  |REPOSITORY | TAG | IMAGE ID | CREATED | SIZE|
  |----- | :---- | :---- | :---- | :---- |
  |anangsu13/jwt-auth-microservice | v1 | 8abedb9a35e4 | 23 minutes ago  |  32MB  |
  | jwt-auth-microservice | multistage | 8abedb9a35e4 | 23 minutes ago   | 32MB |
- Docker Push
  ```
  $ docker push anangsu13/jwt-auth-microservice:v1
  ```
  tunggu beberapa saat hingga proses *push image* ke *docker registry* selesai

  ![Docker Registry](./assets/docker-registry.jpg)

  Sekarang aplikasi *jwt-auth-microservice* telah berhasil di push ke docker hub, Langkah selanjutnya adalah mencoba mengunakan image tersebut di sebuah cluster kubernetes