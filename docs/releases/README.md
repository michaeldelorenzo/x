# Releases
This repo follows [Semver](https://semver.org/) and [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to determine how to version the codebase. Github Actions are used to publish a release when code is pushed to certain branches.

Releases will generate two assets: 
1) A zipped Golang binary asset on the release
2) A docker image in ECR with the same version tag as the release

A prerelease can be triggered on any branch by [manually running](https://docs.github.com/en/free-pro-team@latest/actions/managing-workflow-runs/manually-running-a-workflow) the `release` workflow on a given branch.

A release will be triggered on any merge into `production`.

| Conventional Commit Type | Description                                |
|--------------------------|--------------------------------------------|
| `BREAKING-CHANGE`        | Bump the project's `major` version number. |
| `feat`                   | Bump the project's `minor` version number. |
| `fix`                    | Bump the project's `patch` version number. |
| `build`                  | Bump the project's `patch` version number. |
| `ci`                     | Bump the project's `patch` version number. |
| `refactor`               | Bump the project's `patch` version number. |
| `style`                  | Bump the project's `patch` version number. |
| `perf`                   | Bump the project's `patch` version number. |
| `docs`                   | No version number change.                  |
| `test`                   | No version number change.                  |

> NOTE: BREAKING-CHANGE is not a commit type on its own & must be used with another type. See below for example usage:

> refactor(PROJ-2032)!: reorganize api routes
>
> \- reorganize api routes to make them more restful
>
> BREAKING-CHANGE: API has been refactored to V2. V1 API is no longer supported and clients must update to continue using this service.