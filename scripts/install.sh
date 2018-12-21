#!/bin/bash

echo Deploying Lamp Files

# Copying Depdencies
echo Copying Header Files
cp headers/*.h /usr/local/include/

echo Copying Libraries
cp lib/* /usr/local/lib

# Copy Lamp config and binary
echo Copying Lamp
cp color_lamp /home/pi/
cp config/config.toml /home/pi/

# Copy WPS files
echo Copying WPS Files
cp scripts/wps_script.sh /home/pi/
chmod +x /home/pi/wps_script.sh

# Copy systemd
echo Copying Service Files
cp systemd/* /etc/systemd/system

# Enabling Service
echo Enabling services
systemctl enable color_lamp
systemctl enable wps_scan
systemctl enable  systemd-networkd.service
systemctl enable  systemd-networkd-wait-online.service