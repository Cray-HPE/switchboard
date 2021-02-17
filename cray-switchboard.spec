#
# RPM spec file for switchboard
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
%define cmdname switchboard

Name: cray-%{cmdname}
License: Cray Software License Agreement
Summary: %{cmdname}
Version: %(cat .version)
Release: %(echo ${BUILD_METADATA})
Source: %{name}-%{version}.tar.bz2
Vendor: Cray Inc.
Group: Productivity/Clustering/Computing

BuildRequires: systemd

Requires: craycli

%systemd_requires

%description
This package provides files for deploying the Cray Switchboard
project.

%files
%dir %{_unitdir}
%dir %{_bindir}
%dir %{_sysconfdir}
%dir %{_sysconfdir}/%{cmdname}
%{_unitdir}/cray-%{cmdname}-sshd.service
%{_bindir}/%{cmdname}
%{_sysconfdir}/%{cmdname}/sshd_config
%{_sysconfdir}/%{cmdname}/ssh

%prep
%setup -q

%build
export GO111MODULE=on
go get
go build -o switchboard main.go

%install
%{__mkdir_p} %{buildroot}%{_unitdir}
%{__mkdir_p} %{buildroot}%{_bindir}
%{__mkdir_p} %{buildroot}%{_sysconfdir}/%{cmdname}

%{__install} -m 0644 src/cray-%{cmdname}-sshd.service %{buildroot}%{_unitdir}/cray-%{cmdname}-sshd.service
%{__install} -m 0755 %{cmdname} %{buildroot}%{_bindir}/%{cmdname}
%{__install} -m 0700 src/sshd_config %{buildroot}%{_sysconfdir}/%{cmdname}/sshd_config
%{__install} -m 0755 src/ssh %{buildroot}%{_sysconfdir}/%{cmdname}/ssh

%post
%if 0%{?suse_version}
%service_add_post cray-%{cmdname}-sshd.service
%else
%systemd_post cray-%{cmdname}-sshd.service
%endif

%preun
%if 0%{?suse_version}
%service_del_preun cray-%{cmdname}-sshd.service
%else
%systemd_preun cray-%{cmdname}-sshd.service
%endif

%postun
%if 0%{?suse_version}
%service_del_postun cray-%{cmdname}-sshd.service
%else
%systemd_postun_with_restart cray-%{cmdname}-sshd.service
%endif

%changelog
