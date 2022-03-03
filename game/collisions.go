package game

import "math"

func CheckDynamicCollision(obj, other DynamicObject) Object {
	if obj.GetCollisionType() == Circle && other.GetCollisionType() == Circle {
		p1 := obj.GetCenter()
		p2 := other.GetCenter()
		r1 := obj.GetWidth() / 2
		r2 := other.GetWidth() / 2
		dist := p1.Distance(p2)
		if dist <= r1+r2 {
			n := p1.Sub(p2).Normalize()
			v1 := obj.GetVelocity()
			v2 := other.GetVelocity()
			a1 := v1.Dot(n)
			a2 := v2.Dot(n)
			// Using the optimized version,
			// optimizedP =  2(a1 - a2)
			//              -----------
			//                m1 + m2
			optimizedP := (2.0 * (a1 - a2)) / (obj.GetMass() + other.GetMass())
			u1 := v1.Sub(n.MulScalar(optimizedP * other.GetMass()))
			u2 := v2.Add(n.MulScalar(optimizedP * obj.GetMass()))
			obj.SetVelocity(u1)
			other.SetVelocity(u2)
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
		// TODO: Implement rectangle-rectangle collision
	}
	return nil
}

func CheckStaticCollision(obj DynamicObject, other StaticObject) Object {
	// TODO: Implement static-dynamic collision
	return nil
}
