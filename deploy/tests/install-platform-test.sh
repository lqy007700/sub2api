#!/bin/bash

set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
# Load installer functions without performing an installation.
source <(sed '$d' "$ROOT_DIR/deploy/install.sh")
LANG_CHOICE="en"

fail() {
    printf 'installer platform test failed: %s\n' "$*" >&2
    exit 1
}

uname() {
    case "$1" in
        -s) printf '%s\n' "$TEST_OS" ;;
        -m) printf '%s\n' "$TEST_ARCH" ;;
        *) fail "unexpected uname argument: $1" ;;
    esac
}

for TEST_ARCH in x86_64 amd64; do
    TEST_OS=Linux
    detect_platform >/dev/null
    [[ "$OS" == "linux" && "$ARCH" == "amd64" ]] || fail "incorrect asset platform for $TEST_ARCH"
done

# Any of these calls would mean an unsupported host got past the early guard.
select_language() { fail "prompted for language on unsupported host"; }
check_root() { fail "checked privileges on unsupported host"; }
check_dependencies() { fail "checked dependencies on unsupported host"; }
github_api_curl() { fail "contacted GitHub on unsupported host"; }

assert_rejected_early() {
    TEST_OS=$1
    TEST_ARCH=$2
    local output
    shift 2
    if output=$(main "$@" 2>&1); then
        fail "accepted $TEST_OS/$TEST_ARCH"
    fi
    [[ "$output" == *"requires Linux amd64/x86_64"* ]] || fail "missing platform hint for $TEST_OS/$TEST_ARCH: $output"
    [[ "$output" != *"installer platform test failed:"* ]] || fail "$output"
}

for command in install upgrade rollback list-versions uninstall; do
    assert_rejected_early Darwin arm64 "$command"
    assert_rejected_early Darwin x86_64 "$command"
    assert_rejected_early Linux aarch64 "$command"
    assert_rejected_early Linux arm64 "$command"
    assert_rejected_early Linux i686 "$command"
    assert_rejected_early MINGW64_NT x86_64 "$command"
done
# The no-argument default path must reject unsupported hosts too.
assert_rejected_early Linux riscv64

printf 'installer Linux amd64 platform checks passed\n'
