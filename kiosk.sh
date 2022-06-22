#!/bin/bash

# disable screen blanking
env DISPLAY=:0.0 xset s off -dpms

# flip screen, since we mount the RPi upside down
xrandr -o inverted

# run dashboard
# rm /tmp/emd.log
./dash 2>&1 | tee /tmp/emd.log

# debugging
# glxgears -fullscreen
# lxterminal -e top
