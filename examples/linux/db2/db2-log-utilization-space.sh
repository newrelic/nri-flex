#!/bin/sh

cd /home/$1/sqllib

. ./db2profile

db2 connect to $2

db2 -x "SELECT TOTAL_LOG_USED_KB, TOTAL_LOG_AVAILABLE_KB, TOTAL_LOG_USED_TOP_KB from SYSIBMADM.LOG_UTILIZATION"

exit 0