package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

var (
	ctxTyp = reflect.TypeOf((*gin.Context)(nil))
	errTyp = reflect.TypeOf((*error)(nil)).Elem()

	// Binder iface
	binderTyp = reflect.TypeOf((*interface {
		Bind(ctx *gin.Context) error
	})(nil)).Elem()

	// validator iface
	validatorTyp = reflect.TypeOf((*interface {
		Validate() error
	})(nil)).Elem()
)

func wrapper(fun interface{}, r Renderer) gin.HandlerFunc {
	ftyp := reflect.TypeOf(fun)
	if !isFunc(ftyp) {
		panic(fmt.Errorf("handler isn't a func"))
	}

	sigErr := fmt.Errorf("illegal signature for handler: %s", funcDesc(ftyp))

	ins := funIn(ftyp)
	outs := funOut(ftyp)

	cin := len(ins)
	cout := len(outs)

	// the handler must have parameters or return-value
	assert(cin <= 2 && cout <= 2 && !(cin == 0 && cout == 0), sigErr)

	// it's a gin.HandlerFunc
	if cin == 1 && cout == 0 && ins[0] == ctxTyp {
		return fun.(func(ctx *gin.Context))
	}

	hasCtx := false // the first param is gin.Context

	needBind := false // should execute ctx.Bind automatic

	implBinder := false
	implValidator := false

	if cin > 0 {
		if ins[0] == ctxTyp {
			hasCtx = true
		} else {
			p1 := ins[0]

			assert(p1.Kind() == reflect.Ptr && p1.Elem().Kind() == reflect.Struct, sigErr)
			needBind = true
		}
	}

	// if we don't need the cxt parameter, we should have at least a return-value
	assert(hasCtx || cout > 0, sigErr)

	if cin == 2 {
		p2 := ins[1]
		assert(hasCtx && p2 != ctxTyp, sigErr)
		assert(p2.Kind() == reflect.Ptr && p2.Elem().Kind() == reflect.Struct, sigErr)
		needBind = true
	}

	if needBind {
		i := 0
		if hasCtx {
			i = 1
		}

		p := ins[i]
		if p.Implements(binderTyp) {
			implBinder = true
		}

		if p.Implements(validatorTyp) {
			implValidator = true
		}
	}

	hasError := false
	needRender := false
	if cout > 0 {
		r1 := outs[0]
		if r1.Implements(errTyp) {
			// the error must be the last return-value
			assert(cout == 1, sigErr)
			hasError = true
		} else {
			needRender = true
		}
	}

	if cout == 2 {
		r2 := outs[1]
		assert(!hasError && r2.Implements(errTyp), sigErr)
		hasError = true
	}

	// wrap our handler to a gin.HandlerFunc
	return func(ctx *gin.Context) {
		inParams := make([]reflect.Value, 0, cin)

		if hasCtx {
			inParams = append(inParams, reflect.ValueOf(ctx))
		}

		if needBind {
			typ := ins[0]
			if hasCtx {
				typ = ins[1]
			}

			pv := reflect.New(typ.Elem())
			var err error

			if implBinder {
				// if impl Binder iface, call Bind to bind parameters
				err = pv.Interface().(interface {
					Bind(ctx *gin.Context) error
				}).Bind(ctx)
			} else {
				// use ctx.Bind default
				err = ctx.Bind(pv.Interface())
				if err != nil {
					err = BindError{err}
				}
			}

			if err == nil && implValidator {
				err = pv.Interface().(interface {
					Validate() error
				}).Validate()
			}

			if err != nil {
				r.Render(ctx, err)
				return
			}

			inParams = append(inParams, pv)
		}

		// call original handler
		outs := reflect.ValueOf(fun).Call(inParams)

		if hasError {
			idx := 0
			if needRender {
				idx = 1
			}
			err := outs[idx]
			if !err.IsNil() {
				r.Render(ctx, err.Interface())
				return
			}
		}

		if needRender {
			if out := outs[0].Interface(); out != nil {
				r.Render(ctx, out)
			}
			return
		}

		if !ctx.Writer.Written() {
			r.Render(ctx, nil)
			return
		}
	}
}

func isFunc(typ reflect.Type) bool {
	return typ.Kind() == reflect.Func
}

func funcDesc(typ reflect.Type) string {
	if typ.Kind() != reflect.Func {
		return ""
	}
	return typ.String()
}

func funIn(typ reflect.Type) []reflect.Type {
	ins := make([]reflect.Type, typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		ins[i] = typ.In(i)
	}
	return ins
}

func funOut(typ reflect.Type) []reflect.Type {
	outs := make([]reflect.Type, typ.NumOut())
	for i := 0; i < typ.NumOut(); i++ {
		outs[i] = typ.Out(i)
	}
	return outs
}

func assert(b bool, err error) {
	if !b {
		panic(err)
	}
}

