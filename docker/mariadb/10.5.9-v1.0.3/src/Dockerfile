FROM mariadb:10.5.9

RUN set -eux;  \
    apt-get update;  \
    apt-get -y --no-install-recommends install ncat; \
    rm -rf /var/lib/apt/lists/*;

ADD https://github.com/powerpuffpenguin/cronbackup/releases/download/v1.0.3/linux.amd64.tar.gz /opt/cronbackup/

RUN set -eux;  \
    tar -zxvf /opt/cronbackup/linux.amd64.tar.gz -C /opt/cronbackup/; \
    rm /opt/cronbackup/linux.amd64.tar.gz;