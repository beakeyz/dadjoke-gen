package routing

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

type Router interface{
  Handle(_func, patern string, callbacks []web.Handler)
  Get(patern string, callbacks ...web.Handler)
}

type RouteHandler interface{
  RegisterGet(string, ...web.Handler)

  RegisterPost(string, ...web.Handler)

  RegisterPut(string, ...web.Handler)

  RegisterPatch(string, ...web.Handler)

  RegisterDelete(string, ...web.Handler)

  ArmRoutes(Router)

  Clear()
}

func CraftRegister() *RouteHandlerStruct {
  return &RouteHandlerStruct{
    name: "",
    routes: []route{},
  }
}

type route struct{
  httpmethod string
  identifier string
  handlers []web.Handler
}

type RouteHandlerStruct  struct{
  name string
  routes []route
}

func (self *RouteHandlerStruct) Clear(){
  if self != nil{
    self.routes = nil
  }
}

func (self *RouteHandlerStruct)ArmRoutes(router Router){
  for _, r := range self.routes{

    if (r.httpmethod == http.MethodGet){
      router.Get(r.identifier, r.handlers...)
    }else{
      router.Handle(r.httpmethod, r.identifier, r.handlers)
    }
  }
}

func (self *RouteHandlerStruct) unwrapedRouting(identifier string, httpmethod string, handlers ...web.Handler){
  usedHandlers := make([]web.Handler, 0)
  fmt.Println("Hit unwrapped")
  fmt.Println(reflect.TypeOf(handlers))

  //TODO: add middlewares

  //

  

  usedHandlers = append(usedHandlers, handlers...)

  self.routes = append(self.routes, route{
    httpmethod: httpmethod,
    identifier: self.name + identifier,
    handlers: usedHandlers,
  })
}

func addHandler(a []web.Handler, index int, handlerToAdd web.Handler) []web.Handler {
  if len(a) == index{
    //no need for a funny action
    return append(a, handlerToAdd)
  }
  //shift all values
  a = append(a[:index+1], a[index:]...)
  //insert handler into its position
  a[index] = handlerToAdd
  //return new list
  return a
}

func (self *RouteHandlerStruct) RegisterGet(route string, handler ...web.Handler){
  self.unwrapedRouting(route, http.MethodGet, handler...)
}
func (self *RouteHandlerStruct) RegisterPost(route string, handler ...web.Handler){
  self.unwrapedRouting(route, http.MethodPost, handler...)
}
func (self *RouteHandlerStruct) RegisterDelete(route string, handler ...web.Handler){
  self.unwrapedRouting(route, http.MethodDelete, handler...)
}
func (self *RouteHandlerStruct) RegisterPatch(route string, handler ...web.Handler){
  self.unwrapedRouting(route, http.MethodPatch, handler...)
}
func (self *RouteHandlerStruct) RegisterPut(route string, handler ...web.Handler){
  self.unwrapedRouting(route, http.MethodPut, handler...)
}

//put the other impls here

