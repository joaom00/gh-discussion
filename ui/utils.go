package ui

func truncateString(str string, num int) string {
	truncated := str
	if num <= 3 {
		return str
	}
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		truncated = str[0:num] + "…"
	}
	return truncated
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func (m *Model) prevDisc() {
	m.cursor.currDiscId = max(m.cursor.currDiscId-1, 0)
}

func (m *Model) nextDisc() {
	newDiscID := min(m.cursor.currDiscId+1, len(*m.data)-1)
	newDiscID = max(newDiscID, 0)
	m.cursor.currDiscId = newDiscID
}
