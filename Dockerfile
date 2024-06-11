FROM node:18 as frontend
WORKDIR /app
COPY frontend/package-lock.json /app
COPY frontend/package.json /app
RUN npm install
COPY frontend/src /app/src
COPY frontend/public /app/public
COPY frontend/tsconfig.json /app/tsconfig.json
COPY frontend/vite.config.ts /app/vite.config.ts
COPY frontend/index.html /app/index.html
ENV REACT_APP_ENV=production
RUN npm run build

FROM golang:1.21.1 as backend
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
WORKDIR /app
COPY pkg /app/pkg
COPY cmd /app/cmd
COPY go.mod /app
COPY go.sum /app

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o bin/server cmd/server/main.go
RUN chmod +x bin/server

FROM alpine:3
RUN apk add --no-cache curl
COPY --from=backend /app/bin/server /usr/local/bin
WORKDIR /app
COPY --from=frontend /app/dist /app/public
USER 1001
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD curl -f http://localhost:8080/health || exit 1
ENTRYPOINT ["server"]  