#!/bin/sh
set -e

log_run() {
	echo "running $1"
}

log_ok() {
	echo "✔ $1"
}

run_step() {
	label="$1"
	shift

	logfile="$tmpdir/pre-commit-step.log"
	: >"$logfile"

	echo "running $label"

	if "$@" >"$logfile" 2>&1; then
		log_ok "$label"
	else
		echo "✖ $label"
		if [ -s "$logfile" ]; then
			cat "$logfile"
		fi
		exit 1
	fi
}

# ----------------------------
# staged files
# ----------------------------
tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

root="$(git rev-parse --show-toplevel)"
cd "$root"

git diff --cached --name-only --diff-filter=ACM >"$tmpdir/list"
grep '^backend/.*\.go$' "$tmpdir/list" >"$tmpdir/gofiles" 2>/dev/null || true

# ----------------------------
# gofmt (check mode)
# ----------------------------
if [ "${PRE_COMMIT_GOFMT_CHECK:-}" = "1" ]; then
	log_run "gofmt (check mode)"

	if [ ! -s "$tmpdir/gofiles" ]; then
		echo "- no staged backend .go files, skip gofmt"
	else
		bad=""
		while IFS= read -r f; do
			[ -z "$f" ] && continue
			[ -f "$f" ] || continue

			if gofmt -l "$f" | grep -q .; then
				bad="$bad $f"
			fi
		done <"$tmpdir/gofiles"

		if [ -n "$bad" ]; then
			echo "gofmt would change these files - run: make -C backend fmt"
			exit 1
		fi
	fi

	log_ok "gofmt"

	run_step "go vet" make -C backend vet
	exit 0
fi

# ----------------------------
# gofmt (auto fix)
# ----------------------------
log_run "gofmt"

if [ ! -s "$tmpdir/gofiles" ]; then
	echo "- no staged backend .go files, skip gofmt"
else
	while IFS= read -r f; do
		[ -z "$f" ] && continue
		[ -f "$f" ] || continue

		gofmt -w "$f"
		git add "$f"
	done <"$tmpdir/gofiles"
fi

log_ok "gofmt"

# ----------------------------
# vet
# ----------------------------
run_step "go vet" make -C backend vet
