@REM cd to root
cd %~dp0

set GOOS=js
set GOARCH=wasm

go build -o ./website/public/SOMAS2020.wasm