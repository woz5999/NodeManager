FROM golang:1.5-onbuild

ENV KUBERNETES_VERSION="v1.7.5"
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl


EXPOSE 80
