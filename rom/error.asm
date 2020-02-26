; ********************************************************************************
; ********************************************************************************

; Handle errors from BRK instructions.
; Errors raised with:
;   BRK
;   EQUB error code (unused)
;   EQUS error string
;   BRK or EQUB 0 to terminate string
.errorHandler
    LDX #&FF                        ; Reset the stack
    TXS

    LDA #&7C                        ; Clear escape condition
    JSR osbyte

   ;LDA (brkAddress)                ; error code if we want to compare it for custom handling

    JSR osnewl                      ; Force newline
    LDY #1                          ; Skip the error code
.errorHandler0
    LDA (brkAddress),Y
    BEQ errorHandler1               ; Found end
    JSR osasci
    INY
    BNE errorHandler0
.errorHandler1
    JSR osnewl                      ; Force newline
    JMP cmdLine                     ; Back to command prompt

.errEscape
    BRK
    EQUS &11, "Escape", 0

.errSyntax
    BRK
    EQUS &12, "Syntax", 0
