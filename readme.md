<p align="center">
  <a href="https://travis-ci.org/essentialkaos/sonar"><img src="https://travis-ci.org/essentialkaos/sonar.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/sonar"><img src="https://goreportcard.com/badge/github.com/essentialkaos/sonar" alt="GoReportCard" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.io/ekol.svg" alt="License" /></a>
</p>

<p align="center"><a href="#readme"><img src="https://gh.kaos.io/sonar.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

`Sonar` is a utility for showing user Slack status in Atlassian Jira.

### Installation

#### From ESSENTIAL KAOS Public repo for RHEL6/CentOS6

```bash
[sudo] yum install -y https://yum.kaos.io/6/release/x86_64/kaos-repo-8.0-0.el6.noarch.rpm
[sudo] yum install bastion
```

#### From ESSENTIAL KAOS Public repo for RHEL7/CentOS7

```bash
[sudo] yum install -y https://yum.kaos.io/7/release/x86_64/kaos-repo-8.0-0.el7.noarch.rpm
[sudo] yum install bastion
```

#### Integration with Jira

Go to `atlassian-jira/WEB-INF/classes/templates/plugins/userformat` and modify next files.

##### `actionProfileLink.vm`

```
<a $!{userHoverAttributes} id="$!{id}" href="${baseurl}/secure/ViewProfile.jspa?name=${velocityhelper.urlencode($username)}">${renderedAvatarImg} ${author}</a><img class="slack-status" src="https://sonar/status.svg?user=$username" />
```

##### `profileLinkWithAvatar.vm`

```
${textutils.htmlEncode($fullname)}<img class="slack-status" src="https://sonar/status.svg?user=$username" />
```

##### `avatarFullNameHover.vm`

```
$textutils.htmlEncode($fullname)
<img class="slack-status" src="https://sonar/status.svg?user=$username" />
```

Then restart your Jira instance.

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.org/essentialkaos/bastion.svg?branch=master)](https://travis-ci.org/essentialkaos/bastion) |
| `develop` | [![Build Status](https://travis-ci.org/essentialkaos/bastion.svg?branch=develop)](https://travis-ci.org/essentialkaos/bastion) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.io/ekgh.svg"/></a></p>
