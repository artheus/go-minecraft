package chunk

type Context struct {
	chunk      Chunk
	actionChan chan ChunkAction
}
