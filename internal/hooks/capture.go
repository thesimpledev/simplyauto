//go:build windows

package hooks

import (
	"context"
	"sync"
	"time"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"

	"simplyauto/pkg/events"
)

type EventCallback func(event events.InputEvent)

type InputCapture struct {
	callback     EventCallback
	startTime    time.Time
	filterKeys   map[Key]bool
	ctx          context.Context
	cancel       context.CancelFunc
	running      bool
	captureMouse bool
	captureKbd   bool
	mu           sync.Mutex
}

type CaptureOptions struct {
	CaptureMouse    bool
	CaptureKeyboard bool
	FilterKeys      []Key
}

func NewInputCapture(opts CaptureOptions) *InputCapture {
	filterKeys := make(map[Key]bool)
	for _, k := range opts.FilterKeys {
		filterKeys[k] = true
	}

	return &InputCapture{
		filterKeys:   filterKeys,
		captureMouse: opts.CaptureMouse,
		captureKbd:   opts.CaptureKeyboard,
	}
}

func (c *InputCapture) Start(callback EventCallback) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	c.callback = callback
	c.startTime = time.Now()
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.running = true

	if c.captureKbd {
		kbdChan := make(chan types.KeyboardEvent, 100)
		if err := keyboard.Install(nil, kbdChan); err != nil {
			return err
		}
		go c.keyboardLoop(kbdChan)
	}

	if c.captureMouse {
		mouseChan := make(chan types.MouseEvent, 100)
		if err := mouse.Install(nil, mouseChan); err != nil {
			if c.captureKbd {
				keyboard.Uninstall()
			}
			return err
		}
		go c.mouseLoop(mouseChan)
	}

	return nil
}

func (c *InputCapture) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return nil
	}

	c.cancel()
	c.running = false

	if c.captureKbd {
		keyboard.Uninstall()
	}
	if c.captureMouse {
		mouse.Uninstall()
	}

	return nil
}

func (c *InputCapture) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

func (c *InputCapture) keyboardLoop(ch chan types.KeyboardEvent) {
	for {
		select {
		case <-c.ctx.Done():
			return
		case evt := <-ch:
			if c.filterKeys[Key(evt.VKCode)] {
				continue
			}

			var eventType events.EventType
			switch evt.Message {
			case types.WM_KEYDOWN, types.WM_SYSKEYDOWN:
				eventType = events.EventKeyDown
			case types.WM_KEYUP, types.WM_SYSKEYUP:
				eventType = events.EventKeyUp
			default:
				continue
			}

			c.callback(events.InputEvent{
				Type:      eventType,
				Timestamp: time.Since(c.startTime),
				KeyCode:   uint16(evt.VKCode),
				ScanCode:  uint16(evt.ScanCode),
			})
		}
	}
}

func (c *InputCapture) mouseLoop(ch chan types.MouseEvent) {
	for {
		select {
		case <-c.ctx.Done():
			return
		case evt := <-ch:
			var eventType events.EventType

			switch evt.Message {
			case types.WM_MOUSEMOVE:
				eventType = events.EventMouseMove
			case types.WM_LBUTTONDOWN:
				eventType = events.EventMouseLeftDown
			case types.WM_LBUTTONUP:
				eventType = events.EventMouseLeftUp
			case types.WM_RBUTTONDOWN:
				eventType = events.EventMouseRightDown
			case types.WM_RBUTTONUP:
				eventType = events.EventMouseRightUp
			case types.WM_MBUTTONDOWN:
				eventType = events.EventMouseMiddleDown
			case types.WM_MBUTTONUP:
				eventType = events.EventMouseMiddleUp
			case types.WM_MOUSEWHEEL:
				eventType = events.EventMouseWheel
			default:
				continue
			}

			inputEvent := events.InputEvent{
				Type:      eventType,
				Timestamp: time.Since(c.startTime),
				X:         int(evt.Point.X),
				Y:         int(evt.Point.Y),
			}

			if eventType == events.EventMouseWheel {
				inputEvent.Delta = int(int16(evt.MouseData >> 16))
			}

			c.callback(inputEvent)
		}
	}
}
