; ********************************************************************************
; * The basic Mode 7 based commands
; ********************************************************************************

; Display departure boards
; This stores the requested station CRS, makes the request then displays them
; automatically changing screens.
; Then roughly every minute it re-requests an update.
; It repeats until ESC is pressed.
.departures
    JSR genericCommand
    JSR rotatePages                 ; rotate pages or 60s are up
    BRA departures                  ; reload pages

; Handles a plain simple command on the server, display the results and then return to the
; local command line
.genericCommand
    LDY inputBufferStart            ; Start from the beginning of the input buffer

; Perform a Mode7 search
;
; Entry:
;   A           Command code to use
;   inputBuffer containing the search parameter to send
;
.searchMode7
    JSR resetCommandBuffer
    JSR appendInputBuffer           ; Add command line to payload
    LDA #10
    JSR appendCommand

    JSR useCommandRow               ; Show searching logo
    LDX #<searchingText
    LDY #>searchingText
    JSR writeString

    JSR sendCommand                 ; Send command

    JSR useCommandRow
    JSR cls

    JMP showResponse                ; Mode7 response

.appendSearchCommand
    STX tmpaddr
    STY tmpaddr+1
    JSR resetCommandBuffer
    LDY #0
.appendSearchCommand0
    LDA (tmpaddr),Y
    BEQ appendSearchCommand1
    JSR appendCommand
    INY
    BRA appendSearchCommand0
.appendSearchCommand1
    RTS

.searchingText
    EQUS 12, 131, 157, 129, "Searching...", 0