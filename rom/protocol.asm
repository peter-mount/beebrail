; ********************************************************************************
; * The 6502 protocol client
; ********************************************************************************

; internal, resets the command buffer pointer
.resetCommandBuffer
    PHA                         ; Save A & reset buffer
    LDA page
    STA bufferPos
    LDA page+1
    STA bufferPos+1
    PLA
    RTS

.readNextSendPos
    LDA (sendPos)
; internal increments sendPos
.incSendPos
    PHA
    CLC                         ; Move sendPos 1 byte
    LDA #1
    ADC sendPos
    STA sendPos
    LDA #0
    ADC sendPos+1
    STA sendPos+1
    PLA
    RTS

; internal compares bufferPos and sendPos for equality
.bufferPosSendPosEqual
    LDA sendPos                     ; Compare LSB
    CMP bufferPos
    BNE bufferPosSendPosEqualExit   ; LSB not equal
    LDA sendPos+1                   ; Compare MSB
    CMP bufferPos+1
.bufferPosSendPosEqualExit
    RTS

; Starts a command
; Entry:
;   A   command code
; Exit:
;   X Y preserved
.startCommand
    JSR resetCommandBuffer
    STA lastCommand             ; Save command code for later
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
.appendCommandBuffer            ; entry point used to skip in startCommand only
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
    LDA bufferPos               ; tmpaddr = bufferPos - page
    SBC page
    STA tmpaddr
    LDA bufferPos+1
    SBC page+1
    STA tmpaddr+1

    SEC                         ; command length = tmpaddr - 4
    LDA tmpaddr
    SBC #4
    LDY #2
    STA (page),Y
    LDA tmpaddr+1
    SBC #0
    INY
    STA (page),Y

    LDA page                    ; Set sendPos to the start of the buffer
    STA sendPos
    LDA page+1
    STA sendPos+1

.sendCommandLoop
    LDA (sendPos)               ; get current char
    TAY                         ; Send to serial buffer
    LDA #$8A
    LDX #2
    JSR osbyte
    BCS sendCommandLoop         ; Buffer was full

    JSR incSendPos              ; Add 1 to sendPos
    JSR bufferPosSendPosEqual   ; Check for end of buffer
    BNE sendCommandLoop

    JSR resetCommandBuffer      ; Now wait for a response

    LDX #4                      ; Header size
    LDY #0
    JSR readBuffer              ; read the header

    LDY #2                      ; XY = Payload size
    LDA (page),Y
    TAX
    INY
    LDA (page),Y
    TAY
    JSR readBuffer              ; read the payload

    LDY #0                      ; check returned command is the same
    LDA (page),y
    CMP lastCommand
    BEQ readBufferCheckStatus
.errProtocol
    BRK                         ; fail Protocol error
    EQUS 1,"Protocol error",0

.readBufferCheckStatus
    LDY #1                      ; Read response status
    LDA (page),Y
    BNE readBufferError         ; fail on error

; sets sendPos to the start of the payload
.setSendPosPayload
    CLC
    LDA page                        ; Set sendPos to start of the payload
    ADC #4
    STA sendPos
    LDA page+1
    ADC #0
    STA sendPos+1
    RTS

; Copy the payload to inputBuffer as a BRK using the response code as the error number
.readBufferError
    STZ inputBuffer                 ; BRK
    STA inputBuffer+1               ; Status code
    JSR setSendPosPayload           ; Point to payload address
    LDX #2                          ; Start from inputBuffer+2
.readBufferError1
    JSR bufferPosSendPosEqual       ; Loop until we hit bufferPos
    BEQ readBufferError2
    JSR readNextSendPos             ; read payload
    STA inputBuffer,X
    INX
    BRA readBufferError1
.readBufferError2
    STZ inputBuffer,X               ; Append terminating 0
    JMP inputBuffer                 ; Invoke break

; read XY bytes from RS423 and append to the buffer
.readBuffer
    CLC                         ; sendPos = XY + bufferPos
    TXA
    ADC bufferPos
    STA sendPos
    TYA
    ADC bufferPos+1
    STA sendPos+1
.readBufferLoop
    JSR serRead                 ; Read byte
    JSR appendCommand           ; Append to buffer
    JSR bufferPosSendPosEqual   ; loop until all read
    BNE readBufferLoop
    RTS

; Receive
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
    LDX #1                          ; RS423 input buffer
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

; For simple responses which are plain text but in BBC format just
; Write the received payload direct to the output
.simpleResult
    JSR useEntireScreen             ; Use the entire screen (for now)
.simpleResultLoop
    JSR bufferPosSendPosEqual       ; Loop until sendPos hits bufferPos
    BEQ simpleResultEnd             ; Exit once done
    JSR readNextSendPos
    JSR oswrch
    BRA simpleResultLoop
.simpleResultEnd
    RTS

; appendInputBuffer Copies the rest of the inputBuffer (usually after a command)
; to the payload
;
; Entry:
;   Y   Offset in inputBuffer to copy from
.appendInputBuffer
    LDA inputBuffer,Y               ; Read from inputBuffer
    CMP #' '                        ; Stop on first char < 32
    BMI appendInputBuffer1
    JSR appendCommand               ; Append to command buffer
    INY
    BNE appendInputBuffer
.appendInputBuffer1
    RTS
