FROM docker.io/library/node:25-bookworm

RUN apt-get update -y && apt-get install -y sudo

RUN npm install -g agent-browser

RUN agent-browser install --with-deps

RUN apt-get install -y dropbear

RUN npm install -g @mako10k/mcp-shell-server

CMD sh -c "echo \"root:$SSH_PASSWORD\" | chpasswd && dropbear -F -E -p 22"

