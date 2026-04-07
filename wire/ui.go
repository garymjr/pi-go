package wire

type UIRequest interface {
	Frame
	uiMethod() string
}

type uiRequestBase struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Method string `json:"method"`
}

func (r uiRequestBase) frameType() string { return r.Type }
func (r uiRequestBase) uiMethod() string  { return r.Method }

type SelectUIRequest struct {
	uiRequestBase
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Timeout *int     `json:"timeout,omitempty"`
}

type ConfirmUIRequest struct {
	uiRequestBase
	Title   string `json:"title"`
	Message string `json:"message"`
	Timeout *int   `json:"timeout,omitempty"`
}

type InputUIRequest struct {
	uiRequestBase
	Title       string `json:"title"`
	Placeholder string `json:"placeholder,omitempty"`
	Timeout     *int   `json:"timeout,omitempty"`
}

type EditorUIRequest struct {
	uiRequestBase
	Title   string `json:"title"`
	Prefill string `json:"prefill,omitempty"`
}

type NotifyUIRequest struct {
	uiRequestBase
	Message    string `json:"message"`
	NotifyType string `json:"notifyType,omitempty"`
}

type SetStatusUIRequest struct {
	uiRequestBase
	StatusKey  string `json:"statusKey"`
	StatusText string `json:"statusText,omitempty"`
}

type SetWidgetUIRequest struct {
	uiRequestBase
	WidgetKey       string   `json:"widgetKey"`
	WidgetLines     []string `json:"widgetLines,omitempty"`
	WidgetPlacement string   `json:"widgetPlacement,omitempty"`
}

type SetTitleUIRequest struct {
	uiRequestBase
	Title string `json:"title"`
}

type SetEditorTextUIRequest struct {
	uiRequestBase
	Text string `json:"text"`
}
