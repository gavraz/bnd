package engine

import (
	"fmt"
)

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
	rootParent := obj.GetRootParent()
	for _, child := range obj.GetChildren() {
		childColliders := e.ResolveDynamicCollisions(child)
		for childObj, colliderList := range childColliders {
			colliders[childObj] = append(colliders[child], colliderList...)
		}
	}
	for _, other := range e.dynamicObjects {
		if rootParent == other.GetRootParent() {
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
		childColliders := e.ResolveStaticCollisions(child)
		for childObj, colliderList := range childColliders {
			colliders[childObj] = append(colliders[child], colliderList...)
		}
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

type OnCollisioner interface {
	OnCollision(collider Object)
}

func (e *Environment) Update(dt float64) {
	for _, obj := range e.dynamicObjects {
		obj.ApplyFriction(dt)
		obj.Update(dt)
		if dynamicCollisions := e.ResolveDynamicCollisions(obj); dynamicCollisions != nil {
			for collider, collidees := range dynamicCollisions {
				for _, collidee := range collidees {
					collider.onCollision(collidee)
					if c, ok := collider.(OnCollisioner); ok {
						c.OnCollision(collidee)
					}
					//fmt.Println("Dynamic Collision detected: ", collider, collidee)
				}
			}
		}
		if staticCollisions := e.ResolveStaticCollisions(obj); staticCollisions != nil {
			for collider, collidees := range staticCollisions {
				for _, collidee := range collidees {
					collider.onCollision(collidee)
					if c, ok := collider.(OnCollisioner); ok {
						c.OnCollision(collidee)
					}
					fmt.Println("Static Collision detected: ", collider, collidee)
				}
			}
		}
	}
}

func (e *Environment) ClearGameData() {
	e.dynamicObjects = make(map[string]DynamicObject)
	e.staticObjects = make(map[string]StaticObject)
}
