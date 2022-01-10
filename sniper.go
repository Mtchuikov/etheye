package etheye

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Sniper struct {
	c *ethclient.Client
}

func (sn *Sniper) GetNewBlockHeader(erh ErrHandler, scn NewBlockHeaderHandler) {
	quit := make(chan bool)
	headersFlow := make(chan *types.Header)
	sub, err := sn.c.SubscribeNewHead(context.Background(), headersFlow)
	if err != nil {
		erh(quit, err)
	}

	for {
		select {
		case err = <-sub.Err():
			erh(quit, err)
		case header := <-headersFlow:
			go scn(quit, sn.c, header)
		case <-quit:
			sub.Unsubscribe()
			return
		}
	}
}

func (sn *Sniper) GetNewBlock(erh ErrHandler, scn NewBlockHandler) {
	sn.GetNewBlockHeader(erh, func(q QuitChannel, c *ethclient.Client, h *types.Header) {
		b, err := sn.c.BlockByNumber(context.Background(), h.Number)
		if err != nil {
			erh(q, err)
		}
		scn(q, c, b)
	})
}

func (sn *Sniper) GetNewEvent(erh ErrHandler, scn NewEventHandler, eq ethereum.FilterQuery) {
	quit := make(chan bool)
	logsFlow := make(chan types.Log)
	sub, err := sn.c.SubscribeFilterLogs(context.Background(), eq, logsFlow)
	if err != nil {
		erh(quit, err)
	}

	for {
		select {
		case err = <-sub.Err():
			erh(quit, err)
		case vlog := <-logsFlow:
			go scn(quit, sn.c, vlog)
		case <-quit:
			sub.Unsubscribe()
			return
		}
	}
}
