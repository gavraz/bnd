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

type meleeParticle struct {
	engine.DynamicObject
	parent *meleeObject
}

type meleeObject struct {
	engine.DynamicObject
	angle    float64
	curAngle float64
	lifeTime time.Duration
	radius   float64
	collided bool
}

func addMeleeObject(player *player) {
	userDir, userCenter, userSize := player.GetDirection(), player.GetCenter(), player.GetWidth()
	meleeObj := &meleeObject{
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
	player.AddChild(meleeObj)
	angledDirection := userDir.Rotate(angle)
	meleeObj.SetDirection(angledDirection)
	centerMain := userCenter.Add(angledDirection.MulScalar(userSize))
	meleeObj.SetCenter(centerMain)
	for i := 0; i < numOfChildParticles; i++ {
		particleObj := &meleeParticle{
			DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
				Center:        centerMain.Add(angledDirection.MulScalar(float64(i) * radius / numOfChildParticles)),
				CollisionType: engine.Circle,
				Width:         size,
				Height:        size,
				Mass:          1,
				Direction:     angledDirection,
				Until:         time.Now().Add(lifeTime),
				IsPassthrough: true,
			}),
			parent: meleeObj,
		}
		meleeObj.AddChild(particleObj)
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
	for i, child := range children {
		child.SetCenter(center.Add(dir.MulScalar(float64(i) * m.radius / numOfChildParticles)))
		child.SetDirection(dir)
	}
}

func (mp *meleeParticle) OnCollision(collider engine.Object) {
	if p, ok := collider.(*player); ok {
		if !mp.parent.collided {
			mp.parent.collided = true
			p.applyDamage(1)
		}
	}
}
