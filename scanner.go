package etheye

import (
	"context"
	"log"
  
        "github.com/ethereum/go-ethereum"
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
		case err := <-sub.Err():
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
	logsFlow := make(chan types.Log)
	sub, err := sn.c.SubscribeFilterLogs(context.Background(), eq, logsFlow)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan bool)
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case <-quit:
			sub.Unsubscribe()
			return
		case vlog := <-logsFlow:
			go scn(quit, sn.c, vlog)
		}
	}
}

func NewSniper(rawurl string) (*Sniper, error) {
	c, err := ethclient.Dial(rawurl)
	if err != nil {
		log.Fatal(err)
	}
	return &Sniper{c}, nil
}
