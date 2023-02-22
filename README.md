<p align="center"><a href="#readme"><img src="https://gh.kaos.st/sonar.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/sonar/ci"><img src="https://kaos.sh/w/sonar/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/sonar"><img src="https://kaos.sh/r/sonar.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/sonar"><img src="https://kaos.sh/b/49715c23-4ead-4edb-a351-b4c49cf8d061.svg" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/w/sonar/codeql"><img src="https://kaos.sh/w/sonar/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
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
sudo yum install -y https://yum.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install sonar
```

#### Slack app

Create [classic Slack App](https://api.slack.com/apps?new_classic_app=1) (_this application requires RTM and can't work with socket mode due to lack of presence events_).

Add next OAuth scopes to your app:

- `dnd:read`
- `users:read`

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
| `master` | [![CI](https://kaos.sh/w/sonar/ci.svg?branch=master)](https://kaos.sh/w/sonar/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/sonar/ci.svg?branch=master)](https://kaos.sh/w/sonar/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
