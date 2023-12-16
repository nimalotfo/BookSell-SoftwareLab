FROM golang:1.19.2-bullseye

RUN apt-get update && apt-get install -y git

ARG GITLAB_TOKEN
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux TOKEN=$GITLAB_TOKEN

#RUN git config --global url."https://oauth2:${TOKEN}@gitlab.com".insteadOf "https://gitlab.com"

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan gitlab.com >> ~/.ssh/known_hosts
RUN git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
#RUN go mod download

COPY . .

RUN go env -w GOPRIVATE="gitlab.com/narm-group/*"

#RUN go build -installsuffix cgo -ldflags '-s -w' -o /gateway
RUN --mount=type=ssh go mod download -x && go mod verify
RUN --mount=type=ssh GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -tags musl -a -o /gateway

EXPOSE 8080

CMD [ "/gateway" ]

