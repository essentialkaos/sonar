################################################################################

%global crc_check pushd ../SOURCES ; sha512sum -c %{SOURCE100} ; popd

################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Utility for showing user Slack status in Atlassian Jira
Name:           sonar
Version:        1.8.2
Release:        0%{?dist}
Group:          Applications/System
License:        Apache License, Version 2.0
URL:            https://kaos.sh/sonar

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

Source100:      checksum.sha512

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  golang >= 1.21

Requires:       kaosv >= 2.16
Requires:       systemd

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Utility for showing user Slack status in Atlassian Jira.

################################################################################

%prep
%{crc_check}

%setup -q

%build
if [[ ! -d "%{name}/vendor" ]] ; then
  echo "This package requires vendored dependencies"
  exit 1
fi

pushd %{name}
  %{__make} %{?_smp_mflags} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dDm 755 %{buildroot}%{_bindir}
install -dDm 755 %{buildroot}%{_sysconfdir}/logrotate.d
install -dDm 755 %{buildroot}%{_initddir}
install -dDm 755 %{buildroot}%{_unitdir}
install -dDm 755 %{buildroot}%{_localstatedir}/log/%{name}

install -pm 755 %{name}/%{name} \
                %{buildroot}%{_bindir}/

install -pm 644 %{name}/common/%{name}.knf \
                %{buildroot}%{_sysconfdir}/

install -pm 755 %{name}/common/%{name}.init \
                %{buildroot}%{_initddir}/%{name}

install -pm 644 %{name}/common/%{name}.logrotate \
                %{buildroot}%{_sysconfdir}/logrotate.d/%{name}

install -pDm 644 %{name}/common/%{name}.service \
                 %{buildroot}%{_unitdir}/

%clean
rm -rf %{buildroot}

%pre
getent group %{name} >/dev/null || groupadd -r %{name}
getent passwd %{name} >/dev/null || useradd -r -M -g %{name} -s /sbin/nologin %{name}
exit 0

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%attr(-,%{name},%{name}) %dir %{_localstatedir}/log/%{name}
%config(noreplace) %{_sysconfdir}/%{name}.knf
%config(noreplace) %{_sysconfdir}/logrotate.d/%{name}
%{_unitdir}/%{name}.service
%{_initddir}/%{name}
%{_bindir}/%{name}

################################################################################

%changelog
* Fri Mar 29 2024 Anton Novojilov <andy@essentialkaos.com> - 1.8.2-0
- Improved support information gathering
- Code refactoring
- Dependencies update

* Fri Apr 28 2023 Anton Novojilov <andy@essentialkaos.com> - 1.8.1-0
- Dependencies update
- Typos update

* Thu Mar 31 2022 Anton Novojilov <andy@essentialkaos.com> - 1.8.0-0
- Removed pkg.re usage
- Smaller, optimized SVG images
- Fixed major bug with checking token

* Thu Sep 10 2020 Anton Novojilov <andy@essentialkaos.com> - 1.7.1-0
- Fixed compatibility with the latest version of ek package

* Sat May 23 2020 Anton Novojilov <andy@essentialkaos.com> - 1.7.0-0
- Fixed problem with initial fetching of DND statuses
- nlopes/slack package replaced by slack-go/slack
- fasthttp updated to the latest version
- ek package updated to the latest stable version
- Minor improvements

* Sat Jun 29 2019 Anton Novojilov <andy@essentialkaos.com> - 1.6.1-0
- fasthttp updated to the latest version
- ek package updated to the latest stable version

* Mon Mar 25 2019 Anton Novojilov <andy@essentialkaos.com> - 1.6.0-0
- Improved mechanics of checking users presence
- Added support of debug slack logging
- Added icon for DND + offline status
- slack package switched to original version
- ek package updated to the latest stable version
- fasthttp package updated to the latest stable release

* Wed Dec 19 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.4-0
- ek package updated to v10
- Improved nginx config example

* Thu Oct 25 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.3-0
- fasthttp package updated to the latest release

* Wed Oct 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.2-0
- fasthttp package updated to the latest stable release (1.0.0)

* Wed Jun 27 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.1-1
- Rebuilt with the latest version of the slack package

* Thu Jun 21 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.1-0
- Return transparent image instead of white bullet if user is unknown

* Mon May 28 2018 Anton Novojilov <andy@essentialkaos.com> - 1.5.0-0
- Added systemd unit
- Improved SysV init script
- Rebuilt with the latest version of the slack package

* Wed Mar 28 2018 Anton Novojilov <andy@essentialkaos.com> - 1.4.0-0
- fasthttp package replaced by erikdubbelboer fork
- slack package updated to v3
- Added open files limits to init script
- Added configuration file for log rotation

* Thu Jan 18 2018 Anton Novojilov <andy@essentialkaos.com> - 1.3.1-0
- Fixed subscribing for presence events when new user was added

* Thu Jan 18 2018 Anton Novojilov <andy@essentialkaos.com> - 1.3.0-0
- Fixed compatibility with latest version of Slack API

* Wed Jan 17 2018 Anton Novojilov <andy@essentialkaos.com> - 1.2.1-0
- Updated path to slack package

* Mon Jan 15 2018 Anton Novojilov <andy@essentialkaos.com> - 1.2.0-1
- Rebuilt with latest version of slack package

* Sat Dec 09 2017 Anton Novojilov <andy@essentialkaos.com> - 1.2.0-0
- Added on-call status handling
- New icons
- Improved SVG icons
- Minor performance improvements

* Mon Nov 20 2017 Anton Novojilov <andy@essentialkaos.com> - 1.1.2-0
- Minor fix for status point SVG generation

* Mon Nov 20 2017 Anton Novojilov <andy@essentialkaos.com> - 1.1.1-0
- Minor fix for status point SVG generation

* Mon Nov 20 2017 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Colorblind-friendly design and colors
- slack package replaced by ek fork
- Fixed bug with updating user DND status

* Sat Oct 14 2017 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Initial public release
