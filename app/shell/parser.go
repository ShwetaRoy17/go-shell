package shell


func ParseCmd(command string) (string,[]string) {

	s := []string{}
	singleQ, doubleQ, esc := false, false, false
	curr := ""

	n := len(command)
	for i := 0; i < n-1; i++ {
		ch := command[i]
		if esc && doubleQ {
			if !escCh[ch] {
				curr += "\\"

			}
			curr += string(ch)
			esc = false
		} else if esc {
			esc = false
			curr += string(ch)
		} else if ch == '\'' && !doubleQ {
			singleQ = !singleQ
		} else if ch == '"' && !singleQ {
			doubleQ = !doubleQ
		} else if ch == '\\' && !singleQ {
			esc = true
		} else if ch == ' ' && !singleQ && !doubleQ {
			if curr != "" {
				s = append(s, curr)
				curr = ""
			}
		} else {
			curr += (string)(ch)
		}

	}

	if curr != "" {
		s = append(s, curr)
	}
	

	return s[0],s[1:]
}