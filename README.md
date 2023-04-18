# Main & Develop conventional release action

This is very simple but efficient, push into `develop` branch then you will have a patch release, push into `main` branch then you can have `minor` & `major` release.

## How to setup?

use the example below or the [release.yml](./.github/workflows/release.yml) file inside `.github` directory.

```yml
name: Release

on:
  push:
    branches:
      - main
      - develop

jobs:
  next-version:
    permissions: write-all
    runs-on: ubuntu-latest
    timeout-minutes: 5
    steps:
      - name: Checkout code
        uses: actions/checkout@8e5e7e5ab8b370d6c329ec480221332ada57f0ab
        with:
          fetch-depth: 0

      - name: Release
        uses: mehdi-ra/main-develop-semver@a72752066126879a5ca505f0d0a733ed9e9602e1
        with:
          token: ${{secrets.GITHUB_TOKEN}}
          releaseTitle: Auto release
```

## Inputs

| name         | required | description                     |
| ------------ | -------- | ------------------------------- |
| token        | true     | Token for making release        |
| releaseTitle | false    | Title of releases               |
| repository   | false    | where to get latest tag version |

## Outputs

| name           | description              |
| -------------- | ------------------------ |
| releaseVersion | Token for making release |

## Version fall back

If the action does not detect any change then creates a release with tag `v0.1.0`
