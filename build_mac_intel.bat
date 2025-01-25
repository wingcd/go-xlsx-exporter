:: compile.bat
@echo off

cd /d %~dp0

set prjPath=%~dp0
set binPath=%prjPath%bin\

echo Project Dir:%prjPath%
echo Target Dir:%binPath%

set GOARCH=arm64
set GOOS=darwin3
set GOPATH=%prjPath%../../

go build -o %binPath%gxe-intel
xcopy /y %prjPath%conf.yaml %binPath%
xcopy /y /s /e %prjPath%template %binPath%template

echo build is completed.