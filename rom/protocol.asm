; ********************************************************************************
; * The 6502 protocol client
; ********************************************************************************

; Where the command buffer lies
commandBuffer = &1000
bufferPos = &80
sendPos = &82

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
    PHA                         ; Save A & reset buffer
    LDA #<commandBuffer
    STA bufferPos
    LDA #>commandBuffer
    STA bufferPos+1
    PLA
    JSR appendCommand           ; Append command code
    LDA #0                      ; Append 0 (as this is the status byte)
    JSR appendCommand
    JSR appendCommandBuffer     ; Skip 2 bytes for the buffer size
    BRA appendCommandBuffer     ; We'll set these up later

; Append a byte to the command Buffer
; Entry:
;   A   Value to append
; Exit:
;   A   corrupted
;   X Y preserved
;
.appendCommand
    STA (bufferPos)             ; Store
.appendCommandBuffer
    CLC                         ; Move bufferPos 1 byte
    LDA #1
    ADC bufferPos
    STA bufferPos
    LDA #0
    ADC bufferPos+1
    STA bufferPos+1
    RTS

.sendCommand
    SEC                         ; Set command length
    LDA bufferPos
    SBC #<(commandBuffer+4)
    STA commandBuffer+2
    LDA bufferPos+1
    SBC #>(commandBuffer+4)
    STA commandBuffer+3

    LDA #<commandBuffer         ; Set command pointer
    STA sendPos
    LDA #>commandBuffer
    STA sendPos+1

.sendCommandLoop
    LDA (sendPos)               ; get current char
    TAY                         ; Send to serial buffer
    LDA #$8A
    LDX #2
    JSR osbyte
    BCS sendCommandLoop         ; Buffer was full

    CLC                         ; Move sendPos 1 byte
    LDA #1
    ADC sendPos
    STA sendPos
    LDA #0
    ADC sendPos+1
    STA sendPos+1

    LDA sendPos
    CMP bufferPos
    BNE sendCommandLoop         ; LSB not equal

    LDA sendPos+1
    CMP bufferPos+1
    BNE sendCommandLoop         ; MSB not equal

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
    LDA #&91                        ; Read character from buffer
    LDX #1                          ;
    JSR osbyte
    BCS serRead0                    ; No data read so loop
    TYA                             ; Result in Y
    PLX                             ; Restore X & Y then return
    PLY
    RTS

.enableResult
    PHA:PHX:PHY

    LDA #&02                        ; Use keyboard for input but listen to serial port
    LDX #2
    LDY #0
    JSR osbyte

    LDA #&B5                        ; RS423 input taken as raw data, default but enforce it
    LDX #1
    LDY #0
    JSR osbyte

    LDA #&CC                        ; lets serial data enter the input buffer
    LDX #0
    LDY #0
    JSR osbyte
    PLY:PLX:PLA
    RTS

; For simple responses which are plain text but in BBC format just read from RS423 until we get a 0
.simpleResult
    JSR serRead
    CMP #0
    BEQ simpleResultEnd
    JSR osasci
    bra simpleResult
.simpleResultEnd
    LDA #13
    JMP osasci
