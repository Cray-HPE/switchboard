# Copyright 2020 Hewlett Packard Enterprise Development LP

FROM dtr.dev.cray.com/baseos/sles15sp1 as base

RUN zypper update -y
WORKDIR /app

FROM base as build 

# Build the switchboard binary
COPY . /app
RUN zypper install -y go 
RUN export GO111MODULE=on && \
    go get && \
    go build -o switchboard main.go 

FROM base as application

# Copy switchboard binary and configs from the build layer
COPY --from=build /app/switchboard /usr/bin/switchboard
COPY --from=build /app/broker/nsswitch.conf /etc/nsswitch.conf
COPY --from=build /app/src/sshd_config /etc/switchboard/sshd_config
COPY --from=build /app/broker/entrypoint.sh /app/broker/entrypoint.sh

# Create directories for sssd
RUN mkdir -p /var/lib/sss/db \
             /var/lib/sss/keytabs \
             /var/lib/sss/mc \
             /var/lib/sss/pipes \
             /var/lib/sss/pipes/private \
             /var/lib/sss/pubconf

# Use craycli to drive switchboard logic
# This will go away when the UAS admin API is available
RUN zypper --non-interactive addrepo --no-gpgcheck \
    http://car.dev.cray.com/artifactory/shasta-premium/CLOUD/sle15_sp1_ncn/x86_64/dev/master/ \
    cloud-team && \
    zypper install -y craycli \
                      glibc-locale-base \
                      openssh \
                      sssd \
                      vim

# Set the locale for craycli
ENV LC_ALL=C.UTF-8 LANG=C.UTF-8

ENTRYPOINT ["/app/broker/entrypoint.sh"]
