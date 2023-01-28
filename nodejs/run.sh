#!env sh

set -e

yarn
yarn eslint main.ts
yarn tsc
node main.js
