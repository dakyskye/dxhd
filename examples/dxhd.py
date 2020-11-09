#!/usr/bin/python3

import os

# super + Return
os.popen('alacritty')

#super + d
os.popen('rofi -modi drun -show drun')

#super+{_,shift}+{1-9,0}
os.popen('zenity --no-wrap --info --text="PY i3-msg -t command {_,move container to} workspace {1-9,10}"')

#super+ctrl+{_,shift}+{1-9,0}
os.popen('zenity --no-wrap --info --text="PY i3-msg -t command {_,move container to} workspace {11-19,20}"')
