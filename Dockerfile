FROM docker.io/golang:1.22.12-alpine3.21 AS compiler
WORKDIR /app/
COPY main.go .
RUN go build -o main main.go

FROM scratch AS runner
COPY --from=compiler /app/main /main
COPY index.html .
USER 1001:1001
ENTRYPOINT ["/main"]
EXPOSE 8080/tcp