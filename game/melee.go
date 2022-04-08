package game

import (
	"bnd/engine"
	"time"
)

type meleeObject struct {
	engine.DynamicObject
	angle    float64
	curAngle float64
	lifeTime float64
	radius   float64
}

func newMeleeObject(obj engine.DynamicObject, userDir engine.Vector2, userCenter engine.Vector2, userSize float64, angle float64, lifeTime float64, size float64, radius float64) *meleeObject {
	dir := userDir.Rotate(angle)
	obj.SetDirection(dir)
	centerMain := userCenter.Add(dir.MulScalar(userSize))
	obj.SetCenter(centerMain)
	childNumber := 100 * radius
	for i := 1.0; i <= childNumber; i++ {
		sword := &meleeObject{&engine.GObject{
			Center:        centerMain.Add(dir.MulScalar(i * radius / childNumber)),
			CollisionType: engine.Circle,
			Width:         size,
			Height:        size,
			Mass:          1,
			Direction:     dir,
			ParentObject:  obj,
			Until:         time.Now().Add(time.Duration(lifeTime*1000) * time.Millisecond),
			IsPassthrough: true,
		}, 0, 0, 1, 1}
		obj.AddChild(sword)
	}
	return &meleeObject{obj, angle, angle, lifeTime, radius}
}

func (m *meleeObject) update(dt float64) {
	m.curAngle -= 2 * m.angle * dt / m.lifeTime
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
