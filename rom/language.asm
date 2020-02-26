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

; The main language entry point
.language
    CMP #&01                        ; Accept A=1 only
    BEQ language1
    RTS
.language1
    LDX #&FF                        ; Reset the stack
    TXS

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

    JSR status  ; Temp debug

; cmdLine editor
.cmdLine
    LDA #'R'
    JSR oswrch
    LDA #'>'
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
   ;LDA #0                          ; Set last char (13) to 0
   ;STA inputBuffer,Y               ; There's no STZ inputBuffer,Y

    LDY #0                          ; Find first non-space character
.cmdSearch0
    LDA inputBuffer,Y
    CMP #13                         ; End of command
    BEQ cmdSearchError
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

    STY tmpaddr                     ; Save so we can resume from same position
    LDX #0
.cmdSearch2
    LDY tmpaddr                     ; reset Y to start of input
.cmdSearch3
    LDA langTable,X
    BMI cmdSearchError              ; End of table so error
    BEQ cmdSearchFound              ; 0 so we have found our command
    CMP inputBuffer,Y
    BNE cmdSearchSkip               ; Skip to next command
    INY:INX                         ; Next char
    BNE cmdSearch3                  ; Loop to next char
.cmdSearchError
    JMP errSyntax

.cmdSearchFound
    INX                             ; Skip to address
    JSR cmdSearchExec               ; exec command
    JMP cmdLine                     ; Back to prompt

.cmdSearchSkip
    INX
    LDA langTable,X
    BNE cmdSearchSkip               ; Loop until 0
    INX:INX:INX                     ; Skip 0 & address
    BRA cmdSearch2                  ; Resume loop

; As there's no JSR (langTable,X)
; On entry (&F2),Y will hold offset of first char after the command (if required)
; Y should hold position in inputBuffer of next char after command
.cmdSearchExec
    JMP (langTable,X)

.langTable
    EQUS "STATUS", 0        : EQUW status       ; Debug status
    EQUS "SEARCH ", 0       : EQUW search       ; Search crs by name
    EQUB &FF                                    ; Table terminator

    INCLUDE "rom/error.asm"
    INCLUDE "rom/status.asm"
