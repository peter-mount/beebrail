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
    STY inputBufferPos
    BRA cmdLine                     ; do nothing for now

    INCLUDE "rom/error.asm"
    INCLUDE "rom/status.asm"
