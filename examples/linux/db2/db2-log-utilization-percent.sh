#!/bin/sh

cd /home/$1/sqllib

. ./db2profile

db2 connect to $2

db2 -x "SELECT LOG_UTILIZATION_PERCENT, cast(( TOTAL_LOG_USED_KB/1024) as Integer) as TOTAL_LOG_USED_MB, cast((TOTAL_LOG_AVAILABLE_KB/1024) as Integer) as TOTAL_LOG_AVAILABLE_MB, cast((TOTAL_LOG_USED_TOP_KB/1024) as integer) as TOTAL_LOG_USED_TOP_MB from SYSIBMADM.MON_TRANSACTION_LOG_UTILIZATION"

exit 0