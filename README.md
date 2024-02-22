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
# Test login
$ curl -X POST http://localhost:8080/login
{"token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJkZW1vQGV4YW1wbGUuY29tIiwiZXhwIjoxNzA4NjU1NDQxLCJ1c2VyX2lkIjoiMTEyMjMzNDQifQ.LiboR4YvRXPSmYBuJE-o35wn8h-TUIeCTwzOw2jk4xA9elOGrBU3x6Mn7uG4LPQDDKqGqBeoF5-L2ov_rIEcSjRQEvFK55O-lawahGxKqbhBCsRwBJFq8Qzcx3JQYsW1Af0Ds5U9IIlwWxhTYHvWEYKDhswNgjiAo2xBE_XDl8D0f7ccF4_knm9CUSihg9UEEMw14hABoNZY2Del-_bg24FsjcgA4MoTNc5FhL72BqLh200kG85YXVbc2TLAwOqSFQ5rlB0UU3dxEiyBxmUYJnKotkOC8tSfos_IigU9sI6zz96KBzZLJmaWqgFFbG-t7ofl-gh02ttYDqBlDLw0cQ","user":{"email":"userdemo@example.com","first_name":"Example","last_name":"User","middle_name":"Demo","user_id":"11223344"}}

# Test jwks endpoint
$ curl http://localhost:8080/jwks
{"keys":[{"alg":"RS256","e":"AQAB","exp":1708655531,"kid":"1","kty":"RSA","n":"e0439b8c5b32a3fb681a6dee104821f8ed28a8bb7bf51daf111427cd173984a83d1fd9a273699d9368ae0c583fc3702837f306f17135b07e2f0d9d8e9cc83cbca0e9ef58217912b8b34c7191bea867b42381f0c9a298d591e937106f4fa6ca7270236e6421f00b5c5b18fc6f8d7ad215226a8fecf2adc9f1889214981a038b3647e2016d2b91a57f90fddcdd8d48dbc03ec58adfdb72e8f38e24ae418789e720a6209cc12420af9012a57fd92034e2ee5d587ccb5bc622a53e25873fdabb86b8d8318e1cd98440f990ea947277b443efcc54c101801c9721458a2cc340b890194457fb7598edcc9d8812105a6d983f3d19be91dc2e9b84c7b1507db0347e8a99","use":"sig"}]}
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