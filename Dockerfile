FROM scratch
COPY slackit /
RUN ["/slackit", "mktmp"]
ENTRYPOINT ["/slackit"]
USER 65534
EXPOSE 8080
EXPOSE 8443
LABEL maintainer="as@portworx.com"
