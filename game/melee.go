package game

import (
	"bnd/engine"
	"math"
	"time"
)

const (
	angle               = math.Pi / 4
	lifeTime            = 150 * time.Millisecond
	size                = 0.01
	radius              = 0.1
	numOfChildParticles = 100.0
)

type meleeObject struct {
	engine.DynamicObject
	angle    float64
	curAngle float64
	lifeTime time.Duration
	radius   float64
}

func addMeleeObject(user *player) {
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
	user.AddChild(sword)
	angledDirection := userDir.Rotate(angle)
	sword.SetDirection(angledDirection)
	centerMain := userCenter.Add(angledDirection.MulScalar(userSize))
	sword.SetCenter(centerMain)
	for i := 1.0; i <= numOfChildParticles*radius; i++ {
		swordParticle := &meleeObject{
			DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
				Center:        centerMain.Add(angledDirection.MulScalar(i / numOfChildParticles)),
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
		p.getHit(1)
	}
}
