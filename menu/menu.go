package menu

type action func() *Menu

type Menu struct {
	actions []action
	labels  []string
}

type Handler struct {
	prev          []*Menu
	current       *Menu
	currentChoice int
}

func NewHandler(menu *Menu) *Handler {
	return &Handler{
		current: menu,
	}
}

func (h *Handler) PrevChoice() {
	if h.currentChoice == 0 {
		return
	}

	h.currentChoice--
}

func (h *Handler) NextChoice() {
	if h.currentChoice == len(h.current.actions)-1 {
		return
	}

	h.currentChoice++
}

func (h *Handler) nextMenu(next *Menu) {
	h.prev = append(h.prev, h.current)
	h.current = next
	h.currentChoice = 0
}

func (h *Handler) Choose() {
	act := h.current.actions[h.currentChoice]
	next := act()
	if h.current != next {
		h.nextMenu(next)
	}
}

func (h *Handler) GoBack() {
	if h.prev == nil || len(h.prev) == 0 {
		return
	}

	last := len(h.prev) - 1
	h.current = h.prev[last]
	h.prev[last] = nil
	h.prev = h.prev[:last]
}

func (h *Handler) Choices() []string {
	c := make([]string, len(h.current.labels))
	copy(c, h.current.labels[:])
	return c
}

func (h *Handler) CurrentChoice() int {
	return h.currentChoice
}
