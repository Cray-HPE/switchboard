#
# RPM spec file for switchboard
# Copyright 2019 Cray Inc. All Rights Reserved.
#
%define crayctl_dir /opt/cray/crayctl
%define ansible_dir %{crayctl_dir}/ansible_framework
%define cmdname switchboard

Name: cray-%{cmdname}
License: Cray Software License Agreement
Summary: %{cmdname}
Version: %(cat .version)
Release: %(echo ${BUILD_METADATA})
Source: %{name}-%{version}.tar.bz2
Vendor: Cray Inc.
Group: Productivity/Clustering/Computing

BuildRequires: cme-premium-cf-crayctldeploy-buildmacro
BuildRequires: systemd

Requires: craycli

%systemd_requires

%description
This package provides an ansible roles for deploying the Cray Switchboard
project.

%files
%{ansible_dir}
%{cme_premium_roles_dir}
%{crayctl_dir}
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
# Install ansible files
%{__mkdir_p} %{buildroot}%{cme_premium_roles_dir}
%{__mkdir_p} %{buildroot}%{crayctl_dir}
%{__mkdir_p} %{buildroot}%{_unitdir}
%{__mkdir_p} %{buildroot}%{_bindir}
%{__mkdir_p} %{buildroot}%{_sysconfdir}/%{cmdname}

cp -R ansible %{buildroot}%{ansible_dir}
cp -R ansible/roles/ %{buildroot}%{cme_premium_roles_dir}

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
