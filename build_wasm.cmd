@echo off

@REM cd to root
cd %~dp0
setlocal 
    set GOOS=js
    set GOARCH=wasm

    go build -ldflags="-w -s" -o ./website/public/SOMAS2020.wasm
endlocal