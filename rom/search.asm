; ********************************************************************************
; * The basic Mode 7 based commands
; ********************************************************************************

; Display departure boards
.departures
    LDA #'D'                        ; Command 'D'
    BRA searchMode7

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
    JSR startCommand
    JSR appendInputBuffer           ; Add command line to payload

    JSR useCommandRow               ; Show searching logo
    LDX #<searchingText
    LDY #>searchingText
    JSR writeString

    JSR sendCommand                 ; Send command
    JMP showResponse                ; Mode7 response

.searchingText
    EQUS 12, 131, 157, 129, "Searching...", 0