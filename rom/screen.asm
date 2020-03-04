; ********************************************************************************
; Handles our custom screen
; ********************************************************************************

; Use the entire screen
.useEntireScreen
    LDX #1                  ; Disable cursor
    LDY #0
    JSR vdu23

    LDA #0                  ; Set text window
    LDX #0
    LDY #24
    BRA setTextViewPort

; Restrict scrolling to the top line
.useCommandRow
    LDX #1                  ; Enable cursor
    LDY #1
    JSR vdu23

    LDA #23                 ; Set text window
    LDX #0
    LDY #0
; Sets the text view port. X being the top line number, Y the bottom line number of the page
.setTextViewPort
    PHA                     ; VDU 28,A,Y,39,X       sets text view port
    LDA #28
    JSR oswrch
    PLA
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

; Main home page
.homePage
    LDX #<homePage0         ; Display home page
    LDY #>homePage0
    JMP writeString
.homePage0
    EQUB 22, 128+7,10       ; Mode 7 in shadow
    EQUS 132,157,135,141, 31, 10, 1, "UK Departure Boards", 13, 10
    EQUS 132,157,135,141, 31, 10, 2, "UK Departure Boards", 13, 10
    EQUS 135, "List of available commands:", 13, 10
    EQUS 131, "CRS crs", 31, 13, 4, 134, "Search by CRS", 13, 10
    EQUS 131, "HELP", 31, 13, 5, 134, "Show help", 13, 10
    EQUS 131, "SEARCH name", 31, 13, 6, 134,"Search by name"
    EQUB 0
