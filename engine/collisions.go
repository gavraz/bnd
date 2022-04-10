package engine

import "math"

func (e *Environment) CollidesWith(obj DynamicObject, other Object) bool {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		return e.CircleCollidesWithCircle(obj, other)
	}

	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Rectangle || obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Circle {
		return e.CircleCollidesWithRectangle(obj, other)
	}

	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Rectangle {
		return e.RectangleCollidesWithRectangle(obj, other)
	}

	return false
}

func CheckDynamicCollision(obj, other DynamicObject) DynamicObject {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		r1 := obj.GetWidth() / 2
		r2 := other.GetWidth() / 2
		dist := p1.Distance(p2)
		if dist <= r1+r2 {
			if obj.GetIsPassthrough() || other.GetIsPassthrough() {
				return other
			}
			//n := p1.Sub(p2).Normalize()
			//v1 := obj.GetVelocity()
			//v2 := other.GetVelocity()
			//a1 := v1.Dot(n)
			//a2 := v2.Dot(n)
			//// Using the optimized version,
			//// optimizedP =  2(a1 - a2)
			////              -----------
			////                m1 + m2
			//optimizedP := (2.0 * (a1 - a2)) / (obj.GetMass() + other.GetMass())
			//u1 := v1.Sub(n.MulScalar(optimizedP * other.GetMass()))
			//u2 := v2.Add(n.MulScalar(optimizedP * obj.GetMass()))
			//obj.SetVelocity(u1)
			//other.SetVelocity(u2)
			penetrationDepth := r1 + r2 - p1.Distance(p2)
			direction := p1.Sub(p2).Normalize()
			obj.SetCenter(p1.Add(direction.MulScalar(penetrationDepth)))
			return other
		}
	}
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Rectangle || obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Circle {
		// https://stackoverflow.com/questions/45370692/circle-rectangle-collision-response
		circle, rect := obj, other
		if obj.GetCollisionType() == Rectangle {
			circle, rect = other, obj
		}

		NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
		NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
		dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

		penetrationDepth := circle.GetWidth()/2 - dist.Length()
		penetrationVector := dist.Normalize().MulScalar(penetrationDepth)

		if penetrationDepth > 0.0 {
			if circle.GetIsPassthrough() || rect.GetIsPassthrough() {
				return other
			}

			if circle.GetVelocity().Dot(dist) < 0 {
				tangentVel := dist.Normalize().Dot(circle.GetVelocity())
				combinedMass := circle.GetMass() + rect.GetMass()
				circle.SetVelocity(circle.GetVelocity().Sub(dist.Normalize().MulScalar(tangentVel * rect.GetMass() / combinedMass * 2)))
				rect.SetVelocity(rect.GetVelocity().Add(dist.Normalize().MulScalar(tangentVel * circle.GetMass() / combinedMass * 2)))
			}
			circle.SetCenter(circle.GetCenter().Add(penetrationVector))
			return other
		}
	}
	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Rectangle {
		// check if two rectangles are intersecting
		// if they are, calculate the penetration depth and the penetration vector
		// then move the rectangles apart by the penetration vector
		// return the other rectangle
		p1, p2 := obj.GetCenter(), other.GetCenter()
		w1, w2 := obj.GetWidth(), other.GetWidth()
		h1, h2 := obj.GetHeight(), other.GetHeight()
		v1, v2 := obj.GetVelocity(), other.GetVelocity()
		m1, m2 := obj.GetMass(), other.GetMass()
		if p1.X+w1/2 < p2.X-w2/2 || p1.X-w1/2 > p2.X+w2/2 || p1.Y+h1/2 < p2.Y-h2/2 || p1.Y-h1/2 > p2.Y+h2/2 {
			// No collision
			return nil
		} else {
			// Collision
			if obj.GetIsPassthrough() || other.GetIsPassthrough() {
				return other
			}

			overlapX := math.Min(p1.X+w1/2, p2.X+w2/2) - math.Max(p1.X-w1/2, p2.X-w2/2)
			overlapY := math.Min(p1.Y+h1/2, p2.Y+h2/2) - math.Max(p1.Y-h1/2, p2.Y-h2/2)
			u1 := v1.MulScalar(m1 - m2).Add(v2.MulScalar(2 * m2)).DivScalar(m1 + m2)
			u2 := v1.MulScalar(2 * m1).Sub(v2.MulScalar(m1 - m2)).DivScalar(m1 + m2)

			if overlapX > overlapY {
				obj.SetVelocity(Vector2{X: v1.X, Y: u1.Y})
				other.SetVelocity(Vector2{X: v2.X, Y: u2.Y})
				if p1.Y < p2.Y {
					// Collision on top
					obj.SetCenter(Vector2{X: p1.X, Y: p1.Y - overlapY})
				} else {
					// Collision on bottom
					obj.SetCenter(Vector2{X: p1.X, Y: p1.Y + overlapY})
				}
			} else {
				obj.SetVelocity(Vector2{X: u1.X, Y: v1.Y})
				other.SetVelocity(Vector2{X: u2.X, Y: v2.Y})
				if p1.X < p2.X {
					// Collision on left
					obj.SetCenter(Vector2{X: p1.X - overlapX, Y: p1.Y})

				} else {
					// Collision on right
					obj.SetCenter(Vector2{X: p1.X + overlapX, Y: p1.Y})

				}
			}
			return other
		}
	}
	return nil
}

