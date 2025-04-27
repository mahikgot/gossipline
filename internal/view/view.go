package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mahikgot/gossipline/internal/net"
	"github.com/rivo/tview"
)

type ChatApp struct {
	app            *tview.Application
	pages          *tview.Pages
	flex           *tview.Flex
	username       string
	messageChan    chan []byte
	chatViewIndex  int
	chatInputIndex int
}

func NewChatApp() *ChatApp {
	chatApp := &ChatApp{
		app:            tview.NewApplication(),
		pages:          tview.NewPages(),
		flex:           tview.NewFlex(),
		messageChan:    make(chan []byte, 96),
		chatViewIndex:  0,
		chatInputIndex: 1,
	}
	return chatApp
}

func (ca *ChatApp) createUsernameModal() tview.Primitive {
	usernameInput := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorAntiqueWhite).
		SetFieldTextColor(tcell.ColorBlack).
		SetLabel("username:").
		SetLabelColor(tcell.ColorBlanchedAlmond)

	usernameInput.SetDoneFunc(func(key tcell.Key) {
		if len(usernameInput.GetText()) <= 0 {
			return
		}

		ca.username = usernameInput.GetText()
		message := ca.username + " has logged in!"
		ca.sendMessage(message)
		ca.pages.RemovePage("modal")
		ca.app.SetFocus(ca.flex.GetItem(ca.chatInputIndex))
	})

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	return modal(usernameInput, 21, 3)
}

func (ca *ChatApp) createUI() {
	input := tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorAntiqueWhite).
		SetFieldTextColor(tcell.ColorBlack).
		SetLabel("chat:").
		SetLabelColor(tcell.ColorBlanchedAlmond)
	input.SetDoneFunc(func(key tcell.Key) {
		if len(input.GetText()) <= 0 {
			return
		}
		ca.sendMessage(input.GetText())
		input.SetText("")
	})
	chatView := tview.NewTextView()
	ca.flex.SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(input, 1, 1, true)

	usernameModal := ca.createUsernameModal()
	ca.pages.AddPage("background", ca.flex, true, true).
		AddPage("modal", usernameModal, true, true)
}

func (ca *ChatApp) sendMessage(msg string) {
	ca.messageChan <- []byte(msg + "\n")
}

func (ca *ChatApp) consumeMessages() {
	for {
		data := <-ca.messageChan
		dest := ca.flex.GetItem(ca.chatViewIndex).(*tview.TextView)
		fmt.Fprint(dest, string(data))
		dest.ScrollToEnd()
	}
}

func (ca *ChatApp) Start() {
	ca.createUI()
	net.RecieveMessages(ca.messageChan)
	go ca.consumeMessages()

	if err := ca.app.SetRoot(ca.pages, true).SetFocus(ca.pages).Run(); err != nil {
		panic(err)
	}
}

func Bootstrap() {
	app := NewChatApp()
	app.Start()
}
