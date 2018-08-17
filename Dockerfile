FROM scratch
COPY slackit /
ENTRYPOINT ["/slackit"]
USER 1000
EXPOSE 8080
EXPOSE 8443
LABEL maintainer="as@portworx.com"
