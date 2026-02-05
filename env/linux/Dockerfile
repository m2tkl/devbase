FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV DEVBASE_SKIP_CHSH=1

RUN apt-get update && \
    apt-get install -y sudo ca-certificates curl git && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /work
COPY . /work

RUN chmod +x /work/devbase.sh /work/env/linux/install.sh /work/scripts/apply_common.sh \
  /work/env/linux/settings/sysctl.sh

CMD ["/work/devbase.sh", "linux"]
