
\\-------------------------------------------------------------------------
\\ Subroutine to print a byte in A in hex form (destructive)
\\ Based on PRBYTE from the Apple 1 monitor written by Woz
\\-------------------------------------------------------------------------
.writeHex
PHA				\\ Save A for LSD
LSR A:LSR A:LSR A:LSR A		\\ MSD to LSD position
JSR writeHexChar		\\ Output hex digit
PLA				\\ Restore A
\\ Fall through to print hex routine
\\-------------------------------------------------------------------------
\\ Subroutine to print a hexadecimal digit
\\-------------------------------------------------------------------------
.writeHexChar
AND #%00001111			\\ Mask LSD for hex print
ORA #'0'			\\ Add "0"
CMP #'9'+1			\\ Is it a decimal digit?
BCC writeHexCharEcho		\\ Yes! output it
ADC #6				\\ Add offset for letter A-F
.writeHexCharEcho
JMP oswrch
