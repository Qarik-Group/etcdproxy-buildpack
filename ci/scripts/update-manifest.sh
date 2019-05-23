#!/bin/bash

set -eu

version=$(cat etcd/version)
sha=$(sha256sum etcd/etcd-v*-linux-amd64.tar.gz | awk '{print $1}')

git clone git pushme

cat > pushme/manifest.yml <<YAML
---
language: etcdproxy
default_versions:
- name: etcd
  version: ${version}
dependency_deprecation_dates: []
dependencies:
- name: etcd
  version: ${version}
  uri: https://github.com/etcd-io/etcd/releases/download/v${version}/etcd-v${version}-linux-amd64.tar.gz
  sha256: ${sha}
  cf_stacks:
  - cflinuxfs2
  - cflinuxfs3

include_files:
  - README.md
  - VERSION
  - bin/supply
  - manifest.yml
pre_package: scripts/build.sh
YAML

pushd pushme

if [[ "$(git status -s)X" != "X" ]]; then
  set +e
  if [[ -z $(git config --global user.email) ]]; then
    git config --global user.email "drnic+bot@starkandwayne.com"
  fi
  if [[ -z $(git config --global user.name) ]]; then
    git config --global user.name "CI Bot"
  fi

  set -e
  echo ">> Running git operations as $(git config --global user.name) <$(git config --global user.email)>"
  echo ">> Getting back to master (from detached-head)"
  git merge --no-edit master
  git status
  git --no-pager diff
  git add manifest.yml
  git commit -m "Updated etcd to v${version}"
else
  echo ">> No update needed"
fi

popd
