@REM sc stop "GoApp"
@REM sc delete "GoApp"
nssm.exe install GoApp "C:\Servicios\hello.exe"
nssm.exe set GoApp AppDirectory "C:\Servicios"
nssm.exe set GoApp DisplayName "GoApp"
nssm.exe set "GoApp" Start SERVICE_AUTO_START