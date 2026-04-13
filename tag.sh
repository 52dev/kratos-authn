#!/bin/bash
set -e
BASE_DIR=$(pwd)
VERSION="v1.0.1"

# 强制创建/覆盖根标签
echo "=== Force create root tag: ${VERSION} ==="
git tag -f "${VERSION}"

# 递归查找 1~3 层目录内的 go.mod，自动打子模块标签
echo -e "\n=== Force create submodule tags (recursive 1~3 levels) ==="
find . -maxdepth 3 -type f -name "go.mod" | grep -v '^\./go.mod$' | while read -r modfile; do
    # 获取模块目录（去掉 ./ 和 go.mod）
    moddir=$(dirname "${modfile}")
    reldir=${moddir#./}  # 相对路径，如 transport/activemq

    if [[ -n "${reldir}" && -f "${modfile}" ]]; then
        tag_name="${reldir}/${VERSION}"
        git tag -f "${tag_name}"
        echo "✅ Created tag: ${tag_name}"
    fi
done

# 强制推送覆盖远程所有标签
echo -e "\n=== Push all tags to remote (force overwrite) ==="
git push origin --tags -f

echo -e "\n🎉 All tags updated successfully!"