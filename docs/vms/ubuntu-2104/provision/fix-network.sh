#!/usr/bin/env bash

netplan set "ethernets.enp0s8.nameservers.addresses=[8.8.8.8, 8.8.4.4]"
netplan apply