package network

func IDMatcher(m1, m2 []byte) bool {
	m1ID := m1[0:2]
	m2ID := m2[0:2]

	return m1ID[0] == m2ID[0] && m1ID[1] == m2ID[1]
}
