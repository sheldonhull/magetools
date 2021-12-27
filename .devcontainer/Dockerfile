# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.191.0/containers/codespaces-linux/.devcontainer/base.Dockerfile

FROM mcr.microsoft.com/vscode/devcontainers/universal:1-focal


# ** [Optional] Uncomment this section to install additional packages. **
# USER root
#
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>
#
# USER codespace
USER root
COPY library-scripts/*.sh library-scripts/*.env /tmp/library-scripts/
RUN echo "✅ install starship prompt" \
  && bash /tmp/library-scripts/install-starship.sh \
  && apt-get autoremove && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

USER codespace

# RUN go install github.com/jpoles1/gopherbadger@master \
#   # && go install github.com/golangci/golangci-lint/cmd/golangci-lint \
#   # && GO111MODULE=on go get mvdan.cc/gofumpt \
#   # && GO111MODULE='off'  go get golang.org/x/tools/cmd/goimports \
#   # && GO111MODULE='off' go get github.com/daixiang0/gci \
#   && go install github.com/itchyny/gojq/cmd/gojq@latest \
#   && go install github.com/owenrumney/squealer/cmd/squealer@latest \
#   && GO111MODULE=on go get github.com/zricethezav/gitleaks/v7 \
#   go install github.com/git-town/git-town@latest
