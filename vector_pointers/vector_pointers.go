package vectorpointers

import "math"

type Vector2 struct {
	X float64
	Y float64
}

func NewVector2(x, y float64) *Vector2 {
	return &Vector2{
		X: x,
		Y: y,
	}
}

func (v *Vector2) Add(other *Vector2) *Vector2 {
	v.X += other.X
	v.Y += other.Y
	return v
}

func (v *Vector2) Sub(other *Vector2) *Vector2 {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v *Vector2) Mul(other Vector2) *Vector2 {
	v.X *= other.X
	v.Y *= other.Y
	return v
}

func (v *Vector2) MulScalar(scale float64) *Vector2 {
	v.X *= scale
	v.Y *= scale
	return v
}

func (v *Vector2) Div(other Vector2) *Vector2 {
	v.X /= other.X
	v.Y /= other.Y
	return v
}

func (v *Vector2) DivScalar(scale float64) *Vector2 {
	v.X /= scale
	v.Y /= scale
	return v
}

func (v *Vector2) Squared() *Vector2 {
	v.X = v.X * v.X
	v.Y = v.Y * v.Y
	return v
}

func (v *Vector2) Dot(other *Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v *Vector2) Normalize() *Vector2 {
	length := v.Length()
	if length == 0 {
		return v
	}
	return v.DivScalar(length)
}

func (v *Vector2) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2) Distance(other *Vector2) float64 {
	return v.Sub(other).Length()
}

func (v *Vector2) DistanceSquared(other *Vector2) float64 {
	return v.Sub(other).LengthSquared()
}

func (v *Vector2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v *Vector2) AngleTo(other *Vector2) float64 {
	return math.Atan2(other.Y-v.Y, other.X-v.X)
}

func (v *Vector2) Rotate(angle float64) *Vector2 {
	v.X = v.X*math.Cos(angle) - v.Y*math.Sin(angle)
	v.Y = v.X*math.Sin(angle) + v.Y*math.Cos(angle)
	return v
}
