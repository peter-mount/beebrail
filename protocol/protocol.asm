; ********************************************************************************
; * The 6502 protocol client
; ********************************************************************************

; Where the command buffer lies
commandBuffer = &1000
bufferPos = &80

cmdNOP = 0
cmdCRS = 'C'

; Ends a command
.endCommand
    LDA #0              ; 0 terminator, fall through to startCommand

; Starts a command
; Entry:
;   A   command code
; Exit:
;   X Y preserved
.startCommand
    PHX:PHY
    TAY
    LDA #$8A
    LDX #2
    JSR osbyte
    PLY:PLX
    RTS

; osbyte buffers
; 0 keyboard
; 1 RS423 input
; 2 RS423 output
; 3 Printer
;
; osbyte 15,X flush buffer X
; osbyte 8A,X,Y insert Y into buffer X, C=1 if not inserted, 0 if inserted
; osbyte 91,X read char from buffer into Y, C=1 if empty, 0 if data read

; Read from RS423 into A
; Exit
;   A   character read
;   XY  preserved
.serRead
    PHY                             ; Preserve X & Y
    PHX
.serRead0
    LDA #&91
    LDX #1
    JSR osbyte
    BCS serRead0                    ; No data read
    TYA                             ; Result in Y
    PLX                             ; Restore X & Y then return
    PLY
    RTS