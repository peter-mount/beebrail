; ********************************************************************************
; * The basic Mode 7 based commands
; ********************************************************************************

; Display departure boards
.departures
    LDA inputBuffer,Y               ; Store next 3 chars
    STA currentStation
    INY
    LDA inputBuffer,Y
    STA currentStation+1
    INY
    LDA inputBuffer,Y
    STA currentStation+2
    LDA #10                         ; Char 4 is always terminator
    STA currentStation+3
.departures0                        ; reload departures using currentStation
    ; FIXME optimise this
    LDA #'d':STA inputBuffer
    LDA #'e':STA inputBuffer+1
    LDA #'p':STA inputBuffer+2
    LDA #'a':STA inputBuffer+3
    LDA #'r':STA inputBuffer+4
    LDA #'t':STA inputBuffer+5
    LDA #' ':STA inputBuffer+6
    LDY #0                          ; so copy currentStation into inputBuffer
    LDX #4
.departures1
    LDA currentStation,Y
    STA inputBuffer+7,Y
    INY
    DEX
    BNE departures1

    LDY #0

    JSR searchMode7
    JSR rotatePages                 ; rotate pages or 60s are up
    BRA departures0                 ; reload pages

; Handles the search by crs code
.crsSearch
    LDA #'C'                        ; Command 'C'
    BRA searchMode7

; search performs a station name/crs search just like the field on the
; departureboards.mobi home page.
.search
    LDA #'S'                        ; Command 'S' is for search
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

.searchingText
    EQUS 12, 131, 157, 129, "Searching...", 0