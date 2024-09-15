# generate 56 empty func sigs

instr=$(curl -sL https://www.nesdev.org/obelisk-6502-guide/reference.html |
	grep h3)

# generate func sigs for instructions
while read -r inst; do
	break
	name=$(<<< "$inst" cut -d'"' -f2)
	desc=$(echo "$inst" |
		sed -e 's/<br[^>]*>/\n/g; s/<[^>]*>//g' |
		perl -MHTML::Entities -pe 'decode_entities($_);' | tr -d '')
	echo "// $desc
func (c *Cpu) $name() byte {
// https://www.nesdev.org/obelisk-6502-guide/reference.html#$name
return 0
}"
done <<< "$instr"
