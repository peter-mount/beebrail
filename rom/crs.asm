; ********************************************************************************
; * crs.asm - Handles *CRS searchString
; ********************************************************************************

;
.crsSearch
    LDA #'C'                        ; Command 'C'
    JSR startCommand
    JSR appendInputBuffer           ; Add command line to payload
    JSR sendCommand                 ; Send command
    JMP simpleResult                ; Simple result just plain text
