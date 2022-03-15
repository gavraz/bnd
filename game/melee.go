package game

type meleeObject struct {
	DynamicObject
	angle    float64
	curAngle float64
	lifeTime float64
	radius   float64
}

func NewMeleeObject(obj DynamicObject, userDir Vector2, userCenter Vector2, userSize float64, angle float64, lifeTime float64, size float64, radius float64) *meleeObject {
	dir := userDir.Rotate(angle)
	obj.SetDirection(dir)
	centerMain := userCenter.Add(dir.MulScalar(userSize))
	obj.SetCenter(centerMain)
	childNumber := int(100 * radius)
	for i := 1; i <= childNumber; i++ {
		sword := &meleeObject{&GObject{
			Center:        centerMain.Add(dir.MulScalar(float64(i) * radius / float64(childNumber))),
			CollisionType: Circle,
			Width:         size,
			Height:        size,
			Mass:          1,
			Direction:     dir,
			ParentObject:  obj,
			TimeToLive:    lifeTime,
			IsPassthrough: true,
		}, 0, 0, 1, 1}
		obj.AddChild(sword)
	}
	return &meleeObject{obj, angle, angle, lifeTime, radius}
}

func (m *meleeObject) UpdateMelee(dt float64) {
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
