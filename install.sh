#!/usr/bin/env sh
# Installs the latest abacatepay CLI release from GitHub Releases.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/AbacatePay/abacatepay-cli/main/install.sh | sh
#
# Env:
#   ABACATEPAY_INSTALL_DIR   Install directory (default: /usr/local/bin)
set -eu

REPO="AbacatePay/abacatepay-cli"
BIN_NAME="abacatepay"
INSTALL_DIR="${ABACATEPAY_INSTALL_DIR:-/usr/local/bin}"

detect_os() {
  uname -s | tr '[:upper:]' '[:lower:]'
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64) echo amd64 ;;
    arm64 | aarch64) echo arm64 ;;
    *)
      echo "Unsupported architecture: $(uname -m)" >&2
      exit 1
      ;;
  esac
}

OS="$(detect_os)"
ARCH="$(detect_arch)"

case "$OS" in
  darwin | linux) ;;
  *)
    echo "This script supports macOS and Linux only. On Windows, download a release archive manually:" >&2
    echo "  https://github.com/${REPO}/releases/latest" >&2
    exit 1
    ;;
esac

ASSET="${BIN_NAME}_${OS}_${ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/latest/download"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${ASSET}..."
curl -fsSL "${BASE_URL}/${ASSET}" -o "${TMP_DIR}/${ASSET}"

if curl -fsSL "${BASE_URL}/checksums.txt" -o "${TMP_DIR}/checksums.txt" 2>/dev/null; then
  EXPECTED="$(grep " ${ASSET}\$" "${TMP_DIR}/checksums.txt" | awk '{print $1}')"
  if [ -n "${EXPECTED}" ]; then
    if command -v sha256sum >/dev/null 2>&1; then
      ACTUAL="$(sha256sum "${TMP_DIR}/${ASSET}" | awk '{print $1}')"
    else
      ACTUAL="$(shasum -a 256 "${TMP_DIR}/${ASSET}" | awk '{print $1}')"
    fi
    if [ "${EXPECTED}" != "${ACTUAL}" ]; then
      echo "Checksum mismatch for ${ASSET} (expected ${EXPECTED}, got ${ACTUAL})" >&2
      exit 1
    fi
    echo "Checksum verified."
  fi
else
  echo "Warning: could not fetch checksums.txt, skipping verification" >&2
fi

tar -xzf "${TMP_DIR}/${ASSET}" -C "${TMP_DIR}" "${BIN_NAME}"

if [ -w "${INSTALL_DIR}" ]; then
  mv "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
else
  echo "Installing to ${INSTALL_DIR} requires sudo..."
  sudo mv "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
fi
chmod +x "${INSTALL_DIR}/${BIN_NAME}"

echo "Installed ${BIN_NAME} to ${INSTALL_DIR}"
"${INSTALL_DIR}/${BIN_NAME}" --version 2>/dev/null || true
