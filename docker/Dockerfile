FROM rockylinux:9.3

# Installer les paquets nécessaires, créer les clés SSH et nettoyer yum en une seule couche
RUN yum install -y openssh scap-security-guide sudo gcc make && \
    mkdir -p /root/.ssh && \
    ssh-keygen -q -C "istia@comment" -t rsa -N '' -f /root/.ssh/id_rsa && \
    cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys && \
    ssh-keygen -q -b 4096 -C "istiaOther@comment" -t rsa -N '' -f /root/.ssh/id_rsa_other && \
    cat /root/.ssh/id_rsa_other.pub >> /root/.ssh/authorized_keys && \
    ssh-keygen -q -b 1024 -C "dsa@comment" -t dsa -N '' -f /root/.ssh/id_dsa_other && \
    cat /root/.ssh/id_dsa_other.pub >> /root/.ssh/authorized_keys && \
    ssh-keygen -b 256 -q -C "ecdsa@comment" -t ecdsa -N '' -f /root/.ssh/id_ecdsa_other && \
    cat /root/.ssh/id_ecdsa_other.pub >> /root/.ssh/authorized_keys && \
    yum clean all && rm -rf /var/cache/yum

# Installer Go et nettoyer l'archive
ENV GO_VERSION 1.24.3
RUN curl -L https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm -f /tmp/go.tar.gz

ENV PATH=/usr/local/go/bin:$PATH
ENV GOROOT=/usr/local/go

WORKDIR /tmp/facter
COPY . /tmp/facter

RUN make buildLinux

# Installer trivy
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.62.1

RUN cp /tmp/facter/bin/facter /usr/local/bin/facter

ENTRYPOINT ["/usr/local/bin/facter", "agent", "--config", "/tmp/facter/configs/config-facter-only-mac-docker.yml"]
