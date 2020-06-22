<p align="center"><a href="#readme"><img src="https://gh.kaos.st/sonar.svg"/></a></p>

<p align="center">
  <a href="https://travis-ci.com/essentialkaos/sonar"><img src="https://travis-ci.com/essentialkaos/sonar.svg?branch=master" alt="TravisCI" /></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-sonar-master"><img alt="codebeat badge" src="https://codebeat.co/badges/49715c23-4ead-4edb-a351-b4c49cf8d061" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/sonar"><img src="https://goreportcard.com/badge/github.com/essentialkaos/sonar" alt="GoReportCard" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg" alt="License" /></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`Sonar` is a utility for showing user Slack status in Atlassian Jira.

### Screenshots

<p align="center">
  <img src="https://gh.kaos.st/sonar-preview.png" alt="Sonar Preview">
  <i>Sonar in Jira 6.x (with <a href="https://github.com/essentialkaos/atlassian-remixed-theme">Remixed Theme</a>)</i>
</p>

### Installation

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install sonar
```

#### Integration with Jira

Go to `atlassian-jira/WEB-INF/classes/templates/plugins/userformat` and modify next files.

**actionProfileLink.vm**

```html
<a $!{userHoverAttributes} id="$!{id}" href="${baseurl}/secure/ViewProfile.jspa?name=${velocityhelper.urlencode($username)}">${renderedAvatarImg} ${author}</a><img class="slack-status" src="https://sonar.domain.com/status.svg?token=YOUR_TOKEN_HERE&mail=$user.emailAddress" />
```

**profileLinkWithAvatar.vm**

```html
${textutils.htmlEncode($fullname)}<img class="slack-status" src="https://sonar.domain.com/status.svg?token=YOUR_TOKEN_HERE&mail=$user.emailAddress" />
```

**avatarFullNameHover.vm**

```html
$textutils.htmlEncode($fullname)
<img class="slack-status" src="https://sonar.domain.com/status.svg?token=YOUR_TOKEN_HERE&mail=$user.emailAddress" />
```

Then restart your Jira instance.

Also, you can add `sonar.js` to your announcement banner for a periodic status update.

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![Build Status](https://travis-ci.com/essentialkaos/sonar.svg?branch=master)](https://travis-ci.com/essentialkaos/sonar) |
| `develop` | [![Build Status](https://travis-ci.com/essentialkaos/sonar.svg?branch=develop)](https://travis-ci.com/essentialkaos/sonar) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
