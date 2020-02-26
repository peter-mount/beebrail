; ********************************************************************************
; * search.asm - *SEARCH station name
; ********************************************************************************

; search performs a station name/crs search just like the field on the
; departureboards.mobi home page.
;
; The returned results are written to the console
;
.search
    LDA #'S'                        ; Command 'S' is for search
    JSR startCommand
    JSR appendInputBuffer           ; Add command line to payload
    JSR sendCommand                 ; Send command
    JMP simpleResult                ; Simple result just plain text
