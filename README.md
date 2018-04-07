# dgx-price-feeder

## Prerequisites

1. Docker and Docker compose
2. A json keystore file
3. A passphrase file (the file should contain only the passphrase, note that the app will remove leading spaces, trailing spaces and new lines so please make sure your passphrase doesn't have such characters)

## Install

1. Copy the keystore to `<repo_root>/cmd/keystore`
2. Copy the passphrase to `<repo_root>/cmd/passphrase`
3. Assume you are in `<repo_root>`, run `docker-compose build` to build the image
4. `docker-compose run -d` to run the price feeder

Where: `<repo_root>` is the path to this repo.
Note: `cmd/keystore` and `cmd/passphrase` are ignored by the `.gitignore` to avoid mistakenly committing the credentical to git.

## Log

The log will be written to `<repo_root>/log` and will be rotated daily.
