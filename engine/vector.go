package engine

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

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vector2) AddScalar(scale float64) Vector2 {
	return Vector2{
		X: v.X + scale,
		Y: v.Y + scale,
	}
}

func (v Vector2) Sub(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vector2) SubScalar(scale float64) Vector2 {
	return Vector2{
		X: v.X - scale,
		Y: v.Y - scale,
	}
}

func (v Vector2) Mul(other Vector2) Vector2 {
	return Vector2{
		X: v.X * other.X,
		Y: v.Y * other.Y,
	}
}

func (v Vector2) MulScalar(scale float64) Vector2 {
	return Vector2{
		X: v.X * scale,
		Y: v.Y * scale,
	}
}

func (v Vector2) Div(other Vector2) Vector2 {
	return Vector2{
		X: v.X / other.X,
		Y: v.Y / other.Y,
	}
}

func (v Vector2) DivScalar(scale float64) Vector2 {
	return Vector2{
		X: v.X / scale,
		Y: v.Y / scale,
	}
}

func (v Vector2) Squared() Vector2 {
	return Vector2{
		X: v.X * v.X,
		Y: v.Y * v.Y,
	}
}

func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector2) Normalize() Vector2 {
	return v.DivScalar(v.Length())
}

func (v Vector2) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector2) LengthSquared() float64 {
	return v.Dot(v)
}

func (v Vector2) Distance(other Vector2) float64 {
	return v.Sub(other).Length()
}

func (v Vector2) DistanceSquared(other Vector2) float64 {
	return v.Sub(other).LengthSquared()
}

func (v Vector2) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func (v Vector2) AngleTo(other Vector2) float64 {
	return math.Atan2(other.Y, other.X) - math.Atan2(v.Y, v.X)
}

func (v Vector2) Rotate(angle float64) Vector2 {
	return Vector2{
		X: v.X*math.Cos(angle) - v.Y*math.Sin(angle),
		Y: v.X*math.Sin(angle) + v.Y*math.Cos(angle),
	}
}
