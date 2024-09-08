# decor
It is the  utility to set Windows desktop background

decor.exe [-s] path_to_jpeg_file

It starts with single argument: path to jpeg file.
In case before an argument there is the flag -s program starts in silent mode (without console)

Program reads jpeg file, encodes it, crops image so it's width-to-height ratio matches width-to-height ratio of computer monitor, saves the image as bmp file and sets it as a screen background.

It is pure  Go program. It should be compiled with command  
* go build -ldflags -H=windowsgui

That allows to assign windows console dynamically or suppress console  by  flag -s.


