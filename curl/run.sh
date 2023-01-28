#!env sh

set -e

curl https://peapix.com/bing/feed | jq -r '.[].imageUrl' | xargs curl --remote-name-all
