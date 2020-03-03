; ********************************************************************************
; Handles our custom screen
; ********************************************************************************

; Use the entire screen
.useEntireScreen
    LDX #1                  ; Disable cursor
    LDY #0
    JSR vdu23
    LDX #0
    LDY #24
    BRA setTextViewPort

; Restrict scrolling to the top line
.useCommandRow
    LDX #1                  ; Enable cursor
    LDY #1
    JSR vdu23

    LDX #0
    LDY #0
; Sets the text view port. X being the top line number, Y the bottom line number of the page
.setTextViewPort
    LDA #28                 ; VDU 28,0,Y,39,X       sets text view port
    JSR oswrch
    LDA #0
    JSR oswrch
    TYA
    JSR oswrch
    LDA #39
    JSR oswrch
    TXA
    JSR oswrch

    LDY #16
    LDX #%0010000          ; viewport wraps

; vdu23 handles simple flag settings
;
; Equivalent to VDU 23,X,Y;0;0;0
;
.vdu23
    LDA #23                 ; VDU 23,16,f;0;0;0     Disable scrolling
    JSR oswrch
    TXA
    JSR oswrch
    TYA
    JSR oswrch
    LDY #7                  ; Remaining 7 bytes are 0
    LDA #0
.setTextViewPort0
    JSR oswrch
    DEY
    BNE setTextViewPort0
    RTS
