package types

const TypeA uint16 = 1
const TypeNS uint16 = 2
const TypeCNAME uint16 = 5
const TypeSOA uint16 = 6
const TypePTR uint16 = 12
const TypeMX uint16 = 15
const TypeTXT uint16 = 16
const TypeAAAA uint16 = 28
const TypeSRV uint16 = 33
const TypeOPT uint16 = 41
const TypeAXFR uint16 = 252
const TypeMAILB uint16 = 253
const TypeMAILA uint16 = 254
const TypeAll uint16 = 255

const ClassIN uint16 = 1
const ClassCS uint16 = 2
const ClassCH uint16 = 3
const ClassHS uint16 = 4
const ClassAll uint16 = 255

const RCodeNoError uint8 = 0
const RCodeFormatError uint8 = 1
const RCodeServerFailure uint8 = 2
const RCodeNameError uint8 = 3
const RCodeNotImplemented uint8 = 4
const RCodeRefused uint8 = 5
