; ********************************************************************************
; * The 6502 protocol client
; ********************************************************************************

ASCII_STX   = 2
ASCII_ETX   = 3
ASCII_GS    = 29
ASCII_RS    = 30

.dial
    JSR banner
    JSR showPrompt

    JSR useCommandRow               ; Show dialing logo
    LDX #<dialingText
    LDY #>dialingText
    JSR writeString

    LDA #0                          ; Viewport to show dial progress
    LDX #22
    LDY #24
    JSR setTextViewPort

    LDX #16                         ; Enable scrolling
    LDY #0
    JSR vdu23

    LDA #12                         ; Clear new viewport
    JSR oswrch

    LDY #0                          ; Start dialing
.dial1
    LDA dialText,y
    BNE dial2
    RTS
.dial2
    INY
    CMP #1
    BEQ dialWaitSecond
    JSR serWrite
    CMP #10
    BNE dial1
.dailWait
    LDA #100
    LDX #0
    BRA dialWait0
.dialWaitSecond
    LDA #0
    LDX #1
.dialWait0
    STA currentStation              ; Store timer val
    STX currentStation+1
    STY totalPages                  ; Store Y

    STZ tmpaddr                     ; Reset timer
    STZ tmpaddr+1
    STZ tmpaddr+2
    STZ tmpaddr+3
    STZ tmpaddr+4
    JSR writeTimer

.dialWait1
    LDA #&91                        ; Read character from buffer
    LDX #1                          ; RS423 input buffer
    JSR osbyte
    BCS dialWait2                   ; No data read so skip
    TYA                             ; Write received char to screen
    JSR oswrch

.dialWait2
    JSR readTimer

    LDA tmpaddr+1
    CMP currentStation+1
    BMI dialWait1
    LDA tmpaddr
    CMP currentStation
    BMI dialWait1                   ; Loop until timer hit

    LDY totalPages
    BRA dial1

.readTimer
    LDA #3                          ; Read timer
    BRA readWriteTimer
.writeTimer
    LDA #4
.readWriteTimer
    LDX #<tmpaddr
    LDY #>tmpaddr
    JMP osword

.dialingText
    EQUS 12, 131, 157, 129, "Dialing...", 0

.dialText
    EQUS "+++", 1
    EQUS "ATH", 13, 10
    EQUS "ATZ", 13, 10
    EQUS "ATDTlocalhost:8082", 13, 10
    EQUS "mode bbc api", 13, 10
    EQUB 0

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
    JSR incBufferPos            ; Skip 2 bytes for the buffer size
    BRA incBufferPos            ; We'll set these up later

; Read next byte from command buffer
.readCommandBuffer
    LDA (bufferPos)
    PHA
    JSR incBufferPos
    PLA
    RTS

; Append a byte to the command Buffer
; Entry:
;   A   Value to append
; Exit:
;   A   corrupted
;   X Y preserved
;
.appendCommand
    STA (bufferPos)             ; Store
; increment bufferPos by one
.incBufferPos
    CLC                         ; Move bufferPos 1 byte
    LDA #1
    ADC bufferPos
    STA bufferPos
    LDA #0
    ADC bufferPos+1
    STA bufferPos+1
    RTS

.sendCommand
    LDA page                    ; Set sendPos to the start of the buffer
    STA sendPos
    LDA page+1
    STA sendPos+1

.sendCommandLoop
    LDA (sendPos)               ; get current char
    JSR serWrite

    JSR incSendPos              ; Add 1 to sendPos
    JSR bufferPosSendPosEqual   ; Check for end of buffer
    BNE sendCommandLoop

    JSR resetCommandBuffer      ; Now wait for a response
    JSR readBuffer

    ; FIXME this needs error handling
    RTS
    ;LDY #0                      ; check response code
    ;LDA (page),y
    ;CMP lastCommand
    ;BEQ readBufferCheckStatus
.errProtocol
    BRK                         ; fail Protocol error
    EQUS 1,"Protocol error",0

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

; read a table from RS423
.readBuffer
    JSR serRead                     ; Wait until we get STX
    CMP #ASCII_STX
    BNE readBuffer
    JSR appendCommand
.readBuffer0
    JSR serRead
    PHA
    JSR appendCommand
    PLA
    CMP #ASCII_ETX                  ; Loop until ETX
    BNE readBuffer0
    RTS

.serWrite
    PHA
    PHX
    PHY
    TAY                             ; Send to serial buffer
    LDA #$8A
    LDX #2
    JSR osbyte
    PLY
    PLX
    PLA
    BCS serWrite                    ; Buffer was full
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
    JSR cls
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

; Shows mode7 response
.showResponse
    JSR resetCommandBuffer

; Redisplay the current page
.displayPage
    JSR useCommandRow               ; Clear command row
    LDA #12
    JSR oswrch

    JSR cls                         ; Clear screen
.displayPage0
    JSR readCommandBuffer
    CMP #ASCII_STX                  ; STX Start so skip
    BEQ displayPage0
    CMP #ASCII_ETX                  ; ETX so stop but reset the buffer
    BEQ displayPage1
    CMP #ASCII_GS                   ; GS new table so stop
    BEQ displayPage2
    CMP #ASCII_RS                   ; RS new page in table so stop
    BEQ displayPage2
    JSR osunixlf                    ; Write char
    BRA displayPage0
.displayPage1
    JSR resetCommandBuffer
.displayPage2
    RTS

; rotate pages every 5 seconds
.rotatePages
    LDA #12                         ; 12 * 5 = 60 seconds
    STA reloadCounter
.rotatePages0
    LDA #&81
    LDX #<500
    LDY #>500
    JSR osbyte
    BCC rotatePages1                ; No error
    CPY #&1B                        ; Escape?
    BNE rotatePages1                ; Back to next page
    BRK:BRK:BRK                     ; Error 0 no message which will then bail to command prompt

.rotatePages1
    DEC reloadCounter
    BNE rotatePages2
    RTS

.rotatePages2
    LDA totalPages                  ; just 1 page then do nothing
    CMP #1
    BEQ rotatePages0

    INC curPage                     ; cycle to next page
    JSR displayPage
    BRA rotatePages0
