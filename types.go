package etheye

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type QuitChannel chan bool

type ErrHandler func(q QuitChannel, e error)

type NewBlockHeaderHandler func(q QuitChannel, c *ethclient.Client, h *types.Header)

type NewBlockHandler func(q QuitChannel, c *ethclient.Client, b *types.Block)

type NewEventHandler func(q QuitChannel, c *ethclient.Client, l types.Log)
