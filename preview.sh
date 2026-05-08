#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

cd "$(dirname "$0")"

if ! command -v gem &>/dev/null; then
    echo "Error: 'gem' not found. Install Ruby first." >&2
    exit 1
fi

root="$(pwd)"
export GEM_HOME="$root/vendor/gems"
export GEM_PATH="$root/vendor/gems"
export PATH="$root/vendor/gems/bin:$PATH"

if ! command -v bundle &>/dev/null; then
    echo "Installing bundler..."
    gem install bundler
fi

bundle config set --local path vendor/bundle

if ! bundle check >/dev/null 2>&1; then
    echo "Installing gems..."
    bundle install
fi

exec bundle exec jekyll serve
