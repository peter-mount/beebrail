; ********************************************************************************
; Handles our custom screen
; ********************************************************************************

; Use the entire screen
.useEntireScreen
    LDX #1                  ; Disable cursor
    LDY #0
    JSR vdu23

    LDA #0                  ; Set text window
    LDX #1
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
; Sets the text view port.
; Entry:
;   A   left column
;   X   top line number
;   Y   bottom line number
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

    LDX #16                 ; VDU 23,16,f;0;0;0     Disable scrolling
    LDY #%0010000           ; viewport wraps

; vdu23 handles simple flag settings
;
; Equivalent to VDU 23,X,Y;0;0;0
;
.vdu23
    LDA #23                 ; VDU 23,X,Y;0;0;0     Disable scrolling
    JSR oswrch
    TXA
    JSR oswrch
    TYA
    JSR oswrch
    LDY #7                  ; Remaining 7 bytes are 0
    LDA #0
.vdu23Loop
    JSR oswrch
    DEY
    BNE vdu23Loop
    RTS

; cls - clears the main screen
.cls
    JSR useEntireScreen             ; Use the entire screen (for now)
    LDA #12
    JMP oswrch

; Main home page
.homePage
    JSR useEntireScreen
    LDA #12
    JSR oswrch
    LDX #0
    LDY #homePageEnd-homePageStart
.homePageLoop
    LDA homePageStart,X
    JSR oswrch
    INX
    DEY
    BNE homePageLoop
    RTS
.homePageStart
    EQUS 30
    EQUS 132, 157,135,141, 31, 10, 0, "UK Departure Boards", 13, 10
    EQUS 132, 157,135,141, 31, 10, 1, "UK Departure Boards", 13, 10
    EQUS 135, "List of available commands:", 13, 10
    EQUS 131, "CRS crs", 31, 13, 3, 134, "Search by CRS", 13, 10
    EQUS 131, "HELP", 31, 13, 4, 134, "Show help", 13, 10
    EQUS 131, "SEARCH name", 31, 13, 5, 134,"Search by name"
.homePageEnd