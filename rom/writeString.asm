; ********************************************************************************
; * Misc string utils
; ********************************************************************************

; writeString - writes a null terminated string
;
; on entry:
;  X,Y Address of string
;
; on exit: A & X preserved, Y invalid
.writeString
	PHA
	STX tmpaddr
	STY tmpaddr+1
	LDY #0
.writeStringLoop
	LDA (tmpaddr),Y
	BEQ writeStringEnd
	JSR osasci
	INY
	BNE writeStringLoop
.writeStringEnd
	PLA
	RTS

; Write a space
.writeSpace
	LDA #' '
	JMP oswrch
