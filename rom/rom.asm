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

tmpaddr = &70                           ; 2 bytes for (tmpaddr),Y type calls

; ROM header
.romStart
        EQUB 0,0,0                      ; Language entry point - unused unless bit6 in rom type is set
        JMP serviceEntry                ; Service entry point
        EQUB %10000010                  ; ROM type: Service Entry & 6502 cpu
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
        INCLUDE "rom/service.asm"
;        INCLUDE "rom/osword.asm"
;        INCLUDE "lib/writeString.asm"

; End of the rom.

; Set ORG to the rom end so when the rom image is generated it is
; exactly 16K
        IF FILLBANK
        ORG &C000
        ENDIF
.romEnd
        SAVE "brail", romStart, romEnd
        PUTTEXT "boot", "!BOOT", 1000
