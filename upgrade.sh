#!/bin/bash
set -e
BASE_DIR=$(pwd)
VERSION="v1.0.0"

# 根目录执行
echo "=== Tidy root directory ==="
go mod tidy

# 遍历所有一级子目录
echo -e "\n=== Tidy submodules ==="
for dir in */; do
    if [ -f "${dir}go.mod" ]; then
        mod_name=$(basename "${dir}")
        echo "Processing: ${mod_name}"
        cd "${dir}" && go mod tidy
        cd "${BASE_DIR}"
    fi
done

echo -e "\n✅ All modules tidy finished!"