FROM golang:1.24.6
ENTRYPOINT ["leetcode-tracker"]
CMD ["-c", "/config/config.yaml"]

WORKDIR /code
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install .
