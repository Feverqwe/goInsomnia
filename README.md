<div align="center">
	<img src="assets/screenshot1.png" alt=""/>
	<h1>GoInsomnia</h1>
	<p>
		<b>Prevent computer from going to sleep</b>
	</p>
	<br>
	<br>
	<br>
</div>

Build exe
---
````
go build -ldflags -H=windowsgui -trimpath -o goInsomnia.exe
````

Build mac app
---
```
go build -trimpath -o goInsomnia
go get github.com/strosel/appify
go install github.com/strosel/appify
~/go/bin/appify -menubar -name goInsomnia -author "Anton V" -id com.rndnm.goinsomnia -icon ./assets/icon.icns goInsomnia
rm ./goInsomnia.app/Contents/README
```

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
