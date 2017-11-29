# Cigger - command line client to trigger CI service

[![Build Status](https://travis-ci.org/otiai10/cigger.svg?branch=master)](https://travis-ci.org/otiai10/cigger)

# Usage

```sh
% cigger -s travis -p your/project
```

# Installation

```sh
% go get -u github.com/otiai10/cigger
```

For more information, hit `cigger --help`

# Usecase

Trigger another CI build after one, like this [.travis.yml](https://github.com/otiai10/cwl.go/blob/master/.travis.yml#L10-L11)

```yaml
language: go
go:
  - 1.8
before_script:
  # Get cigger beforehands
  - go get github.com/otiai10/cigger
script:
  # Run the tests of your project
  - go test -v
after_script:
  # Trigger build of another project via `cigger`
  - cigger -s travis -p otiai10/yacle -t ${TRAVIS_API_TOKEN}
```
