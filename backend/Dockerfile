FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o api

# ---- run stage ----
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /home/nonroot
ENV PORT=8080
COPY --from=build /app/api ./api
USER nonroot
EXPOSE 8080
ENTRYPOINT ["./api"]