; ********************************************************************************
; Dialer - handles hanging up and dialing the server using the WiFi modem
; ********************************************************************************

.dial
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
