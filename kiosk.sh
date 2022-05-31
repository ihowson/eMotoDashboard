#!/bin/bash

#glxgears -fullscreen
# lxterminal -e top

# disable screen blanking
env DISPLAY=:0.0 xset s off -dpms

# flip screen, since we mount the RPi upside down
xrandr -o inverted

# run dashboard
./dash