func CheckStaticCollision(obj DynamicObject, other StaticObject) StaticObject {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		r1 := obj.GetWidth() / 2
		r2 := other.GetWidth() / 2
		dist := p1.Distance(p2)
		if dist <= r1+r2 {
			if obj.GetIsPassthrough() {
				return other
			}

			n := p1.Sub(p2).Normalize()
			v1 := obj.GetVelocity()
			u1 := n.MulScalar(v1.Length())
			obj.SetVelocity(u1)
			penetrationDepth := r1 + r2 - p1.Distance(p2)
			direction := p1.Sub(p2).Normalize()
			obj.SetCenter(p1.Add(direction.MulScalar(penetrationDepth)))
			return other
		}
	}
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Rectangle {
		circle, rect := obj, other

		NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
		NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
		dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

		penetrationDepth := circle.GetWidth()/2 - dist.Length()
		penetrationVector := dist.Normalize().MulScalar(penetrationDepth)

		if penetrationDepth > 0.0 {
			if obj.GetIsPassthrough() {
				return other
			}

			if circle.GetVelocity().Dot(dist) < 0 {
				tangentVel := dist.Normalize().Dot(circle.GetVelocity())
				circle.SetVelocity(circle.GetVelocity().Sub(dist.Normalize().MulScalar(tangentVel * 2)))
			}
			circle.SetCenter(circle.GetCenter().Add(penetrationVector))
			return other
		}
	}
	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Circle {
		rect, circle := obj, other

		NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
		NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
		dist := Vector2{X: NearestX - circle.GetCenter().X, Y: NearestY - circle.GetCenter().Y}

		penetrationDepth := circle.GetWidth()/2 - dist.Length()
		penetrationVector := dist.Normalize().MulScalar(penetrationDepth)

		if penetrationDepth > 0.0 {
			if obj.GetIsPassthrough() {
				return other
			}

			tangentVel := dist.Normalize().Dot(rect.GetVelocity())
			rect.SetVelocity(rect.GetVelocity().Sub(dist.Normalize().MulScalar(tangentVel * 2)))
			rect.SetCenter(rect.GetCenter().Add(penetrationVector))
			return other
		}
	}
	if obj.GetCollisionType() == Rectangle && other.GetCollisionType() == Rectangle {
		p1, p2 := obj.GetCenter(), other.GetCenter()
		w1, w2 := obj.GetWidth(), other.GetWidth()
		h1, h2 := obj.GetHeight(), other.GetHeight()
		v1 := obj.GetVelocity()
		if p1.X+w1/2 < p2.X-w2/2 || p1.X-w1/2 > p2.X+w2/2 || p1.Y+h1/2 < p2.Y-h2/2 || p1.Y-h1/2 > p2.Y+h2/2 {
			// No collision
			return nil
		} else {
			// Collision
			if obj.GetIsPassthrough() {
				return other
			}

			overlapX := math.Min(p1.X+w1/2, p2.X+w2/2) - math.Max(p1.X-w1/2, p2.X-w2/2)
			overlapY := math.Min(p1.Y+h1/2, p2.Y+h2/2) - math.Max(p1.Y-h1/2, p2.Y-h2/2)

			if overlapX > overlapY {
				obj.SetVelocity(Vector2{X: v1.X, Y: -v1.Y})
				if p1.Y < p2.Y {
					// Collision on top
					obj.SetCenter(Vector2{X: p1.X, Y: p1.Y - overlapY})
				} else {
					// Collision on bottom
					obj.SetCenter(Vector2{X: p1.X, Y: p1.Y + overlapY})
				}
			} else {
				obj.SetVelocity(Vector2{X: -v1.X, Y: v1.Y})
				if p1.X < p2.X {
					// Collision on left
					obj.SetCenter(Vector2{X: p1.X - overlapX, Y: p1.Y})

				} else {
					// Collision on right
					obj.SetCenter(Vector2{X: p1.X + overlapX, Y: p1.Y})
				}
			}
			return other
		}
	}
	return nil
}

