#!/bin/bash
set -e -x

/home/gpadmin/run_greenplum.sh

pushd /tmp
# Only run backup_merge_test.sh for efficiency
test_script="tests/backup_merge_test.sh"
if [ -f "$test_script" ]; then
  echo
  echo "=============================="
  echo "===== RUNNING $test_script ====="
  echo "=============================="
  set -x
  ./"$test_script";

  set +x
  echo "=============================="
  echo "===== SUCCESS $test_script ====="
  echo "=============================="
  echo
else
  echo "Test script $test_script not found!"
  exit 1
fi
popd
