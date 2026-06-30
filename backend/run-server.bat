@echo off
cd /d "d:\code-space\coding-soler\permission-system-0d41021b\backend"
echo === Start build at %date% %time% === > run-server.log
echo. >> run-server.log
echo === go mod tidy === >> run-server.log
go mod tidy >> run-server.log 2>&1
echo EXIT CODE: %errorlevel% >> run-server.log
echo. >> run-server.log
echo === go build === >> run-server.log
go build -v -mod=mod -o server.exe . >> run-server.log 2>&1
echo EXIT CODE: %errorlevel% >> run-server.log
echo. >> run-server.log
if exist server.exe (
    echo server.exe created successfully >> run-server.log
    echo. >> run-server.log
    echo === Starting server at %date% %time% === >> run-server.log
    start /B "" "d:\code-space\coding-soler\permission-system-0d41021b\backend\server.exe" >> run-server.log 2>&1
) else (
    echo ERROR: server.exe NOT created >> run-server.log
)
echo. >> run-server.log
echo === Done at %date% %time% === >> run-server.log
