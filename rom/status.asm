; General debug code, responds to the STATUS command

.status
    LDX #<status0
    LDY #>status0
    JSR writeString
    LDA page+1
    JSR writeHex
    LDA page
    JSR writeHex

    LDX #<status1
    LDY #>status1
    JSR writeString
    LDA highmem+1
    JSR writeHex
    LDA highmem
    JSR writeHex
    JMP osnewl
.status0
    EQUS "PAGE ",0
.status1
    EQUS 13,10,"HIGHMEM ",0