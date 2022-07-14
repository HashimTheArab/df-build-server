package main

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/we/brush"
	"github.com/df-mc/we/palette"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"math/rand"
	"time"
)

// Fill ...
type Fill struct {
	Palette string `cmd:"palette"`
}

// FillAir ...
type FillAir struct{}

// Run ...
//func (f Fill) Run(s cmd.Source, o *cmd.Output) {
//	p := s.(*player.Player)
//	h, ok := p.Handler().(*Handler)
//	if !ok {
//		panic("should never happen")
//	}
//	palette, ok := palette.LookupHandler(p)
//	if !ok {
//		panic("should never happen")
//	}
//	found, ok := palette.Palette(f.Palette)
//	if !ok || len(found.Blocks()) == 0 {
//		o.Error("Invalid palette, create one using /palette")
//		return
//	}
//	var names []string
//	for _, b := range found.Blocks() {
//		n, _ := b.EncodeBlock()
//		names = append(names, n)
//	}
//
//	base := cube.PosFromVec3(h.Pos1.Add(h.Pos2).Vec3().Mul(0.5))
//	length := int(math.Abs(float64(h.Pos2.X() - h.Pos1.X())))
//	height := int(math.Abs(float64(h.Pos2.Y() - h.Pos1.Y())))
//	width := int(math.Abs(float64(h.Pos2.Z() - h.Pos1.Z())))
//	brush.Perform(base, rectangle{length, height, width}, fillAction{b: found.Blocks()}, s.World())
//	o.Print(text.Colourf("<green>Filled area %v to %v with palette %v", h.Pos1, h.Pos2, names))
//}

// Run ...
// this function is temporary, the commented out one is the proper way to do /fill in dragonfly but minecraft math is hard so this is what ur getting for now!
func (f Fill) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	h, ok := p.Handler().(*Handler)
	if !ok {
		panic("should never happen")
	}
	palette, ok := palette.LookupHandler(p)
	if !ok {
		panic("should never happen")
	}
	found, ok := palette.Palette(f.Palette)
	if !ok || len(found.Blocks()) == 0 {
		o.Error("Invalid palette, create one using /palette")
		return
	}
	var names []string
	for _, b := range found.Blocks() {
		n, _ := b.EncodeBlock()
		names = append(names, n)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	minX, maxX := int(math.Min(float64(h.Pos1.X()), float64(h.Pos2.X()))), int(math.Max(float64(h.Pos1.X()), float64(h.Pos2.X())))
	minY, maxY := int(math.Min(float64(h.Pos1.Y()), float64(h.Pos2.Y()))), int(math.Max(float64(h.Pos1.Y()), float64(h.Pos2.Y())))
	minZ, maxZ := int(math.Min(float64(h.Pos1.Z()), float64(h.Pos2.Z()))), int(math.Max(float64(h.Pos1.Z()), float64(h.Pos2.Z())))
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				s.World().SetBlock(cube.Pos{x, y, z}, found.Blocks()[r.Intn(len(found.Blocks()))], &world.SetOpts{
					DisableBlockUpdates:       true,
					DisableLiquidDisplacement: true,
				})
			}
		}
	}
	o.Print(text.Colourf("<green>Filled area %v to %v with palette %v", h.Pos1, h.Pos2, names))
}

func (f FillAir) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	h, ok := p.Handler().(*Handler)
	if !ok {
		panic("should never happen")
	}

	minX, maxX := int(math.Min(float64(h.Pos1.X()), float64(h.Pos2.X()))), int(math.Max(float64(h.Pos1.X()), float64(h.Pos2.X())))
	minY, maxY := int(math.Min(float64(h.Pos1.Y()), float64(h.Pos2.Y()))), int(math.Max(float64(h.Pos1.Y()), float64(h.Pos2.Y())))
	minZ, maxZ := int(math.Min(float64(h.Pos1.Z()), float64(h.Pos2.Z()))), int(math.Max(float64(h.Pos1.Z()), float64(h.Pos2.Z())))
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				s.World().SetBlock(cube.Pos{x, y, z}, block.Air{}, &world.SetOpts{
					DisableBlockUpdates:       true,
					DisableLiquidDisplacement: true,
				})
			}
		}
	}
	o.Print(text.Colourf("<green>Filled area %v to %v with %v", h.Pos1, h.Pos2, "air"))
}

// Allow ...
func (Fill) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

// fillAction is a fill action used internally by df world edit.
type fillAction struct {
	b []world.Block
}

// At returns the world.Block and world.Liquid behind it that should be placed at a specific x, y and z in the
// *world.World passed.
// At should use the *rand.Rand instance passed to produce random numbers and must only use the at function to
// read blocks at a specific position in the world.
// If At returns a nil world.Block, no block will be placed at that position.
func (f fillAction) At(_, _, _ int, r *rand.Rand, _ *world.World, _ func(x, y, z int) world.Block) (world.Block, world.Liquid) {
	return f.b[r.Intn(len(f.b))], nil
}

// Form is required by the WE action interface for extra data, but we don't need it here so return nil.
func (fillAction) Form(brush.Shape) form.Form {
	return nil
}

// rectangle represents the rectangle shape used for /fill -- length(x, y, z).
type rectangle [3]int

// Inside checks if a specific point is within the cube with the centre coordinates passed.
func (r rectangle) Inside(cx, cy, cz, x, y, z int) bool {
	return x <= cx+r[0] && x >= cx-r[0] && y <= cy+r[1] && y >= cy-r[1] && z <= cz+r[2] && z >= cz-r[2]
}

// Dim returns the width, height and length of the rectangle.
func (r rectangle) Dim() [3]int {
	return [3]int{r[0] + 1, r[1] + 1, r[2] + 1}
}
