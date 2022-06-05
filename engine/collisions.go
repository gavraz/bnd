package engine

import "math"

func (e *Environment) CollidesWith(obj DynamicObject, other Object) bool {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		return CircleCollidesWithCircle(obj, other)
	}

	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Rectangle {
		return CircleCollidesWithRectangle(obj, other)
	}

	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Circle {
		return RectangleCollidesWithCircle(obj, other)
	}

	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Rectangle {
		return RectangleCollidesWithRectangle(obj, other)
	}

	return false
}

func (g *gameObject) onDynamicCollisionCircles(collider DynamicObject) {
	if g.GetIsPassthrough() || collider.GetIsPassthrough() {
		return
	}
	p1 := g.GetCenter()
	p2 := collider.GetCenter()
	r1 := g.GetWidth() / 2
	r2 := collider.GetWidth() / 2
	n := p1.Sub(p2).Normalize()
	v1 := g.GetVelocity()
	v2 := collider.GetVelocity()
	a1 := v1.Dot(n)
	a2 := v2.Dot(n)
	// Using the optimized version,
	// optimizedP =  2(a1 - a2)
	//              -----------
	//                m1 + m2
	optimizedP := (2.0 * (a1 - a2)) / (g.GetMass() + collider.GetMass())
	u1 := v1.Sub(n.MulScalar(optimizedP * collider.GetMass()))
	u2 := v2.Add(n.MulScalar(optimizedP * g.GetMass()))
	g.SetVelocity(u1)
	collider.SetVelocity(u2)

	penetrationDepth := r1 + r2 - p1.Distance(p2)
	direction := p1.Sub(p2).Normalize()
	g.SetCenter(p1.Add(direction.MulScalar(penetrationDepth)))

}
func (g *gameObject) onDynamicCollisionRectangles(collider DynamicObject) {
	if g.GetIsPassthrough() || collider.GetIsPassthrough() {
		return
	}
	p1, p2 := g.GetCenter(), collider.GetCenter()
	w1, w2 := g.GetWidth(), collider.GetWidth()
	h1, h2 := g.GetHeight(), collider.GetHeight()
	v1, v2 := g.GetVelocity(), collider.GetVelocity()
	m1, m2 := g.GetMass(), collider.GetMass()

	overlapX := math.Min(p1.X+w1/2, p2.X+w2/2) - math.Max(p1.X-w1/2, p2.X-w2/2)
	overlapY := math.Min(p1.Y+h1/2, p2.Y+h2/2) - math.Max(p1.Y-h1/2, p2.Y-h2/2)
	u1 := v1.MulScalar(m1 - m2).Add(v2.MulScalar(2 * m2)).DivScalar(m1 + m2)
	u2 := v1.MulScalar(2 * m1).Sub(v2.MulScalar(m1 - m2)).DivScalar(m1 + m2)

	if overlapX > overlapY {
		g.SetVelocity(Vector2{X: v1.X, Y: u1.Y})
		collider.SetVelocity(Vector2{X: v2.X, Y: u2.Y})
		if p1.Y < p2.Y {
			// Collision on top
			g.SetCenter(Vector2{X: p1.X, Y: p1.Y - overlapY})
		} else {
			// Collision on bottom
			g.SetCenter(Vector2{X: p1.X, Y: p1.Y + overlapY})
		}
	} else {
		g.SetVelocity(Vector2{X: u1.X, Y: v1.Y})
		collider.SetVelocity(Vector2{X: u2.X, Y: v2.Y})
		if p1.X < p2.X {
			// Collision on left
			g.SetCenter(Vector2{X: p1.X - overlapX, Y: p1.Y})

		} else {
			// Collision on right
			g.SetCenter(Vector2{X: p1.X + overlapX, Y: p1.Y})

		}
	}
}
func (g *gameObject) onDynamicCollisionCircleRectangle(collider DynamicObject) {
	circle, rect := g, collider
	handleCircleRectangleCollision(circle, rect)
}
func (g *gameObject) onDynamicCollisionRectangleCircle(collider DynamicObject) {
	rect, circle := g, collider
	handleCircleRectangleCollision(circle, rect)
}

