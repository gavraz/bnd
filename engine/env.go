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

func (e *Environment) ResolveDynamicCollisions(obj DynamicObject) DynamicObject {
	var isChild bool
	colliders := make([]DynamicObject, 0)
	for _, other := range e.dynamicObjects {
		isChild = false
		if other == obj {
			continue
		}
		for _, child := range obj.GetChildren() {
			if Object(other) == child {
				isChild = true
				break
			}
		}
		if isChild {
			continue
		}
		for _, child := range other.GetChildren() {
			if Object(obj) == child {
				isChild = true
				break
			}
		}
		if isChild {
			continue
		}
		if collider := CheckDynamicCollision(obj, other); collider != nil {
			colliders = append(colliders, collider)
		}
	}
	return nil
}

func (e *Environment) ResolveStaticCollisions(obj DynamicObject) StaticObject {
	var isChild bool
	for _, other := range e.staticObjects {
		isChild = false
		for _, child := range obj.GetChildren() {
			if Object(other) == child {
				isChild = true
				break
			}
		}
		if isChild {
			continue
		}
		if collider := CheckStaticCollision(obj, other); collider != nil {
			return collider
		}
	}
	return nil
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

		if collider := e.ResolveDynamicCollisions(obj); collider != nil {
			fmt.Println("Dynamic Collision detected: ", obj.GetCenter(), collider.GetCenter())
		}
		if collider := e.ResolveStaticCollisions(obj); collider != nil {
			fmt.Println("Static Collision detected: ", obj.GetCenter(), collider.GetCenter())
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
}

func (e *Environment) ClearGameData() {
	e.dynamicObjects = make(map[string]DynamicObject)
	e.staticObjects = make(map[string]StaticObject)
}
