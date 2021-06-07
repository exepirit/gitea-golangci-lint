# GITEA-GOLANGCI-LINT

This lib can send golangci-lint issues to gitea as pull request reviews.
You can visit <https://golangci-lint.run/> to find the golangci-lint configurations.

## Build
```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o gitea-golangci-lint
```

## Configurations
There are 6 configurations to be configured when you want to run this tool.

### Gitea Server Url
```shell
export GITEA_URL=https://git.giteaserver.com
```
> There has not one slash(/) at the end of the url.

### Gitea Username
```shell
export GITEA_USER=golanglinter
```

### Gitea AccessToken
```shell
export GITEA_TOKEN=your_gitea_user_accesstoken
```

### Git REPO
```shell
export GIT_REPO=octocat/hello-worId
```

### Pull Request ID
```shell
export PULL_REQUEST=123
```

### Http Timeout
```shell
export HTTP_TIMEOUT=30
```

## How to use?
### Run it in the shell
```shell
golangci-lint run | gitea-golangci-lint
```

### Run with Drone docker pipeline
Below is a drone task configuration, it may help you to config your drone task.

```yaml
---
kind: pipeline
type: docker
name: default

volumes:
- name: deps
  temp: {}

steps:
- name: linter
  image: golangci/golangci-lint:latest-alpine
  pull: if-not-exists
  environment:
    GOPROXY:
      from_secret: GOPROXY
  volumes:
    - name: deps
      path: /go
  commands:
    - golangci-lint run | tee .golangci-lint.log
    - |
      [[ -z "$${DRONE_PULL_REQUEST}" ]] && [[ -s .golangci-lint.log ]] && exit 1
    - exit 0
- name: review
  image: newbing/gitea-golangci-lint:latest
  pull: if-not-exists
  environment:
    GITEA_URL:
      from_secret: GITEA_URL
    GITEA_USER:
      from_secret: GITEA_CI_USER
    GITEA_TOKEN:
      from_secret: GITEA_CI_TOKEN
  commands:
    - export GIT_REPO=$DRONE_REPO
    - export PULL_REQUEST=$DRONE_PULL_REQUEST
    - cat .golangci-lint.log | gitea
  when:
    event:
    - pull_request
```

#### `DroneCI` 配置 `secrets` 变量说明

> 此仓库中 `.drone.yml` 文件配置的 `pipeline` 正确运行，需要配置以下 `secrets` 变量。

1. `GITEA_URL`, `GITEA_URL` 是`Gitea`仓库的地址，例如:`https://git.example.com`, 末尾不需要`/`; (`review step` 使用)
2. `GITEA_CI_USER`, `GITEA_CI_USER` 是`Gitea`仓库的用户名，例如:`gitea`, 需要有对应仓库的 `read` 权限; (`review step` 使用)
3. `GITEA_CI_TOKEN`, `GITEA_CI_TOKEN` 是`Gitea`仓库用户的访问 `Token`, 需要有对应仓库的 `read` 权限; (`review step` 使用)