package rpc

import (
	"github.com/artheus/go-minecraft/core/block"
	"github.com/artheus/go-minecraft/core/ctx"
	"github.com/artheus/go-minecraft/core/game/store"
	"github.com/artheus/go-minecraft/core/types"
	. "github.com/artheus/go-minecraft/math/f32"

	"flag"
	"log"
	"net"
	"net/rpc"
	"strings"

	gocraft "github.com/icexin/gocraft-server/client"
	"github.com/icexin/gocraft-server/proto"
)

var (
	serverAddr = flag.String("s", "", "server address")

	Client *gocraft.Client
)

func InitClient(ctx *ctx.Context) error {
	if *serverAddr == "" {
		return nil
	}
	addr := *serverAddr
	if strings.Index(addr, ":") == -1 {
		addr += ":8421"
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	Client = gocraft.NewClient()
	Client.RegisterService("Block", &BlockService{ctx: ctx})
	Client.RegisterService("Player", &PlayerService{ctx: ctx})
	Client.Start(conn)
	return nil
}

func ClientFetchChunk(id Vec3, f func(bid Vec3, w *block.Block)) {
	if Client == nil {
		return
	}
	req := proto.FetchChunkRequest{
		P:       int(id.X),
		Q:       int(id.Z),
		Version: store.Storage.GetChunkVersion(id),
	}
	rep := new(proto.FetchChunkResponse)
	err := Client.Call("Block.FetchChunk", req, rep)
	if err == rpc.ErrShutdown {
		return
	}
	if err != nil {
		log.Panic(err)
	}
	/*for _, b := range rep.Blocks {
		f(Vec3{X: float32(b[0]), Y: float32(b[1]), Z: float32(b[2])}, b[3])
	}*/
	if req.Version != rep.Version {
		store.Storage.UpdateChunkVersion(id, rep.Version)
	}
}

func ClientUpdateBlock(id Vec3, w *block.Block) {
	if Client == nil {
		return
	}
	cid := id.ChunkID()
	req := &proto.UpdateBlockRequest{
		Id: Client.ClientId,
		P:  int(cid.X),
		Q:  int(cid.Z),
		X:  int(id.X),
		Y:  int(id.Y),
		Z:  int(id.Z),
		//W:  w,
	}
	rep := new(proto.UpdateBlockResponse)
	err := Client.Call("Block.UpdateBlock", req, rep)
	if err == rpc.ErrShutdown {
		return
	}
	if err != nil {
		log.Panic(err)
	}
	store.Storage.UpdateChunkVersion(cid, rep.Version)
}

func ClientUpdatePlayerState(ctx *ctx.Context, state types.PlayerState) {
	if Client == nil {
		return
	}
	req := &proto.UpdateStateRequest{
		Id: Client.ClientId,
	}
	s := &req.State
	s.X, s.Y, s.Z, s.Rx, s.Ry = state.X, state.Y, state.Z, state.Rx, state.Ry
	rep := new(proto.UpdateStateResponse)
	err := Client.Call("Player.UpdateState", req, rep)
	if err == rpc.ErrShutdown {
		return
	}
	if err != nil {
		log.Panic(err)
	}

	for id, player := range rep.Players {
		ctx.Game().PlayerRenderer().UpdateOrAdd(id, player)
	}
}

type BlockService struct {
	ctx *ctx.Context
}

func (s *BlockService) UpdateBlock(req *proto.UpdateBlockRequest, rep *proto.UpdateBlockResponse) error {
	log.Printf("rpc::UpdateBlock:%v", *req)
	bid := Vec3{float32(req.X), float32(req.Y), float32(req.Z)}
	//s.ctx.Game().World().UpdateBlock(bid, req.W)
	s.ctx.Game().ChunkRenderer().DirtyChunk(bid.ChunkID())
	return nil
}

type PlayerService struct {
	ctx *ctx.Context
}

func (s *PlayerService) RemovePlayer(req *proto.RemovePlayerRequest, rep *proto.RemovePlayerResponse) error {
	s.ctx.Game().PlayerRenderer().Remove(req.Id)
	return nil
}
