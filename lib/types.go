package lib

import (
	goirc "github.com/thoj/go-ircevent"
)

// TODO: give modules their own sections in config?
type Config struct {
	User struct {
		Nick string
		User string
	}
	Server struct {
		Host string
		Port string
	}
	Channel struct {
		Channelname string
		Greeting    string
	}
	Connection struct {
		Verbose_callback_handler bool
		Debug                    bool
	}
	Module struct {
		Insult_swearfile string
	}
}

type Module struct {
	Configure func(*Config)
	GetHelp   func() []string
	Run       func(*Privmsg) bool
}

type Privmsg struct {
	Msg     string
	MsgArgs []string
	Dest    string
	Event   *goirc.Event
	Conn    *goirc.Connection
	SayChan chan SayMsg
}

type SayMsg struct {
	Conn *goirc.Connection
	Dest string
	Msg  string
}
