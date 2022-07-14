package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/we"
	"github.com/df-mc/we/brush"
	"github.com/df-mc/we/palette"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// register commands
	cmd.Register(cmd.New("brush", "", nil, brush.BindCommand{}, brush.UnbindCommand{}, brush.UndoCommand{}))
	cmd.Register(cmd.New("palette", "", nil, palette.DeleteCommand{}, palette.SaveCommand{}, palette.SetCommand{}))
	cmd.Register(cmd.New("fill", "fill an area", nil, Fill{}))
	cmd.Register(cmd.New("fillair", "/fill but it fills with air bc u cant have an air palette", nil, FillAir{}))
	cmd.Register(cmd.New("gamemode", "", nil, GameMode{}))

	srv := server.New(&config, log)
	srv.World().SetTime(0)
	srv.World().StopTime()
	srv.CloseOnProgramEnd()
	if err := srv.Start(); err != nil {
		log.Fatalln(err)
	}

	// this is the loop to accept players, basically a join event
	for srv.Accept(acceptPlayer) {
	}
}

// handleJoin is a join event, this code is ran for each player that joins the server.
func acceptPlayer(p *player.Player) {
	p.Handle(&Handler{
		p: p,
		w: we.NewHandler(p),
	})
	p.ShowCoordinates()
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
