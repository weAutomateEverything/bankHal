FROM alpine:latest

ARG DT_API_URL="https://vzb12882.live.dynatrace.com/api"
ARG DT_API_TOKEN="5WUwr7a7TtOG4hSe_BC70"
ARG DT_ONEAGENT_OPTIONS="flavor=musl"
ENV DT_HOME="/opt/dynatrace/oneagent"

RUN mkdir -p "$DT_HOME" && \
    wget -O "$DT_HOME/oneagent.zip" "$DT_API_URL/v1/deployment/installer/agent/unix/paas/latest?Api-Token=$DT_API_TOKEN&$DT_ONEAGENT_OPTIONS" && \
    unzip -d "$DT_HOME" "$DT_HOME/oneagent.zip" && \
    rm "$DT_HOME/oneagent.zip"
ENTRYPOINT [ "/opt/dynatrace/oneagent/dynatrace-agent64.sh" ]

WORKDIR /app
# Now just add the binary
COPY cacert.pem /etc/ssl/certs/ca-bundle.crt
COPY cacert.pem /
COPY bankhal /app/
COPY swagger.json /app/
RUN apk add --no-cache openssh-client
ENTRYPOINT ["/app/bankhal"]
EXPOSE 8000 8080