/*
  btcbot is a Bitcoin trading bot for HUOBI.com written
  in golang, it features multiple trading methods using
  technical analysis.

  Disclaimer:

  USE AT YOUR OWN RISK!

  The author of this project is NOT responsible for any damage or loss caused
  by this software. There can be bugs and the bot may not perform as expected
  or specified. Please consider testing it first with paper trading /
  backtesting on historical data. Also look at the code to see what how
  it's working.

  Weibo:http://weibo.com/bocaicfa
*/

package main

import (
	. "config"
	"fmt"
	"huobi"
	"logger"
	"okcoin"
	"strconv"
	"time"
)

func backtesting() {
	fmt.Println("back testing begin...")
	huobi := huobi.NewHuobi()
	huobi.Disable_trading = 1

	peroids := []int{1, 5, 15, 30, 60, 100}
	for _, v := range peroids {
		huobi.Peroid = v
		if huobi.TradeKLinePeroid(huobi.Peroid) == true {

		} else {
			logger.Errorln("TradeKLine failed.")
		}
	}

	fmt.Println("生成 1/5/15/30/60分钟及1天 周期的后向测试报告于log/reportxxx.log文件中,请查看")

	fmt.Println("back testing end ...")
}

func testKLineAPI(done chan bool) {
	ticker := time.NewTicker(2000 * time.Millisecond) //2s

	huobi := huobi.NewHuobi()
	huobi.Peroid, _ = strconv.Atoi(Option["tick_interval"])
	totalHour, _ := strconv.ParseInt(Option["totalHour"], 0, 64)
	if totalHour < 1 {
		totalHour = 1
	}
	slippage, err := strconv.ParseFloat(Config["slippage"], 64)
	if err != nil {
		logger.Debugln("config item slippage is not float")
		slippage = 0
	}
	huobi.Slippage = slippage

	huobi.Disable_trading = 0

	go func() {
		for _ = range ticker.C {
			if huobi.Peroid == 1 {
				huobi.TradeKLineMinute()
			} else {
				huobi.TradeKLinePeroid(huobi.Peroid)
			}
		}
	}()

	oneHour := 60 * 60 * 1000 * time.Millisecond

	logger.Infof("程序将持续运行%d小时后停止", time.Duration(totalHour))

	time.Sleep(time.Duration(totalHour) * oneHour)

	ticker.Stop()
	fmt.Println("程序到达设定时长%d小时，停止运行。", time.Duration(totalHour))
	done <- true
}

func testHuobiAPI() {
	tradeAPI := huobi.NewHuobiTrade(SecretOption["access_key"], SecretOption["secret_key"])
	accout_info, _ := tradeAPI.Get_account_info()
	fmt.Println(accout_info)

	//	fmt.Println(tradeAPI.Get_account_info())
	if false {
		buyId := tradeAPI.Buy("1000", "0.001")
		sellId := tradeAPI.Sell("10000", "0.001")

		//fmt.Println(tradeAPI.Get_delegations())
		if tradeAPI.Cancel_delegation(buyId) {
			fmt.Printf("cancel %s success \n", buyId)
		} else {
			fmt.Printf("cancel %s falied \n", buyId)
		}

		if tradeAPI.Cancel_delegation(sellId) {
			fmt.Printf("cancel %s success \n", sellId)
		} else {
			fmt.Printf("cancel %s falied \n", sellId)
		}
	}

	fmt.Println(tradeAPI.Get_delegations())
}

func testOkcoinAPI() {
	tradeAPI := okcoin.NewOkcoinTrade(SecretOption["ok_partner"], SecretOption["ok_secret_key"])
	accout_info, _ := tradeAPI.Get_account_info()
	fmt.Println(accout_info)

	buyret := tradeAPI.BuyBTC("1000", "0.01")
	fmt.Println(buyret)
	sellret := tradeAPI.SellBTC("10000", "0.01")
	fmt.Println(sellret)

	var orderTable okcoin.OrderTable
	orderTable, ret := tradeAPI.Get_BTCorder("-1")
	fmt.Println(ret, orderTable)

	time.Sleep(2000 * time.Millisecond)

	orderTable, ret = tradeAPI.Get_LTCorder("-1")
	fmt.Println(ret, orderTable)

	ret = tradeAPI.Cancel_BTCorder("-1")
	fmt.Println(ret)

	time.Sleep(2000 * time.Millisecond)

	ret = tradeAPI.Cancel_LTCorder("-1")
	fmt.Println(ret)
}
func tradeService() {

	done := make(chan bool, 1)

	fmt.Println("robot working...")

	//backtesting()

	go testKLineAPI(done)
	<-done

	fmt.Println("done")

	return
	//doTradeDelegation()
}