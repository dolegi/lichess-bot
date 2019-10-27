FROM alpine

WORKDIR /app

# Install dependencies
RUN apk add make g++ wget git go

# Install stockfish
RUN wget https://github.com/official-stockfish/Stockfish/archive/sf_10.tar.gz && \
      tar xvf sf_10.tar.gz && \
      cd Stockfish-sf_10/src && \
      make build ARCH=x86-64-modern && \
      mv stockfish /app && \
      rm -rf /app/Stockfish-sf_10 /app/sf_10.tar.gz

# Install lichess-bot
COPY . ./files
RUN cd files && \
  go build -o /app/lichess-bot *.go && \
  rm -rf /app/files

ENTRYPOINT ["./lichess-bot"]
