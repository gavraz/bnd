package game

import (
	"bnd/engine"
	"fmt"
	"time"
)

type meleeObject struct {
	engine.DynamicObject
	angle    float64
	curAngle float64
	lifeTime time.Duration
	radius   float64
}

func newMeleeObject(user *player, angle float64, lifeTime time.Duration, size float64, radius float64) {
	userDir, userCenter, userSize := user.GetDirection(), user.GetCenter(), user.GetWidth()
	sword := &meleeObject{
		DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
			CollisionType: engine.Circle,
			Width:         size,
			Height:        size,
			Mass:          1,
			Until:         time.Now().Add(lifeTime),
			IsPassthrough: true,
		}),
		angle:    angle,
		curAngle: angle,
		lifeTime: lifeTime,
		radius:   radius,
	}
	sword.SetRootParent(user)
	user.AddChild(sword)
	angledDirection := userDir.Rotate(angle)
	sword.SetDirection(angledDirection)
	centerMain := userCenter.Add(angledDirection.MulScalar(userSize))
	sword.SetCenter(centerMain)
	childNumber := 100.0
	for i := 1.0; i <= childNumber*radius; i++ {
		swordParticle := &meleeObject{
			DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
				Center:        centerMain.Add(angledDirection.MulScalar(i / childNumber)),
				CollisionType: engine.Circle,
				Width:         size,
				Height:        size,
				Mass:          1,
				Direction:     angledDirection,
				Until:         time.Now().Add(lifeTime),
				IsPassthrough: true,
			}),
			angle:    angle,
			curAngle: angle,
			lifeTime: lifeTime,
			radius:   radius,
		}
		swordParticle.SetRootParent(sword)
		sword.AddChild(swordParticle)
	}
}

func (m *meleeObject) Update(dt float64) {
	m.curAngle -= 2 * m.angle * dt / m.lifeTime.Seconds()
	parent := m.GetParent()
	center := parent.GetCenter().Add(m.GetDirection().MulScalar(parent.GetWidth()))
	dir := parent.GetDirection().Rotate(m.curAngle)
	m.SetDirection(dir)
	m.SetCenter(center)
	children := m.GetChildren()
	size := float64(len(children))
	for i, child := range children {
		child.SetCenter(center.Add(dir.MulScalar(float64(i+1) * m.radius / size)))
		child.SetDirection(dir)
	}
}

func (m *meleeObject) OnCollision(collider engine.Object) {
	if p, ok := collider.(*player); ok {
		if time.Now().After(p.hitCooldown) {
			p.hp -= 1
			fmt.Println("Hit! \nCurrent hp: ", p.hp)
			p.hitCooldown = time.Now().Add(1000 * time.Millisecond)
		}
	}
}
