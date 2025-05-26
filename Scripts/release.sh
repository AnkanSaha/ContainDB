#!/bin/bash

set -e

# === CONFIG ===
APP_NAME="containdb"
ARCH="amd64"
VERSION_FILE="./VERSION"
VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
DEB_FILE="./Debian/${APP_NAME}_${VERSION}_${ARCH}.deb"
TAG="v$VERSION"
COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_MSG=$(git log -1 --pretty=%B)

# === Environment Variables ===
REPO="${GIT_REPOSITORY}"  # GitHub Actions sets this automatically
TOKEN="${GIT_TOKEN}"      # GitHub Actions provides this

# === Build steps ===
./Scripts/BinBuilder.sh
./Scripts/DebBuilder.sh

# === Checks ===
if [ ! -f "$VERSION_FILE" ]; then
  echo "❌ VERSION file not found"
  exit 1
fi

if [ ! -f "$DEB_FILE" ]; then
  echo "❌ .deb file not found at $DEB_FILE"
  exit 1
fi

if ! command -v gh &> /dev/null; then
  echo "❌ GitHub CLI (gh) not installed"

  # Update package list and install dependencies
  apt update
  apt install -y curl gnupg software-properties-common

  # Add GitHub CLI's official package repository
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | \
    dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg

  sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | \
    tee /etc/apt/sources.list.d/github-cli.list > /dev/null

  # Install gh
  apt update
  apt install -y gh


fi

# === Authenticate GitHub CLI ===
echo "🔑 Authenticating GitHub CLI..."
echo "${TOKEN}" | gh auth login --with-token

# === Create Release ===
echo "📦 Creating GitHub release for tag $TAG..."

gh release create "$TAG" "$DEB_FILE" \
  --title "$TAG" \
  --notes "🔨 Commit: $COMMIT_HASH

📝 Message:
$COMMIT_MSG"

echo "✅ GitHub release published with .deb asset"
