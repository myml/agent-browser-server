FROM docker.io/library/golang:bookworm as builder
COPY . /src
WORKDIR /src
RUN go build -o server

FROM docker.io/library/node:25-bookworm
RUN apt-get update -y && apt-get install -y sudo
RUN npm install -g agent-browser
RUN agent-browser install --with-deps
RUN apt-get install -y dropbear
COPY --from=builder /src/server /server
ENV SSH_PASSWORD=root
CMD sh -c "echo \"root:$SSH_PASSWORD\" | chpasswd && dropbear -E -p 22 && /server"

