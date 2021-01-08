#!/bin/bash

lsusb | \
while read line; \
do 
  PART1=$(echo $line | awk '{gsub(/:$/,"",$4); print $2";"$4";"$6}');
  PART2=$(echo $line | cut -d" " -f 7-);
  echo $PART1\;$PART2;
done

