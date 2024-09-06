# decor
It is the Windows utility to set desktop background

decor [-s] path_to_jpeg_file

It starts with single argument: path to jpeg file.
In case  flag -s used before an argument program starts in silent mode (without console)

Program read jpeg file, encode it, crop image so it's dimensions corresponds to dimensions of computer monitor, save the image as bmp file and set it as a screen background.

It is pure  Go program. It should be compiled with command  
* go build -ldflags -H=windowsgui

That allows to assign windows console dynamically or suppress console  by  flag -s.


