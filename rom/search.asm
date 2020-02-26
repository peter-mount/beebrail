; ********************************************************************************
; * search.asm - *SEARCH station name
; ********************************************************************************

; search performs a station name/crs search just like the field on the
; departureboards.mobi home page.
;
; The returned results are written to the console
;
.search
    JSR enableResult
    LDA #'S'                        ; Command 'S' is for search
    JSR startCommand
.search0
    LDA (&F2),Y                     ; Read from OSCLI buffer
    CMP #' '                        ; Stop on first char < 32
    BMI search1
    JSR appendCommand               ; Append to command buffer
    INY
    BNE search0
.search1
    JSR sendCommand                 ; Send command
    JMP simpleResult                ; Simple result just plain text
    RTS
