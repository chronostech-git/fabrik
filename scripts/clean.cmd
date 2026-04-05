@echo off
setlocal

echo.
    set /p DATADIR=Enter data directory:

if "%DATADIR%"=="" (
    echo No directory supplied.
    exit /b 1
)

echo.
    echo Running make clean...
make clean

echo.
echo Deleting data directory: %DATADIR%

if exist "%DATADIR%" (
    rmdir /s /q "%DATADIR%"
    echo Directory removed.
) else (
    echo Directory "%DATADIR%" does not exist.
)

echo.
    echo Done.
endlocal
