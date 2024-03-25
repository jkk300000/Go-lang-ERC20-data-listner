package main

import "log"

import (

    "fmt"
    "context"
    "encoding/hex"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/core/types"
)

func BlockListner() error{

  client,err := ethclient.Dial("wss://mainnet.infura.io/ws/v3/85b90e9f4c784126b7fa89e4f90b2259");
  if err != nil{
    log.Fatal(err)
  }
  headers := make(chan *types.Header)

  sub, err := client.SubscribeNewHead(context.Background(),headers)
  if err != nil{
    log.Fatal(err)
  }

  for {
    select {
    case err := <- sub.Err():
      log.Fatal(err)
    case header := <-headers:
      block, err := client.BlockByHash(context.Background(), header.Hash())
      if err != nil{
        log.Fatal(err)
      }

      fmt.Println("############################################################################")
      fmt.Println("Block Hash : " + block.Hash().Hex())
      fmt.Println("Block Number : ", block.Number().Uint64())
      fmt.Println("Block Timestamp : ", block.Time())
      fmt.Println("Block Nonce : ", block.Nonce())
      fmt.Println("Total Transaction : ", len(block.Transactions()))
      fmt.Println("############################################################################")
      baseFee := block.BaseFee()

      if (len(block.Uncles())>0){
        for _, uncle := range block.Uncles(){
          uncleFee := float64((uncle.Number.Uint64() +8 - block.Number().Uint64())*2) /8.0
          fmt.Println("Uncle Block Length :", len(block.Uncles()))
          fmt.Println("Uncle Miner Address :", uncle.Coinbase.Hex())
          fmt.Println("Uncle Block Number :", uncle.Number.Uint64())
          fmt.Println("Uncle Block Reward :", uncleFee)
        }
      }


      for _, tx := range block.Transactions(){
        fmt.Println("***********************************************************************************")
        fmt.Println("Transaction Hash :", tx.Hash().Hex())
        if tx.To() != nil{
          fmt.Println("To Address :", tx.To())
        }else{
          fmt.Println("To Address : Contract Creation")
          //TODO
            contractAddress :=GetContractAddress(client, tx.Hash())
            fmt.Println("To Address :", contractAddress)

        }
        fmt.Println("Transfer value(wei)" + tx.Value().String())
        fmt.Println("Transaction nonce :", tx.Nonce())
        fmt.Println("Transaction Gas Limit :", tx.Gas()) //예상 gas limit
        //TODO
        realGasLimit := GetRealGasUsed(client,tx.Hash())
        fmt.Println("Transaction Real Gas Limit :", realGasLimit)
        fmt.Println("Transaction GasFeeCap :", tx.GasFeeCap().Uint64())
        fmt.Println("Transaction GasTipCap :", tx.GasTipCap().Uint64())

        // TODO
        realGasPrice := GetRealGasPrice(baseFee.Uint64(),tx.GasFeeCap().Uint64(), tx.GasTipCap().Uint64())
        fmt.Println("Transaction gasPrice :", realGasPrice)
        fmt.Println("Transaction Input Data :", hex.EncodeToString(tx.Data()))

        if (len(tx.Data()) !=0){
            to,value :=ERC20Transaction(hex.EncodeToString(tx.Data()))
            if to !=""{


            symbol, name, decimal := GetContractInfo(client,tx.To())
            fmt.Println("ERC20 Contract Address : ",tx.To().Hex())
            fmt.Println("ERC20 Contract Name : ",name)
            fmt.Println("ERC20 Contract Symbol : ",symbol)
            fmt.Println("ERC20 Contract Decimal : ",decimal)
            fmt.Println("ERC20 Transfer To Address : ",to)
            fmt.Println("ERC20 Transfer value : ",value) // value/ (10**decimal)
          }
        }
      }
      //TODO ERC20

    }
  }
}
