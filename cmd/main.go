package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/KyberNetwork/dgx-price-feeder"
	rsblockchain "github.com/KyberNetwork/dgx-price-feeder/blockchain"
	"github.com/KyberNetwork/dgx-price-feeder/feed"
	"github.com/KyberNetwork/dgx-price-feeder/runner"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/robfig/cron"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

//set config log
func configLog() {
	logger := &lumberjack.Logger{
		Filename: "/go/src/github.com/KyberNetwork/dgx-price-feeder/log/log.log",
		// MaxSize:  1, // megabytes
		MaxBackups: 0,
		MaxAge:     0, //days
		// Compress:   true, // disabled by default
	}

	mw := io.MultiWriter(os.Stdout, logger)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(mw)

	c := cron.New()
	c.AddFunc("@daily", func() { logger.Rotate() })
	c.Start()
}

func main() {

	configLog()

	endpoints := []string{
		"https://semi-node.kyber.network",
		"https://mainnet.infura.io",
		"https://api.mycryptoapi.com/eth",
		"https://api.myetherapi.com/eth",
		"https://mew.giveth.io/",
	}
	chainType := "byzantium"
	operators := map[string]*blockchain.Operator{}
	bc, err := blockchain.NewMinimalBaseBlockchain(
		endpoints, operators, chainType,
	)
	if err != nil {
		panic(err)
	}
	passphrase, err := ioutil.ReadFile(
		"/go/src/github.com/KyberNetwork/dgx-price-feeder/cmd/passphrase",
	)
	if err != nil {
		panic(err)
	}
	reserve := rsblockchain.NewDGXReserve(
		bc,
		ethereum.HexToAddress("0xce076f8ab3f5af34ecf70b99995b11039190edc1"),
		"/go/src/github.com/KyberNetwork/dgx-price-feeder/cmd/keystore",
		strings.TrimSpace(string(passphrase)),
	)
	feedCorpus := feed.NewFeedCorpus()
	runner := runner.NewTickerRunner(30 * time.Minute)
	feeder := dgxpricing.NewPriceFeeder(runner, reserve, feedCorpus)
	feeder.Run()
}
