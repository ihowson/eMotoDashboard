# eMoto Dashboard

I built an electric motorcycle. I mounted a Raspberry Pi 3B onto the handlebars to use as a dashboard. This is the software that runs on it.

It has a few major functions.

- Talks to the major hardware components -- BMS, controller, Cycle Analyst -- and collects their status.
- Reports the status through a GUI.
- While the bike is charging, it shows details on how the charge process is progressing.
- Exposes some metrics through HTTP (eventually Prometheus or HomeKit) so the bike can be remotely monitored and controlled.

# Hardware configuration

Schematic is at TODO

You can modify the port mappings in system_target.go. Notably, you can change which device is attached to which serial port. If you don't have the interface board that I built, you probably want to use USB-to-serial adapters for your devices.

# Attributions

- Font Awesome by Dave Gandy - http://fontawesome.io
JBD interface docs by XYZ - URL
