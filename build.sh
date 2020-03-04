#!/bin/bash

DEST=$(pwd)/dest
DISK=${DEST}/rom.ssd

rm -rf ${DEST}
mkdir -p ${DEST}

# The server
CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOARM="" \
    go build -o dest/beebserver ./bin

# Now the BBC rom

# Write version to be the build date
echo " EQUS \"0.01 ($(date "+%d %b %Y"))\"" >version.asm

# Copyright year
echo " EQUS \"$(date "+%Y")\"" >copyright.asm

rm -f $DISK $TMPDISK boot

# Create the !BOOT script to load the rom into sideways bank 4
(
  # The date here is so we know we are loading the correct image
  # The way I copy an ssd to the gotek is via ssh to a
  # Raspberry PI 0W which is masquerading as the USB stick
  # and it doesn't always pick up the new image immediately
  echo "REM Image build date"
  echo "REM $(date)"
  echo "*SRLOAD BRAIL 8000 4"
) >boot

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
