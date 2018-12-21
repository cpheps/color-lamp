#!/bin/bash
if ! ifconfig wlan0 | grep -q "inet addr:" ; then
    wpa_cli scan_results | grep WPS | sort -r -k3 | tail -1 | awk '{print $1;}' > /tmp/wifi
    read ssid < /tmp/wifi
    wpa_cli -i wlan0 wps_pbc $ssid
fi