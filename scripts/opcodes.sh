# generate 128 opcodes (as map items)

# empty user agent header is fine, apparently
opcodes=$(curl -sL 'http://www.6502.org/tutorials/6502opcodes.html' -H 'User-Agent:' |
	grep -P '\$[0-9A-F]{2}  \d')

# doesn't work for stack instructions, but there are only 6 of them, so writing
# them manually is trivial
while read -r op; do
	addr=$(<<< "$op" cut -b-11 | tr -d ' ,')
	inst=$(<<< "$op" cut -b15-18)
	byte=$(<<< "$op" awk '{print $(NF-2)}' | sed 's/\$/0x/')
	cyc=$(<<< "$op" awk '{print $NF}' | tr -d '+')
	echo "$byte: {Instruction: (*Cpu).$inst, Cycles: $cyc, AddressingMode: $addr},"
done <<< "$opcodes"
