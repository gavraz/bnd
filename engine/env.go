package engine

import "fmt"

type Environment struct {
	dynamicObjects map[string]DynamicObject
	staticObjects  map[string]StaticObject
}

func NewEnvironment() *Environment {
	return &Environment{
		dynamicObjects: make(map[string]DynamicObject),
		staticObjects:  make(map[string]StaticObject),
	}
}

func (e *Environment) AddDynamicObject(name string, object DynamicObject) {
	e.dynamicObjects[name] = object
}

func (e *Environment) AddStaticObject(name string, object StaticObject) {
	e.staticObjects[name] = object
}

func (e *Environment) ForEachGameObject(do func(object Object)) {
	for _, obj := range e.dynamicObjects {
		do(obj)
		obj.ForEachChild(do)
	}
	for _, obj := range e.staticObjects {
		do(obj)
	}
}

type DynamicObjectCollisionMap map[DynamicObject][]DynamicObject
type StaticObjectCollisionMap map[DynamicObject][]StaticObject

func (e *Environment) ResolveDynamicCollisions(obj DynamicObject) DynamicObjectCollisionMap {
	colliders := make(DynamicObjectCollisionMap)
	for _, child := range obj.GetChildren() {
		colliders[child] = append(colliders[child], e.ResolveDynamicCollisions(child)[child]...)
	}

	for _, other := range e.dynamicObjects {
		if obj == other || obj.GetRootParent() == other.GetRootParent() {
			continue
		}
		if e.CollidesWith(obj, other) {
			colliders[obj] = append(colliders[obj], other)
		}
	}
	return colliders
}

func (e *Environment) ResolveStaticCollisions(obj DynamicObject) StaticObjectCollisionMap {
	colliders := make(StaticObjectCollisionMap)
	for _, child := range obj.GetChildren() {
		colliders[child] = append(colliders[child], e.ResolveStaticCollisions(child)[child]...)
	}

	for _, other := range e.staticObjects {
		if obj == other {
			continue
		}
		if e.CollidesWith(obj, other) {
			colliders[obj] = append(colliders[obj], other)
		}
	}
	return colliders
}

func (e *Environment) ObjectByName(name string) Object {
	if obj, ok := e.dynamicObjects[name]; ok {
		return obj
	}

	if obj, ok := e.staticObjects[name]; ok {
		return obj
	}
	return nil
}

func (e *Environment) Update(dt float64) {
	for _, obj := range e.dynamicObjects {
		if obj.GetParent() != nil {
			continue
		}
		obj.ApplyFriction(dt)
		obj.Update(dt)

		if colliders := e.ResolveDynamicCollisions(obj); colliders != nil {
			for collider, collidee := range colliders {
				for _, c := range collidee {
					collider.OnCollision(c, dt)
					fmt.Println("Dynamic Collision detected: ", collider, c)
				}
			}
		}
		if colliders := e.ResolveStaticCollisions(obj); colliders != nil {
			for collider, collidee := range colliders {
				for _, c := range collidee {
					collider.OnCollision(c, dt)
					fmt.Println("Static Collision detected: ", collider, c)
				}
			}
		}

		// TODO: melee collision
		//Might need to move it into a whole outside function that deals with such object types in the future
		//do := func(child Object) {
		//	if ObjectType(child) != Melee {
		//		return
		//	}
		//	if collider := e.ResolveDynamicCollisions(child.(DynamicObject)); collider != nil && collider != child.(DynamicObject).GetParent() && ObjectType(collider) != Melee {
		//		collider.(DynamicObject).GetHit()
		//	}
		//}
		//obj.ForEachChild(do)
	}
	for _, obj := range e.dynamicObjects {
		if obj.GetParent() != nil {
			continue
		}
		obj.Update(dt)
	}
}

func (e *Environment) ClearGameData() {
	e.dynamicObjects = make(map[string]DynamicObject)
	e.staticObjects = make(map[string]StaticObject)
}
