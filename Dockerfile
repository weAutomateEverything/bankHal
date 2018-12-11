FROM ubuntu:14.04

ARG DT_API_URL="https://vzb12882.live.dynatrace.com/api"
ARG DT_ONEAGENT_OPTIONS="flavor=default&include=all"
ARG DT_API_TOKEN="5WUwr7a7TtOG4hSe_BC70"
ENV DT_HOME="/opt/dynatrace/oneagent"

RUN  apt-get update \
  && apt-get install -y wget openssh-client unzip \
  && rm -rf /var/lib/apt/lists/*

RUN wget -O /usr/local/share/ca-certificates/sbsapko.pem http://pko.standardbank.co.za/05766pkojnb0001_Standard%20Bank%20ROOT%20CA.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca11.pem http://pko.standardbank.co.za/05766pkojnb0011_Standard%20Bank%20Policy%20CA%2011.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca21.pem http://pko.standardbank.co.za/05766pkojnb0021_Standard%20Bank%20Policy%20CA%2021.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca111.pem http://pko.standardbank.co.za/05766PKOJNB0111.sbicdirectory.com_Standard%20Bank%20CA%20111.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca112.pem http://pko.standardbank.co.za/05766PKOJNB0112.sbicdirectory.com_Standard%20Bank%20CA%20112.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca113.pem http://pko.standardbank.co.za/05766PKOJNB0113.sbicdirectory.com_Standard%20Bank%20CA%20113.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca114.pem http://pko.standardbank.co.za/05766PKOJNB0114.sbicdirectory.com_Standard%20Bank%20CA%20114.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca211.pem http://pko.standardbank.co.za/05766PKOJNB0211.corpdirectory.com_Standard%20Bank%20Certificate%20Authority%20211.crt && \
      wget -O /usr/local/share/ca-certificates/sbsaca212.pem http://pko.standardbank.co.za/05766PKOJNB0212.corpdirectory.com_Standard%20Bank%20Certificate%20Authority%20212.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsapko.pem -inform der -out /usr/local/share/ca-certificates/sbsapko.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca11.pem -inform der -out /usr/local/share/ca-certificates/sbsaca11.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca21.pem -inform der -out /usr/local/share/ca-certificates/sbsaca21.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca111.pem -inform der -out /usr/local/share/ca-certificates/sbsaca111.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca112.pem -inform der -out /usr/local/share/ca-certificates/sbsaca112.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca113.pem -inform der -out /usr/local/share/ca-certificates/sbsaca113.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca114.pem -inform der -out /usr/local/share/ca-certificates/sbsaca114.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca211.pem -inform der -out /usr/local/share/ca-certificates/sbsaca211.crt && \
      openssl x509 -in /usr/local/share/ca-certificates/sbsaca212.pem -inform der -out /usr/local/share/ca-certificates/sbsaca212.crt && \
      update-ca-certificates

RUN mkdir -p "$DT_HOME" && \
    wget -O "$DT_HOME/oneagent.zip" "$DT_API_URL/v1/deployment/installer/agent/unix/paas/latest?Api-Token=$DT_API_TOKEN&$DT_ONEAGENT_OPTIONS" && \
    unzip -d "$DT_HOME" "$DT_HOME/oneagent.zip" && \
    rm "$DT_HOME/oneagent.zip" && \
    mkdir -p  /var/lib/dynatrace/oneagent/agent/customkeys

RUN apt-get update -qq && \
 DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
 build-essential \
 ca-certificates \
 cmake \
 curl \
 git \
 make \
 language-pack-en \
 libcurl4-openssl-dev \
 libffi-dev \
 libsqlite3-dev \
 libzmq3-dev \
 pandoc \
 python \
 python3 \
 python-dev \
 python3-dev \
 sqlite3 \
 texlive-fonts-recommended \
 texlive-latex-base \
 texlive-latex-extra \
 zlib1g-dev && \
 apt-get clean && \
 rm -rf /var/lib/apt/lists/*

RUN curl -SL -o nss_wrapper.tar.gz https://ftp.samba.org/pub/cwrap/nss_wrapper-1.1.2.tar.gz && \
 mkdir nss_wrapper && \
 tar -xC nss_wrapper --strip-components=1 -f nss_wrapper.tar.gz && \
 rm nss_wrapper.tar.gz && \
 mkdir nss_wrapper/obj && \
 (cd nss_wrapper/obj && \
 cmake -DCMAKE_INSTALL_PREFIX=/usr/local -DLIB_SUFFIX=64 .. && \
 make && \
 make install) && \
 rm -rf nss_wrapper

WORKDIR /app
# Now just add the binary
COPY cacert.pem /
COPY bankhal /app/
COPY swagger.json /app/
COPY custom.pem  /var/lib/dynatrace/oneagent/agent/customkeys/
COPY entrypoint.sh /app/

EXPOSE 8000 8080 9162

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/bankhal" ]
