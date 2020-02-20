; ********************************************************************************
; * mos.asm - The BBC MOS
; ********************************************************************************

; OS calls
oswrch = &FFEE
osasci = &FFE3
osbyte = &FFF4

; Zero page
oswordReason = &EF          ; EF contains OSWORD reason code
oswordData   = &F0          ; F0/F1 contains parameter block address
pagedRomID   = &F4          ; F4 contains the currently active paged rom
