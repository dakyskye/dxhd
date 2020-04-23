#!/bin/sh

## launch termite
# Mod4 + Return
termite

## launch rofi
# Mod4 + d
rofi -modi drun -show drun

## launch window switcher
# Mod4 + Tab
rofi -modi window -show window

## let's script something
# Mod4 + space
FOO=5
BAR="baz"

[ $FOO -eq 5 ] && zenity --no-wrap --info --text="hello dxhd"

BAR="$BAR$FOO"

[ $BAR = "baz5" ] && zenity --no-wrap --info --text="scripting in dxhd is easay"
