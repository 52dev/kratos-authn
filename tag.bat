@echo off
setlocal
set "BASE_DIR=%cd%"
set "VERSION=v1.0.0"

git tag %VERSION%

for /d %%D in (*) do (
    if exist "%%D\go.mod" (
        git tag %%D/%VERSION%
    )
)

git push origin --tags
echo All tags pushed!
pause