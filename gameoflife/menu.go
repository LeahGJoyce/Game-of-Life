package main

func GetMenu() []string {
	return []string{"Start Simulation", "Quit"}
}

func RenderMenu(m model) string {
	s := ""
	for i, choice := range m.menu {
		if m.cursor == i {
			s += ">"
		} else {
			s += " "
		}
		s += choice + "\n"
	}
	return s
}
