; ********************************************************************************
; The language entry point, for now a simple command line interface
; ********************************************************************************

; switchLanguage Switches to our language
; The language ROM is entered via its entry point with A=1.
; Locations &FD and &FE in zero page are set to point to the copyright message in the ROM.
.switchLanguage
    LDA #&8E                        ; Enter language ROM
    LDX pagedRomID                  ; Use our ROM number
    JMP osbyte

.showPrompt
    JSR usePrompt
    LDX #<prompt
    LDY #>prompt
    JSR writeString

; The main language entry point
.language
    CMP #&01                        ; Accept A=1 only
    BEQ language1
    RTS
.language1
    LDX #&FF                        ; Reset the stack
    TXS

    LDA #22                         ; Switch to shadow mode 7
    JSR oswrch
    LDA #128+7
    JSR oswrch

    LDA #&84                        ; set HIGHMEM
    JSR osbyte
    STX highmem
    STY highmem+1

    DEC A                           ; set PAGE
    JSR osbyte
    STX page
    STY page+1

    JSR enableResult                ; Enable RS423

    LDA #<errorHandler              ; Setup error handler
    STA BRKV
    LDA #>errorHandler
    STA BRKV+1

    CLI                             ; Enable IRQ's

    JSR dial                        ; Dial the server

    JSR homePage                    ; show home page

; cmdLine editor
.cmdLine
    JSR showPrompt
    JSR useCommandRow               ; Use command row only
    LDA #12                         ; Clear command row
    JSR oswrch

    LDA #>inputBuffer               ; Page 7 is our input buffer
    STZ tmpaddr                     ; Input buffer
    STA tmpaddr+1
    LDA #&EE                        ; Max length, EE = basic
    STA tmpaddr+2
    LDA #' '                        ; Lowest acceptable character
    STA tmpaddr+3
    LDY #&FF                        ; Highest acceptable character
    STY tmpaddr+4
    INY                             ; XY=tmpAddr
    LDX #tmpaddr
    TYA                             ; Call osword 0
    JSR osword
    BCC cmdLineNoEscape             ; CC Escape not pressed
    JMP errEscape                   ; Escape

.cmdLineNoEscape
    STY inputBufferPos              ; Save end of command line

    LDY #0                          ; Find first non-space character
.cmdSearch0
    LDA inputBuffer,Y
    CMP #13                         ; End of command, so nothing entered
    BEQ cmdLine
    CMP #33
    BPL cmdSearch1                  ; Found something
    INY
    BNE cmdSearch0
    BRA cmdSearchError              ; overflow

; Handles * commands
; Entry:
;   Y   Offset in inputBuffer of *
.cmdOscli
    INY                             ; Skip *
    TYA
    TAX
    LDY #>inputBuffer
    JSR oscli
    BRA cmdLine

.cmdSearch1
    CMP #'*'
    BEQ cmdOscli                    ; oscli command

    STY inputBufferStart            ; Save so we can resume from same position
    LDX #0
.cmdSearch2
    LDY inputBufferStart            ; reset Y to start of input
.cmdSearch3
    LDA langTable,X
    BMI cmdSearchError              ; End of table so error
    BNE cmdSearch4                  ; Not at end of command

    LDA inputBuffer,Y               ; next char in inputBuffer must be ' ' or control char
    CMP #' '
    BEQ cmdSearchFoundLong          ; Long command with space
    BMI cmdSearchFoundShort         ; Short command no args

; Skip until we hit the next entry in the table
.cmdSearchSkip
    INX
    LDA langTable,X
    BNE cmdSearchSkip               ; Loop until 0
    INX:INX:INX                     ; Skip 0 & address
    BRA cmdSearch2                  ; Resume loop

.cmdSearch4
    LDA inputBuffer,Y               ; Get char from inputBuffer
    JSR toUpper                     ; Convert to uppercase
    CMP langTable,X
    BNE cmdSearchSkip               ; Skip to next command
    INY                             ; Next char
    INX
    BNE cmdSearch3                  ; Loop to next char
.cmdSearchError
    JMP errSyntax

; We have found the command so call the command handler.
; Entry:
;   Y   contains the offset in inputBuffer of the character at the end of the command
.cmdSearchFoundLong                 ; Entry point if command terminated by space
    INY                             ; Skip whitespace to point to first non-space arg
    LDA inputBuffer,Y
    CMP #' '
    BEQ cmdSearchFoundLong          ; loop until we hit non-space
.cmdSearchFoundShort                ; Entry point if command not terminated by space
    PHY                             ; Save Y then find command terminator which for BBC is CR (13)

    PHY                             ; Next we lower case the command - as server side they are all lower case
.cmdSearchFound0
    DEY
    BMI cmdSearchFound1
    LDA inputBuffer,Y
    JSR toLower
    STA inputBuffer,Y
    BRA cmdSearchFound0
.cmdSearchFound1
    PLY

.cmdSearchFound2
    LDA inputBuffer,Y               ; Look for first control char
    JSR oswrch
    CMP #' '
    BMI cmdSearchFound3             ; Found it
    INY
    BRA cmdSearchFound2             ; Loop until we find it

.cmdSearchFound3
    LDA #10                         ; Set the terminator to UNIX LF (10)
    STA inputBuffer,Y

    PLY                             ; Restore Y to point to start of argument

    INX                             ; Skip to handler address
    JSR cmdSearchExec               ; exec command
    JMP cmdLine                     ; Back to prompt

; As there's no JSR (langTable,X)
; On entry (&F2),Y will hold offset of first char after the command (if required)
; Y should hold position in inputBuffer of next char after command
.cmdSearchExec
    JMP (langTable,X)

; toUpper converts A to uppercase
.toUpper
    CMP #'a'                        ; Check for lower case letter
    BMI toUpper0                    ; if not then skip case conversion
    CMP #'z'+1
    BPL toUpper0
    AND #&5f                        ; convert to uppercase
.toUpper0
    RTS

; toLower converts A to lowercase
.toLower
    CMP #'A'                        ; Check for lower case letter
    BMI toLower0                    ; if not then skip case conversion
    CMP #'Z'+1
    BPL toLower0
    ORA #&20                        ; convert to lowercase
.toLower0
    RTS

.prompt
    EQUS 30, 134, "departureboards.mobi", 135, '*', 0

.langTable
    EQUS "DEPART", 0        : EQUW departures       ; Departures
    EQUS "CRS", 0           : EQUW genericCommand   ; CRS search
    EQUS "HELP", 0          : EQUW homePage         ; Help
    EQUS "SEARCH", 0        : EQUW genericCommand   ; Search crs by name
    EQUS "STATUS", 0        : EQUW status           ; Debug status
    EQUB &FF                                        ; Table terminator

    INCLUDE "rom/error.asm"
    INCLUDE "rom/status.asm"
    INCLUDE "rom/search.asm"
