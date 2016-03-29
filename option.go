// Package composite is an higher level wrapper of gopherjs-vue,
// by providing a more gopher friendly API, this package tries to hide
// the JavaScript details for VueJS easy usage in GopherJS world.
package vue

import (
	"github.com/gopherjs/gopherjs/js"
)

type LifeCycleEvent string

const (
	EvtInit          LifeCycleEvent = "init"
	EvtCreated       LifeCycleEvent = "created"
	EvtBeforeCompile LifeCycleEvent = "beforeCompile"
	EvtCompiled      LifeCycleEvent = "compiled"
	EvtReady         LifeCycleEvent = "ready"
	EvtAttached      LifeCycleEvent = "attached"
	EvtDetached      LifeCycleEvent = "detached"
	EvtBeforeDestroy LifeCycleEvent = "beforeDestroy"
	EvtDestroyed     LifeCycleEvent = "destroyed"
)

// Option are used to organize mutiple sub component together to
// construct VueJS apps or (higher level) components.
type Option struct {
	*js.Object

	// Type: String
	// Restriction: only respected when used in Vue.extend().
	// Details:
	//
	// Allow the component to recursively invoke itself in its template.
	// Note that when a component is registered globally with
	// Vue.component(), the global ID is automatically set as its name.
	//
	// Another benefit of specifying a name option is console inspection.
	// When inspecting an extended Vue component in the console,
	// the default constructor name is VueComponent,
	// which isn’t very informative. By passing in an optional name option to
	// Vue.extend(), you will get a better inspection output so that
	// you know which component you are looking at.
	// The string will be camelized and used as the component’s constructor name.
	Name string `js:"name"`
	///////////////////////////// Instance Properties

	// 	Type: Object | Function
	//
	// Restriction: Only accepts Function when used in a component definition.
	//
	// Details:
	//
	// The data object for the Vue instance. Vue.js will recursively convert
	// its properties into getter/setters to make it “reactive”. The object
	// must be plain: native objects, existing getter/setters and prototype
	// properties are ignored. It is not recommended to observe complex
	// objects.
	//
	// Once the instance is created, the original data object can be accessed
	// as vm.$data. The Vue instance also proxies all the properties
	// found on the data object.
	//
	// Properties that start with _ or $ will not be proxied on the Vue
	// instance because they may conflict with Vue’s internal properties and
	// API methods. You will have to access them as vm.$data._property.
	//
	// When defining a component, data must be declared as a function that
	// returns the initial data object, because there will be many instances
	// created using the same definition. If we still use a plain object for
	// data, that same object will be shared by reference across all instance
	// created! By providing a data function, every time a new instance
	// is created, we can simply call it to return a fresh copy of
	// the initial data.
	//
	// If required, a deep clone of the original object can be obtained by
	// passing vm.$data through JSON.parse(JSON.stringify(...)).
	Data interface{} `js:"data"`

	// Type: String | HTMLElement | Function
	// Restriction: only accepts type Function when used in a component definition.
	//
	//Details:
	//
	// Provide the Vue instance an existing DOM element to mount on.
	// It can be a CSS selector string, an actual HTMLElement,
	// or a function that returns an HTMLElement.
	// Note that the provided element merely serves as a mounting point;
	// it will be replaced if a template is also provided,
	// unless replace is set to false. The resolved element will
	// be accessible as vm.$el.
	//
	// When used in Vue.extend, a function must be provided
	// so each instance gets a separately created element.
	// If this option is available at instantiation,
	// the instance will immediately enter compilation;
	// otherwise, the user will have to explicitly call
	// vm.$mount() to manually start the compilation.
	El interface{} `js:"el"`

	// Type: String
	//
	// Details:
	//
	// A string template to be used as the markup for the Vue instance. By
	// default, the template will replace the mounted element. When the replace
	// option is set to false, the template will be inserted into the mounted
	// element instead. In both cases, any existing markup inside the mounted
	// element will be ignored, unless content distribution slots are present
	// in the template.
	//
	// If the string starts with # it will be used as a querySelector and use
	// the selected element’s innerHTML as the template string. This allows the
	// use of the common <script type="x-template"> trick to include templates.
	Template string `js:"template"`

	// 	Type: Boolean
	//
	// Default: true
	//
	// Restriction: only respected if the template option is also present.
	//
	// Details:
	//
	// Determines whether to replace the element being mounted on with the
	// template. If set to false, the template will overwrite the element’s inner
	// content without replacing the element itself.
	Replace bool `js:"replace"`

	// parent
	//
	// Type: Vue instance
	//
	// Details:
	//
	// Specify the parent instance for the instance to be created. Establishes
	// a parent-child relationship between the two. The parent will be
	// accessible as this.$parent for the child, and the child will be pushed
	// into the parent’s $children array.
	Parent *js.Object `js:"parent"`

	// map to sub component
	mcom map[string]*Component
	// map to event handler
	mevt map[string]interface{}
	// properties
	props []string
	// mixins
	mixins []js.M
}

