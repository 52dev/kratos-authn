#!/bin/bash
set -e
BASE_DIR=$(pwd)
VERSION="v1.0.0"

# 强制创建/覆盖根标签
echo "=== Force create root tag: ${VERSION} ==="
git tag -f "${VERSION}"

# 强制创建/覆盖子模块标签
echo -e "\n=== Force create submodule tags ==="
for dir in */; do
    if [ -f "${dir}go.mod" ]; then
        tag_name=$(basename "${dir}")/"${VERSION}"
        git tag -f "${tag_name}"
        echo "Created tag: ${tag_name}"
    fi
done

# 强制推送覆盖远程标签
echo -e "\n=== Push all tags to remote ==="
git push origin --tags -f

echo -e "\n✅ All tags updated!"