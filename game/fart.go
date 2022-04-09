package game

import (
	"bnd/engine"
	"time"
)

type fartObject struct {
	engine.DynamicObject
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

func (f *fartObject) Update(dt float64) {

}
