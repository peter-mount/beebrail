; ********************************************************************************
; * oscli - Implements our OSCLI commands
; ********************************************************************************

; unknown oscli command (&F2),Y holds the start of the command
.oscliHandler
    PHA:PHY                         ; Save A & Y
    LDX #0
.oscliHandlerLoop
    LDA oscliTable,X
    BMI oscliHandlerFail            ; End of table
    BEQ oscliHandlerFound           ; 0 so we have found our command
    CMP (&F2),Y
    BNE oscliHandlerSkip            ; Skip entry
    INY:INX                         ; next char
    BNE oscliHandlerLoop

.oscliHandlerFail
    PLY:PLA                         ; Restore A & Y as its unclaimed
    RTS

.oscliHandlerFound
    PLA:PLA                         ; Dump original A & Y to bitbucket
    INX                             ; Skip to address
    JSR oscliHandlerExec            ; Invoke command
    LDA #0                          ; Claim command
    RTS

.oscliHandlerSkip
    INX
    LDA oscliTable,X
    BNE oscliHandlerSkip            ; Loop until 0
    INX:INX:INX                     ; Skip 0 & address
    PLY:PHY                         ; Get initial Y back
    BRA oscliHandlerLoop            ; Resume loop

; As there's no JSR (oscliTable,X)
; On entry (&F2),Y will hold offset of first char after the command (if required)
.oscliHandlerExec
    JMP (oscliTable,X)

.oscliTable
    EQUS "RAIL",0       : EQUW switchLanguage   ; Language entry point
    EQUS "CRS ",0       : EQUW crsSearch    ; Display CRS
    EQUS "SEARCH ",0    : EQUW search       ; Search for name
    EQUB 0                                  ; Table terminator

    INCLUDE "rom/crs.asm"
    INCLUDE "rom/search.asm"
