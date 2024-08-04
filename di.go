package tinydi

import (
	"github.com/fatih/structs"
	"log"
	"reflect"
)

// defaultInjector is used whenever a nil argument is passed to the tinydi functions.
var defaultInjector = New()

type Injector struct {
	transients map[reflect.Type]func(i *Injector) any

	singletons         map[reflect.Type]func(i *Injector) any
	singletonInstances map[reflect.Type]any
}

// New returns a new Injector context, which can be used to separate different injector
// contexts.
func New() *Injector {
	injector := &Injector{
		transients: make(map[reflect.Type]func(i *Injector) any),

		singletons:         make(map[reflect.Type]func(i *Injector) any),
		singletonInstances: make(map[reflect.Type]any),
	}

	return injector
}

// Add defines a transient service of type T to be included in the DI context.
// When a service is registered using Add, it will create a new instance every time.
// If you specify a nil value as the injector instance, it will use the default, global one.
func Add[T any](inj *Injector, provider func(i *Injector) *T) {
	if inj == nil {
		inj = defaultInjector
	}

	var t T
	dependencyType := reflect.TypeOf(&t)
	//dependencyType := reflect.TypeOf((*T)(nil)).Elem()
	inj.transients[dependencyType] = func(i *Injector) any {
		return provider(i)
	}
}

// AddSingleton defines a singleton service of type T to be included in the DI context.
// When a service is registered using AddSingleton, it will return the same instance every time.
// If you specify a nil value as the injector instance, it will use the default, global one.
func AddSingleton[T any](inj *Injector, provider func(i *Injector) *T) {
	if inj == nil {
		inj = defaultInjector
	}

	var t T
	dependencyType := reflect.TypeOf(&t)
	//dependencyType := reflect.TypeOf((*T)(nil)).Elem()
	inj.singletons[dependencyType] = func(i *Injector) any {
		return provider(i)
	}
}

// MustGet will provide the dependency, or panic.
// If you specify a nil value as the injector instance, it will use the default, global one.
func MustGet[T any](i *Injector) *T {
	var instance *T
	instanceType := reflect.TypeOf(instance)

	dependency := initializeDependencyTree(i, instanceType)
	if dependency == nil {
		log.Panicf("No interface of type %s could be provided.", instanceType)
		return nil
	}
	instance = dependency.(*T)
	return instance
}

func initializeDependencyTree(i *Injector, instanceType reflect.Type) any {
	if i == nil {
		i = defaultInjector
	}

	dependency := initializeDependency(i, instanceType)
	if dependency != nil {
		initializeDependencyDependencies(i, dependency)
		return dependency
	}

	return nil
}

func initializeDependency(inj *Injector, typeOf reflect.Type) any {
	if initializer, ok := inj.singletons[typeOf]; ok {
		if instance, instOk := inj.singletonInstances[typeOf]; instOk {
			return instance
		}
		result := initializer(inj)
		inj.singletonInstances[typeOf] = result
		return result
	}

	if initializer, ok := inj.transients[typeOf]; ok {
		result := initializer(inj)
		return result
	}

	return nil
}

// initializeDependencyDependencies exists to make sure sub-dependencies are also initialized using the correct context.
func initializeDependencyDependencies[T any](inj *Injector, dep T) {
	fields := structs.Fields(dep)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		if field.Tag("di") != "inject" {
			continue
		}

		fieldType := reflect.ValueOf(field.Value()).Type()
		instance := initializeDependencyTree(inj, fieldType)
		if instance == nil {
			continue
		}

		_ = field.Set(instance)
	}
}
