#!/bin/sh

# super + Return
alacritty

# super + d
rofi -modi drun -show drun

#super + {_,shift + }{1-9,0}
zenity --no-wrap --info --text="i3-msg -t command {_,move container to} workspace {1-9,10}"

#super+ctrl+{_,shift}+{1-9,0}
zenity --no-wrap --info --text="i3-msg -t command {_,move container to} workspace {11-19,20}"
