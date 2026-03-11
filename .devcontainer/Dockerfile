FROM mcr.microsoft.com/devcontainers/go:dev-1.25-bookworm

# [Optional] Uncomment this section to install additional OS packages.
# [Optional] Install PostgreSQL client tools (interact with PostgreSQL from container)
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends postgresql-client

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
# RUN go get -x <your-dependency-or-tool>
RUN go install github.com/pressly/goose/v3/cmd/goose@latest   # database migration tool
RUN go install github.com/air-verse/air@latest                # cli tool for live reloading go server in dev
USER root

# RUN go mod tidy
# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1