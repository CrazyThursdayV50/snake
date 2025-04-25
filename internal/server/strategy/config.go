package strategy

import (
	"snake/internal/strategy/clients"
	"snake/internal/strategy/repository"

	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
)

type Config struct {
	Log        *defaultlogger.Config
	Repository *repository.Config
	Clients    *clients.Config
}