func NewOption() *Option {
	c := &Option{
		Object: js.Global.Get("Object").New(),
	}
	c.mcom = make(map[string]*Component, 0)
	c.mevt = make(map[string]interface{}, 0)
	c.props = []string{}
	c.mixins = []js.M{}
	return c
}

// NewViewModel create the VueJS instance for finally use
// the VueJS instance becomes usable only after this call
func (o *Option) NewViewModel() *ViewModel {
	return newViewModel(vue.New(o.prepare()))
}

func (o *Option) NewComponent() *Component {
	if _, ok := o.El.(string); ok {
		panic("Option.El in component must be a function")
	}
	return newComponent(vue.Call("extend", o.prepare()))
}

// prepare set the proper options into js.Object
func (c *Option) prepare() (opts *js.Object) {
	if len(c.mcom) > 0 {
		c.Set("components", c.mcom)
	}
	if len(c.mevt) > 0 {
		c.Set("events", c.mevt)
	}
	if len(c.props) > 0 {
		c.Set("props", c.props)
	}
	if len(c.mixins) > 0 {
		c.Set("mixins", c.mixins)
	}
	return c.Object
}

// SetDataWithMethods set data and methods of the genereated VueJS instance
// based on `structPtr` and `js.MakeWrapper(structPtr)`
func (c *Option) SetDataWithMethods(structPtr interface{}) *Option {
	c.Set("data", structPtr)
	c.Set("methods", js.MakeWrapper(structPtr))
	return c
}

// func (c *Option) OnLifeCycleEvent(evt LifeCycleEvent, fn interface{}) *Option {
// 	c.addMixin(string(evt), fn)
// 	return c
// }

func (o *Option) OnLifeCycleEvent(evt LifeCycleEvent, fn func(vm *ViewModel)) *Option {
	return o.addMixin(
		string(evt),
		js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			vm := newViewModel(this)
			fn(vm)
			return nil
		}),
	)
}

// The mixins option accepts an array of mixin objects.
// These mixin objects can contain instance options just
// like normal instance objects, and they will be
// merged against the eventual options using the same
// option merging logic in Vue.extend(). e.g.
// If your mixin contains a created hook and
// the component itself also has one, both functions will be called.
//
// Mixin hooks are called in the order they are provided,
// and called before the component’s own hooks.
func (c *Option) Mixin(val js.M) *Option {
	c.mixins = append(c.mixins, val)
	return c
}

func (c *Option) addMixin(name string, val interface{}) *Option {
	return c.Mixin(js.M{
		name: val,
	})
}

// AddComponent add sub component to the genereated VueJS instance (optional)
func (c *Option) AddSubComponent(name string, sub *Component) *Option {
	c.mcom[name] = sub
	return c
}

// On add EventHandler to VueJS-generated component-oriented event
// for cross component message passing
func (c *Option) On(name string, fn interface{}) *Option {
	c.mevt[name] = fn
	return c
}

// AddProp add props to the genereated VueJS instance (optional)
// 	props is a list/hash of attributes that are exposed to accept data from
// 	the parent component. It has a simple Array-based syntax and
// 	an alternative Object-based syntax that allows advanced configurations
// 	such as type checking, custom validation and default values.
func (c *Option) AddProp(name ...string) *Option {
	c.props = append(c.props, name...)
	return c
}
