; ********************************************************************************
; Displays our banner/splash page
; ********************************************************************************

.banner
    JSR cls
    LDA #<(bannerStart)
    STA tmpaddr
    LDA #>(bannerStart)
    STA tmpaddr+1
    STZ tmpaddr+2
    LDA #&7C
    STA tmpaddr+3

    LDA #&6C                    ; Switch to shadow ram
    LDX #1
    JSR osbyte

.banner0
    LDA (tmpaddr)
    STA (tmpaddr+2)

    CLC
    LDA tmpaddr
    ADC #1
    STA tmpaddr
    LDA tmpaddr+1
    ADC #0
    STA tmpaddr+1

    CLC
    LDA tmpaddr+2
    ADC #1
    STA tmpaddr+2
    LDA tmpaddr+3
    ADC #0
    STA tmpaddr+3

    LDA tmpaddr+1
    CMP #>bannerEnd
    BNE banner0
    LDA tmpaddr
    CMP #<bannerEnd
    BNE banner0

    LDA #&6C                    ; Switch to main ram
    LDX #0
    JSR osbyte

    RTS

    ; Class 165 from https://zxnet.co.uk/teletext/gallery/index.php?gallery=gallery2
.bannerStart
    INCBIN "rom/banner.m7"
.bannerEnd
