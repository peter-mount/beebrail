# The wire protocol

The BBC talks to the server over RS423 using a simple packet protocol.

| offset | code | description |
| ------ | ---- | ----------- |
| 00 | cmd | Command code |
| 01 | status | Ignored on requests, set on response |
| 02 | len | Length of data |
| 04 | data | Packet data |
| 04 + len | checsum | Simple checksum |

The commands are:

| code | Name | Operation |
| ---- | ---- | --------- |
| 00 | Nop | Does nothing |
| 01 | Time | Requests the time |

## 00 Nop

A packet with 0 data
