#!/bin/bash

# disable screen saver
env DISPLAY=:0.0 xset s off -dpms

# flip screen, since we mount the RPi upside down
xrandr -o inverted

# run dashboard
# rm /tmp/emd.log
while true
do
    ./dash-armhf 2>&1 | tee /tmp/emd.log
    sleep 1
done

# debugging
# glxgears -fullscreen
# x11vnc &
# lxterminal -e top
