package game

import (
	"bnd/engine"
	"time"
)

type fartObject struct {
	engine.DynamicObject
	hasCollided      bool
	disableCollision bool
}

func newFartObject(player engine.DynamicObject, width, height float64) *fartObject {
	return &fartObject{
		DynamicObject: engine.NewDynamicObject(engine.GameObjectConf{
			CollisionType: engine.Circle,
			Center:        player.GetCenter(),
			Width:         width,
			Height:        height,
			IsPassthrough: true,
			Until:         time.Now().Add(100 * time.Millisecond),
		}),
	}
}

func (f *fartObject) OnCollision(collider engine.Object, dt float64) {
	if f.disableCollision {
		return
	}
	pushVector := collider.GetCenter().Sub(f.GetCenter()).Normalize().DivScalar(dt)
	if u, ok := collider.(engine.DynamicObject); ok {
		u.AddForce(pushVector)
	}
	f.hasCollided = true
}

func (f *fartObject) Update(dt float64) {
	if f.hasCollided {
		f.disableCollision = true
	}
}
