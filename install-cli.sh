#!/bin/sh

set -eu

INSTALLROOT=${INSTALLROOT:-"${HOME}/.tobs"}
TOBS_VERSION=${TOBS_VERSION:-0.1.0-beta.1}

happyexit() {
  echo ""
  echo "Add the tobs CLI to your path with:"
  echo ""
  echo "  export PATH=\$PATH:${INSTALLROOT}/bin"
  echo ""
  echo "After starting your cluster, run"
  echo ""
  echo "  tobs install"
  echo ""
  exit 0
}

validate_checksum() {
  filename=$1
  checksumlist=$(curl -sfL "${url}/checksums.txt")
  echo ""
  echo "Validating checksum..."

  checksum=$($checksumbin -a256 "${filename}")

  if grep -Fxq "${checksum}" <<< "${checksumlist}"; then
    echo "Checksum valid."
    return 0
  else
    echo "Checksum validation failed." >&2
    return 1
  fi
}

OS=$(uname -s)
arch=$(uname -m)
case $OS in
  Darwin)
    ;;
  Linux)
    case $arch in
      x86_64)
        ;;
      i386)
        ;;
      *)
        echo "The Observability Stack does not support $OS/$arch. Please open an issue with your platform details."
        exit 1
        ;;
    esac
    ;;
  *)
    echo "The Observability Stack does not support $OS/$arch. Please open an issue with your platform details."
    exit 1
    ;;
esac

tarbin=$(command -v tar) || {
  echo "Failed to find unpacking binary. Please install tar."
  exit 1
}

checksumbin=$(command -v shasum) || {
  echo "Failed to find checksum binary. Please install shasum."
  exit 1
}

tmpdir=$(mktemp -d /tmp/tobs.XXXXXX)
srcfile="tobs_${TOBS_VERSION}_${OS}_${arch}.tar.gz"
dstfile="${INSTALLROOT}/bin/tobs-${TOBS_VERSION}"
url="https://github.com/timescale/timescale-observability/releases/download/${TOBS_VERSION}"

(
  cd "$tmpdir"

  echo "Downloading ${srcfile}..."
  curl -fLO "${url}/${srcfile}"
  echo "Download complete!"

  if ! validate_checksum "${srcfile}"; then
    exit 1
  fi

  $tarbin -xvf $srcfile

  echo ""
)

(
  mkdir -p "${INSTALLROOT}/bin"
  mv "${tmpdir}/tobs" "${dstfile}"
  chmod +x "${dstfile}"
  rm -f "${INSTALLROOT}/bin/tobs"
  ln -s "${dstfile}" "${INSTALLROOT}/bin/tobs"
)

rm -r "$tmpdir"
echo "tobs ${TOBS_VERSION} was successfully installed 🎉"
echo ""
happyexit
