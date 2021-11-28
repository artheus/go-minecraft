package block

const (
	AirID        = "core:air"
	GrassBlockID = "core:grass_block"
	DirtID       = "core:dirt"
	StoneID      = "core:stone"
	SandID       = "core:sand"
	WoodID       = "core:wood"
	LeavesID     = "core:leaves"
	GrassID      = "core:grass"
	DandelionID  = "core:dandelion"
	CloudID      = "core:cloud"
)

func LoadBlocks() {
	// TODO: Load from assets
	_ = AddBlock(NewBlock(AirID))

	_ = AddBlock(
		NewBlock(GrassBlockID).
			breakable().
			obstacle().
			visible().
			durability(0.5).
			hardness(0.5).
			material("grass").
			strength(0.5).
			stepSound("grass"),
		)

	_ = AddBlock(
		NewBlock(DirtID).
			breakable().
			visible().
			obstacle().
			durability(0.5).
			hardness(0.5).
			material("dirt").
			strength(0.5).
			stepSound("dirt"),
	)

	_ = AddBlock(
		NewBlock(StoneID).
			breakable().
			visible().
			obstacle().
			durability(1.5).
			hardness(1.5).
			material("stone").
			strength(1.5).
			stepSound("stone"),
	)

	_ = AddBlock(
		NewBlock(SandID).
			breakable().
			visible().
			obstacle().
			durability(0.5).
			hardness(0.5).
			material("sand").
			strength(0.5).
			stepSound("sand"),
	)

	_ = AddBlock(
		NewBlock(LeavesID).
			breakable().
			visible().
			obstacle().
			transparent().
			durability(0.5).
			hardness(0.5).
			material("leaves").
			strength(0.5).
			stepSound("leaves"),
	)

	_ = AddBlock(
		NewBlock(WoodID).
			breakable().
			visible().
			obstacle().
			durability(0.5).
			hardness(0.5).
			material("wood").
			strength(0.5).
			stepSound("wood"),
	)

	_ = AddBlock(
		NewBlock(GrassID).
			breakable().
			visible().
			plant().
			transparent().
			durability(0.1).
			hardness(0.1).
			material("grass").
			strength(0.1),
	)

	_ = AddBlock(
		NewBlock(DandelionID).
			breakable().
			visible().
			plant().
			transparent().
			durability(0.1).
			hardness(0.1).
			material("dandelion").
			strength(0.1),
	)

	_ = AddBlock(
		NewBlock(CloudID).
			visible().
			material("cloud"),
	)
}
