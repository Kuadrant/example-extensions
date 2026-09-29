FROM --platform=$BUILDPLATFORM golang:1.26 AS builder
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -o threat-policy .

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/threat-policy .
ENTRYPOINT ["/threat-policy"]
