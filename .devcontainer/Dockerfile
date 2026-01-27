FROM mcr.microsoft.com/devcontainers/go:1.25

# Arguments
ARG USERNAME=vscode
ARG KUBECTL_VERSION=1.35
ARG HELM_VERSION=3.19.0
# https://github.com/kubernetes-sigs/kubebuilder/releases
ARG KUBEBUILDER_VERSION=4.10.1
# https://github.com/kubernetes-sigs/kind/releases
ARG KIND_VERSION=v0.31.0
# Build target arch passed by BuildKit
ARG TARGETARCH

# Install tools
RUN apt-get -y update \
  && apt-get install --no-install-recommends -y \
  apt-transport-https bash-completion ca-certificates curl dnsutils file gnupg

# Install kubectl
RUN mkdir -p -m 755 /etc/apt/keyrings \
  && curl -fsSL "https://pkgs.k8s.io/core:/stable:/v${KUBECTL_VERSION}/deb/Release.key" | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg \
  && chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg \
  && echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v${KUBECTL_VERSION}/deb/ /" | tee /etc/apt/sources.list.d/kubernetes.list \
  && chmod 644 /etc/apt/sources.list.d/kubernetes.list \
  && apt-get update \
  && apt-get install -y kubectl

# Install Helm
RUN curl -fsSL "https://get.helm.sh/helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz" -o "helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz" && \
  tar -xzf "helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz" && \
  cp "linux-${TARGETARCH}/helm" /usr/local/bin/helm && \
  rm "helm-v${HELM_VERSION}-linux-${TARGETARCH}.tar.gz" && \
  rm -fr "linux-${TARGETARCH}/"

# Download kubebuilder and install locally.
RUN curl -fsSL -o /tmp/kubebuilder "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_$(go env GOOS)_$(go env GOARCH)" \
  && chmod +x /tmp/kubebuilder && mv /tmp/kubebuilder /usr/local/bin/

# Install kind
RUN go install "sigs.k8s.io/kind@${KIND_VERSION}" \
  && ln -s "/home/${USERNAME}/go/bin/kind" /usr/local/bin/kind

# Configure completion for kubectl
RUN kubectl completion bash | tee /etc/bash_completion.d/kubectl >/dev/null \
  && chmod a+r /etc/bash_completion.d/kubectl \
  && touch /etc/profile.d/kubectl.sh \
  && chmod 644 /etc/profile.d/kubectl.sh \
  && echo 'alias k=kubectl' >>/etc/profile.d/kubectl.sh \
  && echo 'complete -o default -F __start_kubectl k' >>/etc/profile.d/kubectl.sh \
  && echo ". /etc/profile.d/kubectl.sh" >>/etc/bash.bashrc

# Configure completion for kubebuilder
RUN kubebuilder completion bash | tee /etc/bash_completion.d/kubebuilder >/dev/null \
  && chmod a+r /etc/bash_completion.d/kubebuilder

# Configure completion for kind
RUN kind completion bash | tee /etc/bash_completion.d/kind >/dev/null \
  && chmod a+r /etc/bash_completion.d/kind

# Fix permissions for the /go/pkg directory
RUN chmod -R 2775 /go/pkg

ENV HISTFILE=/dotfiles/.bash_history
RUN mkdir -p -m 755 /dotfiles \
  && touch ${HISTFILE} \
  && chown ${USERNAME} ${HISTFILE} \
  && chmod 644 ${HISTFILE}

# Cleanup cache
RUN apt-get clean \
  && rm -rf /var/lib/apt/lists/*
