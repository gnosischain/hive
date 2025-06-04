#!/usr/bin/env bash
# compare_dirs.sh DIR1 DIR2 – counts identical / changed / unique files
# and prints the lists.

set -euo pipefail

[ $# -eq 2 ] || { echo "Usage: $0 DIR1 DIR2" >&2; exit 1; }
d1="${1%/}"
d2="${2%/}"

# Run diff once, capture output
out=$(diff -r -q -s "$d1" "$d2" || true)   # diff exits 1 when differences exist

identical=$(printf '%s\n' "$out" | grep -c " are identical$")
changed=$( printf '%s\n' "$out" | grep -c " differ$")
unique=$(  printf '%s\n' "$out" | grep -c "^Only in ")

printf 'Identical files : %s\n' "$identical"
printf 'Changed  files  : %s\n' "$changed"
printf 'Unique   files  : %s\n\n' "$unique"

printf '=== Changed files ===\n'
printf '%s\n' "$out" | grep " differ$"  || true
printf '\n=== Unique files (exist only in one tree) ===\n'
printf '%s\n' "$out" | grep "^Only in " || true

