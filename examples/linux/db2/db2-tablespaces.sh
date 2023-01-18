#!/bin/sh

cd /home/$1/sqllib

. ./db2profile

db2 connect to $2

db2 -x "SELECT varchar(tbsp_name, 30) as tbsp_name, reclaimable_space_enabled, tbsp_free_pages, tbsp_page_top, tbsp_usable_pages FROM TABLE(MON_GET_TABLESPACE('',-2)) AS t ORDER BY tbsp_free_pages ASC"

exit 0