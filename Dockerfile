# MIT License
#
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

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

RUN zypper install -y openssh \
                      sssd \
                      vim

# Set the locale for craycli
ENV LC_ALL=C.UTF-8 LANG=C.UTF-8

ENTRYPOINT ["/app/broker/entrypoint.sh"]
