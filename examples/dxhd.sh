#!/bin/sh

# super + Return
alacritty

# super + d
rofi -modi drun -show drun

# super + ctrl + n
redshift -P -O 4500

# super + ctrl + d
redshift -P -O 6500

#super + {_,shift + }{1-9,0}
i3-msg -t command {_,move container to} workspace {1-9,10}

#super+ctrl+{_,shift}+{1-9,0}
i3-msg -t command {_,move container to} workspace {11-19,20}
