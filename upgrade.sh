#!/bin/bash
set -e
BASE_DIR=$(pwd)

# 根目录执行
echo "=== Tidy root directory ==="
go mod tidy

# 递归查找 1~3 层内所有 go.mod，自动执行 go mod tidy
echo -e "\n=== Tidy all submodules (1~3 levels) ==="
find . -maxdepth 3 -type f -name "go.mod" | grep -v '^\./go.mod$' | sort | while read -r modfile; do
    moddir=$(dirname "${modfile}")
    echo "====================================="
    echo "Processing: ${moddir}"
    cd "${moddir}"
    go mod tidy
    cd "${BASE_DIR}"
done

echo -e "\n====================================="
echo -e "✅ All modules tidy finished!"