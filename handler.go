package main

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/we"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Handler is the handler for the player.
type Handler struct {
	player.NopHandler

	p *player.Player
	w *we.Handler

	// These are the users selected positions, we can do this since handlers are per player.
	Pos1, Pos2 cube.Pos
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(ctx *event.Context) {
	h.w.HandleItemUse(ctx)
}

// HandleItemUseOnBlock ...
func (h *Handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, vec mgl64.Vec3) {
	h.w.HandleItemUseOnBlock(ctx, pos, face, vec)
	if held, _ := h.p.HeldItems(); held.Item() == (item.Stick{}) && pos != h.Pos2 {
		h.SetPos2(pos)
	}
}

// HandleBlockBreak ...
func (h *Handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	h.w.HandleBlockBreak(ctx, pos, drops)
	if held, _ := h.p.HeldItems(); held.Item() == (item.Stick{}) {
		ctx.Cancel()
		h.SetPos1(pos)
	}
}

// HandleQuit ...
func (h *Handler) HandleQuit() {
	h.w.HandleQuit()
}

/* User Functions ------------------------------ BELOW  --------------------------------------- User Functions */

// SetPos1 sets the first position saved for the player.
func (h *Handler) SetPos1(p cube.Pos) {
	h.Pos1 = p
	h.p.Message(text.Colourf("<green>Pos1 has been set to %v</green>", p))
}

// SetPos2 sets the second position saved for the player.
func (h *Handler) SetPos2(p cube.Pos) {
	h.Pos2 = p
	h.p.Message(text.Colourf("<green>Pos2 has been set to %v</green>", p))
}
