#!/bin/bash

DISK=$(pwd)/rom.ssd
TMPDISK=${DISK}.1

# Write version to be the build date
echo " EQUS \"0.01 ($(date "+%d %b %Y"))\"" >version.asm

# Copyright year
echo " EQUS \"$(date "+%Y")\"" >copyright.asm

rm -f $DISK $TMPDISK

# Compile a local rom image, for programming into an actual EEPROM
# -D FILLBANK=1 tells rom.asm to ensure the image is 16K
beebasm -w \
  -i rom/rom.asm \
  -D FILLBANK=1

# Compile a second time to a DFS disk image
# -D FILLBANK=0 tells rom.asm not to pad the image out which saves time
# when writing to sideways ram
beebasm -w \
  -i rom/rom.asm \
  -D FILLBANK=0 \
  -title BeebRail \
  -opt 3 \
  -do $DISK
