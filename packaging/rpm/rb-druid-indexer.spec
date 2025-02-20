Name:      rb-druid-indexer
Version:   %{__version}
Release:   %{__release}%{?dist}
BuildArch: x86_64
Summary:   rb-druid-indexer service to create indexing tasks in druid supervisors

License:   AGPL-3.0
URL:       https://github.com/redBorder/rb-druid-indexer
Source0:   %{name}-%{version}.tar.gz

Requires: pmacct arpwatch rsyslog
BuildRequires: golang

%global debug_package %{nil}

%description
%{summary}

%prep
%setup -q

go mod tidy

%build
go build -o bin/rb-druid-indexer ./main.go

%install
install -D -m 0755 bin/rb-druid-indexer %{buildroot}/usr/bin/rb-druid-indexer
install -D -m 0644 packaging/rpm/rb-druid-indexer.service %{buildroot}/usr/lib/systemd/system/rb-druid-indexer.service

%pre
getent group rb-druid-indexer >/dev/null || groupadd -r rb-druid-indexer
getent passwd rb-druid-indexer >/dev/null || useradd -r -g rb-druid-indexer -d /var/lib/rb-druid-indexer -s /sbin/nologin -c "rb-druid-indexer user" rb-druid-indexer

%post
systemctl daemon-reload

%files
%defattr(755,root,root)
/usr/bin/rb-druid-indexer
%defattr(644,root,root)
/usr/lib/systemd/system/rb-druid-indexer.service

%doc

%changelog
* Thu Feb 20 2025 Miguel Álvarez <malvarez@redborder.com>
- First version of rb-druid-indexer (Go-based)
