; ********************************************************************************
; * search.asm - *SEARCH station name
; ********************************************************************************

;
.search
    JSR enableResult
    LDA #'S'
    JSR startCommand
    LDX #1
.search0
    LDA (&F2),Y
    CMP #' '
    BMI search1
    JSR startCommand
    INY
    BNE search0
.search1
    JSR endCommand
    JMP simpleResult
    RTS
