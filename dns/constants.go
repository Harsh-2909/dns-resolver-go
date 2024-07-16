package dns

// Root DNS server
const (
	RootDNS     = "198.41.0.4"
	RootDNSPort = 53
)

// DNS record types
const (
	TypeA     uint16 = 1   // IPv4 address record
	TypeNS    uint16 = 2   // authoritative name server record
	TypeCNAME uint16 = 5   // canonical name record
	TypeSOA   uint16 = 6   // start of authority record
	TypePTR   uint16 = 12  // pointer record
	TypeMX    uint16 = 15  // mail exchange record
	TypeTXT   uint16 = 16  // text record
	TypeAAAA  uint16 = 28  // IPv6 address record
	TypeSRV   uint16 = 33  // service locator record
	TypeOPT   uint16 = 41  // option record
	TypeAXFR  uint16 = 252 // transfer of an entire zone record
	TypeMAILB uint16 = 253 // mailbox-related records (MB, MG, MR)
	TypeMAILA uint16 = 254 // mail agent RRs (Obsolete - see MX)
	TypeAll   uint16 = 255 // all records
)

// DNS record classes
const (
	ClassIN  uint16 = 1   // Internet class
	ClassCS  uint16 = 2   // CSNET class (Obsolete)
	ClassCH  uint16 = 3   // CHAOS class
	ClassHS  uint16 = 4   // Hesiod [Dyer 87]
	ClassAll uint16 = 255 // all classes
)

// DNS response codes
const (
	RCodeNoError        uint8 = 0 // No error condition
	RCodeFormatError    uint8 = 1 // Format error - The name server was unable to interpret the query
	RCodeServerFailure  uint8 = 2 // Server failure - The name server was unable to process this query due to a problem with the name server
	RCodeNameError      uint8 = 3 // Name error - Meaningful only for responses from an authoritative name server, this code signifies that the domain name referenced in the query does not exist
	RCodeNotImplemented uint8 = 4 // Not implemented - The name server does not support the requested kind of query
	RCodeRefused        uint8 = 5 // Refused - The name server refuses to perform the specified operation for policy reasons
)
