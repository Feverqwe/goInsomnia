Build app
---
````
go build -ldflags -H=windowsgui -trimpath
````

Build resources with go-bindata
---
````
go-bindata .\icon.ico
````

File icon, use rsrc 
---
````
.\rsrc_windows_amd64.exe -ico .\icon.ico -o FILE.syso
````
