#!/bin/bash
set -e -x

pushd /tmp

# Only run the backup-merge test
i="tests/backup_merge_test.sh"
echo
echo "===== RUNNING $i ====="
set -x

./"$i";

set +x
echo "===== SUCCESS $i ====="
echo

popd
