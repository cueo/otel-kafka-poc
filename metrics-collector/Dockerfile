# below stuff
FROM ubuntu:22.04 AS build

RUN apt-get update && apt-get install -y \
  build-essential \
  ca-certificates \
  clang-15 \
  curl \
  git \
  libelf-dev \
  libncursesw5-dev \
  libssl-dev \
  m4 \
  pkg-config \
  python3 \
  zlib1g-dev

# Install rustup
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > /rustup.sh
RUN chmod +x /rustup.sh
RUN bash /rustup.sh -y

RUN git clone https://github.com/facebookincubator/below.git && cd below && git checkout v0.7.1
# Build below
WORKDIR below
ENV CLANG=clang-15
RUN /root/.cargo/bin/cargo build --release --all-targets

FROM golang:latest

RUN apt-get update && apt-get install -y \
  libelf-dev

# Set destination for COPY
WORKDIR /app

COPY --from=build /below/target/release/below /bin/below

# Run below
#RUN nohup below record --retain-for-s 604800 --compress &

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
RUN go env -w GO111MODULE=on

# Copy the source code. Note the slash at the end, as explained in
ADD . .

# Build
RUN go build -o /collector

# Run
CMD ["./start.sh"]
