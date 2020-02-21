; ********************************************************************************
; * crs.asm - Handles *CRS searchString
; ********************************************************************************

;
.crsSearch
    LDA #cmdCRS
    JSR startCommand
    LDX #1
.crsSearch0
    LDA (&F2),Y
    CMP #' '
    BMI crsSearch1
    JSR startCommand
    INY
    BNE crsSearch0
.crsSearch1
    JSR endCommand
    ; receive result
    RTS