func (g *gameObject) onStaticCollisionCircles(collider StaticObject) {
	if g.GetIsPassthrough() {
		return
	}
	p1 := g.GetCenter()
	p2 := collider.GetCenter()
	r1 := g.GetWidth() / 2
	r2 := collider.GetWidth() / 2

	n := p1.Sub(p2).Normalize()
	v1 := g.GetVelocity()
	u1 := n.MulScalar(v1.Length())
	g.SetVelocity(u1)
	penetrationDepth := r1 + r2 - p1.Distance(p2)
	direction := p1.Sub(p2).Normalize()
	g.SetCenter(p1.Add(direction.MulScalar(penetrationDepth)))
}
func (g *gameObject) onStaticCollisionRectangles(collider StaticObject) {
	if g.GetIsPassthrough() {
		return
	}
	p1, p2 := g.GetCenter(), collider.GetCenter()
	w1, w2 := g.GetWidth(), collider.GetWidth()
	h1, h2 := g.GetHeight(), collider.GetHeight()
	v1 := g.GetVelocity()

	overlapX := math.Min(p1.X+w1/2, p2.X+w2/2) - math.Max(p1.X-w1/2, p2.X-w2/2)
	overlapY := math.Min(p1.Y+h1/2, p2.Y+h2/2) - math.Max(p1.Y-h1/2, p2.Y-h2/2)

	if overlapX > overlapY {
		g.SetVelocity(Vector2{X: v1.X, Y: -v1.Y})
		if p1.Y < p2.Y {
			// Collision on top
			g.SetCenter(Vector2{X: p1.X, Y: p1.Y - overlapY})
		} else {
			// Collision on bottom
			g.SetCenter(Vector2{X: p1.X, Y: p1.Y + overlapY})
		}
	} else {
		g.SetVelocity(Vector2{X: -v1.X, Y: v1.Y})
		if p1.X < p2.X {
			// Collision on left
			g.SetCenter(Vector2{X: p1.X - overlapX, Y: p1.Y})

		} else {
			// Collision on right
			g.SetCenter(Vector2{X: p1.X + overlapX, Y: p1.Y})
		}
	}

}
func (g *gameObject) onStaticCollisionCircleRectangle(collider StaticObject) {
	circle, rect := g, collider
	if g.GetIsPassthrough() {
		return
	}
	if collider.GetIsPassthrough() {
		return
	}
	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()
	penetrationVector := dist.Normalize().MulScalar(penetrationDepth)

	if circle.GetVelocity().Dot(dist) < 0 {
		tangentVel := dist.Normalize().Dot(circle.GetVelocity())
		circle.SetVelocity(circle.GetVelocity().Sub(dist.Normalize().MulScalar(tangentVel * 2)))
	}
	circle.SetCenter(circle.GetCenter().Add(penetrationVector))

}
func (g *gameObject) onStaticCollisionRectangleCircle(collider StaticObject) {
	rect, circle := g, collider
	if g.GetIsPassthrough() {
		return
	}
	if collider.GetIsPassthrough() {
		return
	}
	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: NearestX - circle.GetCenter().X, Y: NearestY - circle.GetCenter().Y}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()
	penetrationVector := dist.Normalize().MulScalar(penetrationDepth)

	tangentVel := dist.Normalize().Dot(rect.GetVelocity())
	rect.SetVelocity(rect.GetVelocity().Sub(dist.Normalize().MulScalar(tangentVel * 2)))
	rect.SetCenter(rect.GetCenter().Add(penetrationVector))

}

func CircleCollidesWithCircle(obj DynamicObject, other Object) bool {
	p1 := obj.GetCenter()
	p2 := other.GetCenter()
	r1 := obj.GetWidth() / 2
	r2 := other.GetWidth() / 2
	dist := p1.Distance(p2)
	return dist <= r1+r2
}
func CircleCollidesWithRectangle(obj DynamicObject, other Object) bool {
	circle, rect := obj, other

	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()

	return penetrationDepth > 0.0
}
func RectangleCollidesWithCircle(obj DynamicObject, other Object) bool {
	rect, circle := obj, other

	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()

	return penetrationDepth > 0.0

}
func RectangleCollidesWithRectangle(obj DynamicObject, other Object) bool {
	// check if two rectangles are intersecting
	p1, p2 := obj.GetCenter(), other.GetCenter()
	w1, w2 := obj.GetWidth(), other.GetWidth()
	h1, h2 := obj.GetHeight(), other.GetHeight()
	if p1.X+w1/2 < p2.X-w2/2 || p1.X-w1/2 > p2.X+w2/2 || p1.Y+h1/2 < p2.Y-h2/2 || p1.Y-h1/2 > p2.Y+h2/2 {
		// No collision
		return false
	} else {
		// Collision
		return true
	}
}

func handleCircleRectangleCollision(circle DynamicObject, rect DynamicObject) {
	if circle.GetIsPassthrough() || rect.GetIsPassthrough() {
		return
	}
	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()
	penetrationVector := dist.Normalize().MulScalar(penetrationDepth)

	if penetrationDepth > 0.0 {
		if circle.GetVelocity().Dot(dist) < 0 {
			tangentVel := dist.Normalize().Dot(circle.GetVelocity())
			combinedMass := circle.GetMass() + rect.GetMass()
			circle.SetVelocity(circle.GetVelocity().Sub(dist.Normalize().MulScalar(tangentVel * rect.GetMass() / combinedMass * 2)))
			rect.SetVelocity(rect.GetVelocity().Add(dist.Normalize().MulScalar(tangentVel * circle.GetMass() / combinedMass * 2)))
		}
		circle.SetCenter(circle.GetCenter().Add(penetrationVector))
	}
}
