@echo off
setlocal

set APP_NAME=StickerDownloadBot

for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set "dt=%%I"
set BUILD_TIME=%dt:~0,4%-%dt:~4,2%-%dt:~6,2% %dt:~8,2%:%dt:~10,2%:%dt:~12,2%

echo Building %APP_NAME%.exe ...
echo buildTime: %BUILD_TIME%

go build -trimpath -ldflags="-s -w -X 'main.buildTime=%BUILD_TIME%'" -o %APP_NAME%.exe .

if %ERRORLEVEL% equ 0 (
    echo Done: %APP_NAME%.exe
) else (
    echo Build failed.
)

endlocal
