FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY cacert.pem /etc/ssl/certs/ca-bundle.crt
COPY cacert.pem /
COPY bankHal /app/
COPY swagger.json /app/
COPY sampling-rules.json /
RUN apk add --no-cache openssh-client
ENTRYPOINT ["/app/bankHal"]
EXPOSE 8000 8080