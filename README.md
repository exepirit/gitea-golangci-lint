# gitea-golangci-lint

![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/exepir1t/gitea-golangci-lint)

This tool can send golangci-lint issues to Gitea as pull request reviews. You can visit <https://golangci-lint.run/> to
find the golangci-lint configurations.

## Docker Image
To check the docker image, visit <https://hub.docker.com/r/exepir1t/gitea-golangci-lint>
> docker pull exepir1t/gitea-golangci-lint

## Build

```shell
CGO_ENABLED=0 go build -o gitea-golangci-lint
```

## Configuration

There are 6 environment variables to be configured when you want to run this tool.

| Variable | Description | Example |
| --- | --- | --- |
| `GITEA_URL` | Gitea server url | `https://try.gitea.io` |
| `GITEA_USER` | Gitea username | `golanglinter` |
| `GITEA_TOKEN` | Gitea access token | `your_gitea_user_accesstoken` |
| `GITEA_REPO` | Repository name, which is inspected | `octocat/hello_world` |
| `PULL_REQUEST` | Pull request ID | `123` |
| `HTTP_TIMEOUT` | HTTP requests timeout in seconds | `30` |

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
    temp: { }

steps:
  - name: linter
    image: golangci/golangci-lint:latest-alpine
    pull: if-not-exists
    volumes:
      - name: deps
        path: /go
    commands:
      - golangci-lint run | tee .golangci-lint.log
      - |
        [[ -z "$${DRONE_PULL_REQUEST}" ]] && [[ -s .golangci-lint.log ]] && exit 1
      - exit 0
  - name: push linter review
    image: exepir1t/gitea-golangci-lint:latest
    pull: if-not-exists
    environment:
      GITEA_URL:
        from_secret: GITEA_URL
      GITEA_USER:
        from_secret: GITEA_CI_USER
      GITEA_TOKEN:
        from_secret: GITEA_CI_TOKEN
    commands:
      - cat .golangci-lint.log | /bin/gitea-golangci-lint
    when:
      event:
        - pull_request
```

Make the `pipeline` in the `.drone.yml` to be run correctly，you should add `secrets` as below:

1. `GITEA_URL`: origin of Gitea server，Like: `https://git.example.com`;
2. `GITEA_CI_USER`: user of Gitea server，Like: `gitea`, which has `read` authorization to the
   repo;
3. `GITEA_CI_TOKEN`: token of Gitea user, which has `read` authorization to the repo;
