FROM ubuntu:14.04

#ARG DT_API_URL="https://vzb12882.live.dynatrace.com/api"
#ARG DT_API_TOKEN="5WUwr7a7TtOG4hSe_BC70"
#ARG DT_ONEAGENT_OPTIONS="flavor=default&technology=go"
#ENV DT_HOME="/opt/dynatrace/oneagent"

#RUN  apt-get update \
#  && apt-get install -y wget openssh-client unzip \
#  && rm -rf /var/lib/apt/lists/*

#RUN mkdir -p "$DT_HOME" && \
#    wget -O "$DT_HOME/oneagent.zip" "$DT_API_URL/v1/deployment/installer/agent/unix/paas/latest?Api-Token=$DT_API_TOKEN&$DT_ONEAGENT_OPTIONS" && \
#    unzip -d "$DT_HOME" "$DT_HOME/oneagent.zip" && \
#    rm "$DT_HOME/oneagent.zip"
#ENTRYPOINT [ "/opt/dynatrace/oneagent/dynatrace-agent64.sh" ]

RUN wget -O Dynatrace-OneAgent-Linux-1.149.188.sh "https://vzb12882.live.dynatrace.com/api/v1/deployment/installer/agent/unix/default/latest?Api-Token=Zx-S_6-MTcO7xhPe4_t1l&arch=x86&flavor=default"
RUN sudo Dynatrace-OneAgent-Linux-1.149.188.sh  APP_LOG_CONTENT_ACCESS=1 INFRA_ONLY=0

WORKDIR /app
# Now just add the binary
COPY cacert.pem /etc/ssl/certs/ca-bundle.crt
COPY cacert.pem /
COPY bankhal /app/
COPY swagger.json /app/
CMD ["/app/bankhal"]
EXPOSE 8000 8080