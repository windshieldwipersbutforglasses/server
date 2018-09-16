package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/dragonboard"
)

var step1 *gpio.StepperDriver
var step2 *gpio.StepperDriver

func main() {
	drag := dragonboard.NewAdaptor()
	drag.Connect()
	pins1 := [4]string{"a0", "a1", "a2", "a3"}
	pins2 := [4]string{"0", "1", "2", "3"}

	step1 = gpio.NewStepperDriver(drag, pins1, gpio.StepperModes.SinglePhaseStepping, 10)
	step2 = gpio.NewStepperDriver(drag, pins2, gpio.StepperModes.SinglePhaseStepping, 10)

	for i := 0; i < 5; i++ {
		step1.Move(10)
		step2.Move(10)
		time.Sleep(1 * time.Second)
		step1.Move(-10)
		step2.Move(-10)
		time.Sleep(1 * time.Second)
	}

	s := NewServer()
	log.Fatal(s.Start())
}

type Server struct {
	router *mux.Router
}

func NewServer() *Server {
	r := mux.NewRouter()
	return &Server{
		router: r,
	}
}

func (s *Server) bindHandlers() {
	s.router.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
		for i := 0; i < 5; i++ {
			step1.Move(10)
			step2.Move(10)
			time.Sleep(1 * time.Second)
			step1.Move(-10)
			step2.Move(-10)
			time.Sleep(1 * time.Second)
		}
	})
}

func (s *Server) Start() error {
	s.bindHandlers()
	return http.ListenAndServe(":3000", s.router)
}
