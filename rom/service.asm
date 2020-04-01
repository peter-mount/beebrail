; ********************************************************************************
; * service.asm - rom service handler
; * See http://beebwiki.mdfs.net/Paged_ROM_service_calls
; ********************************************************************************

; service call handler
;
; On entry:
;         A Service call number, reason code
;         X Page ROM slot number (also in memory at address &F4)
;         Y Parameter to service call
;
; On exit:
;        A 0 if call has been dealt with, otherwise preserved
;        Y return value
;
; Note: The individual handlers will have X corrupted so use &F4 to get the rom
;         slot number if required
.serviceEntry
        PHA                                 ; Preserve a and use as lookup ID
        LDX #0                              ; Start at serviceLookup Table
.serviceEntryLookup
        STA tmpaddr
.serviceEntryLoop
        LDA serviceLookup,X
        BEQ serviceEntryNotSupported
        CMP tmpaddr
        BEQ serviceEntryFound
        INX:INX:INX                         ; Skip to next table entry
        BNE serviceEntryLoop

.serviceEntryNotSupported                   ; No entry so restore A & return
        PLA
        RTS

.serviceEntryFound                          ; Service entry matched
        INX                                 ; Skip to entrypoint
        PLA                                 ; Restore A
        JMP (serviceLookup,X)               ; Jump to handler. This should set A=0 if consumed, preserve A if not

        ; Service lookup table: 3 bytes per entry, the service call number, the address (2 bytes).
        ; A 0 terminates the table.
        ; For performance, put the most used calls first in the table!
.serviceLookup
        EQUB &04:EQUW oscliHandler          ; 04 unrecognised OSCLI
        ;EQUB &08:EQUW oswordHandler         ; 08 Unrecognised OSWORD
        ;EQUB &27:EQUW serviceReset          ; 27 Reset
        EQUB 0                              ; end of table
