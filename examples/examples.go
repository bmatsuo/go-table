package examples

import "fmt"

type sprinter struct {
	flag string
	v    interface{}
}

func newSprinter(flag string, v interface{}) *sprinter { return &sprinter{flag, v} }

func (s *sprinter) String() string { return fmt.Sprintf("fmt.Sprintf(%q, %v)", s.flag, s.v) }
func (s *sprinter) Out() string    { return fmt.Sprintf(s.flag, s.v) }
