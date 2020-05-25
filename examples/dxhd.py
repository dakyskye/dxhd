#!/usr/bin/python3

# super + Return
import os
os.popen('alacritty')

#super + d
import os
os.popen('rofi -modi drun -show drun')

#super+{_,shift}+{1-9,0}
import os
os.popen('zenity --no-wrap --info --text="i3-msg -t command {_,move container to} workspace {1-9,10}"')

#super+ctrl+{_,shift}+{1-9,0}
import os
os.popen('zenity --no-wrap --info --text="i3-msg -t command {_,move container to} workspace {11-19,20}"')

