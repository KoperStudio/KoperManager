@ECHO OFF
FOR /f "tokens=2 delims==" %%f IN ('wmic os get osarchitecture /value ^| find "="') DO SET "OS_ARCH=%%f"
IF "%OS_ARCH%"=="32-bit" GOTO :32bit
IF "%OS_ARCH%"=="64-bit" GOTO :64bit

ECHO OS Architecture %OS_ARCH% is not supported!
EXIT 1

:32bit
ECHO 32 bit Operating System
installer\installer_windows_386.exe
GOTO :SUCCESS

:64bit
ECHO 64 bit Operating System
installer\installer_windows_amd64.exe
GOTO :SUCCESS

:SUCCESS
EXIT 0
