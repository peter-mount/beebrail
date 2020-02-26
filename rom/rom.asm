; ********************************************************************************
; * beebrail - Live Departure boards for the BBC Master 128
; ********************************************************************************

        ; Select 65c02
        CPU        1

        ; Paged rom boundaries
        ORG     &8000
        GUARD   &C000

        ; MOS constants
        INCLUDE "rom/mos.asm"

; ********************************************************************************
; Zero page allocations
bufferPos       = &00                   ; Current position of command buffer
sendPos         = &02                   ; Used in command processing
highmem         = &04                   ; HIGHMEM
page            = &06                   ; PAGE
inputBufferPos  = &08                   ; line length in inputBuffer excluding CR
tmpaddr         = &70                   ; 2 bytes for (tmpaddr),Y type calls, BASIC friendly
                                        ; 5 bytes when used for OSWORD &00

; Page 4,5,6 & 7 are for the current language

inputBuffer = &700                      ; Page 7 for the command line input buffer

; ********************************************************************************
; ROM header
.romStart
        JMP language                    ; Language entry point - unused unless bit6 in rom type is set
        JMP serviceEntry                ; Service entry point
        EQUB %11000010                  ; ROM type: Service Entry, Language & 6502 cpu
        EQUB copyright-romStart
        EQUB 1                          ; Version
.title
        EQUS "DepartureBoards.mobi", 0
.version
        INCLUDE "version.asm"           ; Version date is the build date
.copyright
        EQUS 0, "(C)"                   ; Must start with 0 to be valid
        INCLUDE "copyright.asm"         ; Pulls in the build year
        EQUS " Area51.dev", 0

; Core modules
        INCLUDE "rom/writeString.asm"
        INCLUDE "rom/hex.asm"
        INCLUDE "rom/service.asm"
        INCLUDE "rom/protocol.asm"
        INCLUDE "rom/language.asm"

; End of the rom.

; Set ORG to the rom end so when the rom image is generated it is
; exactly 16K
        IF FILLBANK
        ORG &C000
        ENDIF
.romEnd
        SAVE "brail", romStart, romEnd
        PUTTEXT "boot", "!BOOT", 1000