func (g *gameObject) onDynamicCollisionCircles(collider DynamicObject, dt float64) {
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
func (g *gameObject) onDynamicCollisionRectangles(collider DynamicObject, dt float64) {
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
func (g *gameObject) onDynamicCollisionCircleRectangle(collider DynamicObject, dt float64) {
	circle, rect := g, collider
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
func (g *gameObject) onDynamicCollisionRectangleCircle(collider DynamicObject, dt float64) {
	rect, circle := g, collider
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
func (g *gameObject) onStaticCollisionCircles(collider StaticObject, dt float64) {
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
func (g *gameObject) onStaticCollisionRectangles(collider StaticObject, dt float64) {
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
func (g *gameObject) onStaticCollisionCircleRectangle(collider StaticObject, dt float64) {
	circle, rect := g, collider
	if g.GetIsPassthrough() {
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
func (g *gameObject) onStaticCollisionRectangleCircle(collider StaticObject, dt float64) {
	rect, circle := g, collider
	if g.GetIsPassthrough() {
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

func (e *Environment) CircleCollidesWithCircle(obj DynamicObject, other Object) bool {
	p1 := obj.GetCenter()
	p2 := other.GetCenter()
	r1 := obj.GetWidth() / 2
	r2 := other.GetWidth() / 2
	dist := p1.Distance(p2)
	return dist <= r1+r2
}
func (e *Environment) CircleCollidesWithRectangle(obj DynamicObject, other Object) bool {
	circle, rect := obj, other

	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()

	return penetrationDepth > 0.0
}
func (e *Environment) RectangleCollidesWithCircle(obj DynamicObject, other Object) bool {
	rect, circle := obj, other

	NearestX := math.Max(rect.GetCenter().X-rect.GetWidth()/2, math.Min(circle.GetCenter().X, rect.GetCenter().X+rect.GetWidth()/2))
	NearestY := math.Max(rect.GetCenter().Y-rect.GetHeight()/2, math.Min(circle.GetCenter().Y, rect.GetCenter().Y+rect.GetHeight()/2))
	dist := Vector2{X: circle.GetCenter().X - NearestX, Y: circle.GetCenter().Y - NearestY}

	penetrationDepth := circle.GetWidth()/2 - dist.Length()

	return penetrationDepth > 0.0

}
func (e *Environment) RectangleCollidesWithRectangle(obj DynamicObject, other Object) bool {
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
