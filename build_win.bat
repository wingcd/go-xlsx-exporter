:: compile.bat
@echo off

cd /d %~dp0

set prjPath=%~dp0
set binPath=%prjPath%bin\

echo Project Dir:%prjPath%
echo Target Dir:%binPath%

set GOARCH=amd64
set GOOS=windows
set GOPATH=%prjPath%../../

go build -o %binPath%gxe.exe
xcopy /y %prjPath%conf.yaml %binPath%
xcopy /y %prjPath%conf.template.yaml %binPath%
xcopy /y /s /e %prjPath%template %binPath%template

echo build is completed.