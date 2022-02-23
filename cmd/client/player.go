package main

type Ulti int // for now
type point [2]int

type Player struct {
	hp        int
	position  point
	direction float32 // maybe make it 360 degrees and for now just use 0,90,180,270 for simplicity and update later if needed
	velocity  float32
	ultiz     [2]Ulti
	alive     bool // Maybe change to state if we have other changes to status (like slow/stun...)
}

func (p Player) Attack() {

}

func (p Player) Move(dir float32) {

}

func (p Player) TakeDamage(dmg int) bool {
	p.hp -= dmg
	p.alive = p.hp > 0
	return p.alive
}

func (p Player) Heal() {

}

func (p Player) OpenCrate(c crate) {
	// Need crate object
}
