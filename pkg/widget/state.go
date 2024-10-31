package widget

// BaseState provides a default implementation of State
type BaseState struct {
	widget  StatefulWidget
	element StatefulElement
	context BuildContext
}

func (s *BaseState) InitState()               {}
func (s *BaseState) Dispose()                 {}
func (s *BaseState) Widget() StatefulWidget   { return s.widget }
func (s *BaseState) Element() StatefulElement { return s.element }
func (s *BaseState) Context() BuildContext    { return s.context }

func (s *BaseState) MountState(element StatefulElement) {
	s.element = element
	s.widget = element.Widget().(StatefulWidget)
	s.context = element.BuildContext()
	s.InitState()
}

func (s *BaseState) SetState(fn func()) {
	if fn != nil {
		fn()
	}
	s.element.MarkNeedsBuild()
}
