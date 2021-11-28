package block

type Block struct {
	ID          string
	Breakable   bool    `json:"breakable,omitempty"`
	Durability  float32 `json:"durability,omitempty"`
	Hardness    float32 `json:"hardness,omitempty"`
	Liquid      bool    `json:"liquid,omitempty"`
	Material    string  `json:"material,omitempty"`
	Strength    float32 `json:"strength,omitempty"`
	StepSound   string  `json:"stepSound,omitempty"`
	Transparent bool    `json:"transparent,omitempty"`
	Visible     bool    `json:"visible,omitempty"`
	Obstacle    bool    `json:"obstacle,omitempty"`
	Plant       bool    `json:"plant,omitempty"`
}

func NewBlock(id string) *Block {
	return &Block{ID: id}
}

func (b *Block) breakable() *Block {
	b.Breakable = true
	return b
}

func (b *Block) durability(f float32) *Block {
	b.Durability = f
	return b
}

func (b *Block) hardness(f float32) *Block {
	b.Hardness = f
	return b
}

func (b *Block) material(f string) *Block {
	b.Material = f
	return b
}

func (b *Block) strength(f float32) *Block {
	b.Strength = f
	return b
}

func (b *Block) stepSound(f string) *Block {
	b.StepSound = f
	return b
}

func (b *Block) liquid() *Block {
	b.Liquid = true
	return b
}

func (b *Block) transparent() *Block {
	b.Transparent = true
	return b
}

func (b *Block) visible() *Block {
	b.Visible = true
	return b
}

func (b *Block) obstacle() *Block {
	b.Obstacle = true
	return b
}

func (b *Block) plant() *Block {
	b.Plant = true
	return b
}
