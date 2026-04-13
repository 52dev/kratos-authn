@echo off
setlocal
set "BASE_DIR=%cd%"
go mod tidy
for /d %%D in (*) do (
    if exist "%%D\go.mod" (
        cd "%%D"
        go mod tidy
        cd "%BASE_DIR%"
    )
)
echo All done!
pause