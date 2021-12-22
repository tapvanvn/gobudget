package main

import (
	"fmt"
	"time"

	"github.com/tapvanvn/gobudget"
	engines "github.com/tapvanvn/godbengine"
	"github.com/tapvanvn/godbengine/engine"
	"github.com/tapvanvn/godbengine/engine/adapter"
)

func startEngine(eng *engine.Engine) {

	memPool := &adapter.LocalMemDB{}

	err := memPool.Init("")

	if err != nil {

		fmt.Println("cannot init memdb")
	}
	eng.Init(memPool, nil, nil)
}

func claim(beginTime int64, event *gobudget.Event, routeName gobudget.RouteName, require map[gobudget.BudgetName]int64) bool {
	fmt.Println("claim at:", time.Now().Unix()-beginTime)
	if success, err := event.Claim(routeName, require); err != nil {

		if err == gobudget.ErrOutOfBudget {
			fmt.Println("\tout of budget")
		} else {
			panic(err)
		}
		return false
	} else if !success {

		fmt.Println("\tnot success")
		return false
	} else {
		fmt.Println("\tsuccess")
	}
	return true
}

//Total budget is 110 coin
//each 10 second can release 5 package
//each 20 second can release 9 package
func main() {
	engines.InitEngineFunc = startEngine
	eng := engines.GetEngine()

	err := gobudget.InitGoBudget(eng)
	if err != nil {
		panic(err)
	}
	event := gobudget.NewEvent("test")
	event.AddNewBudget("10s", 10, 5)
	event.AddNewBudget("20s", 20, 9)
	event.AddNewBudget("reward", gobudget.NoPeriod, 110)
	event.AddNewRoute("reward", []gobudget.BudgetName{
		"10s", "20s", "reward",
	})

	rewardRequire := map[gobudget.BudgetName]int64{
		"10s":    1,
		"20s":    1,
		"reward": 5,
	}
	claimLoop(event, rewardRequire)
}

func claimLoop(event *gobudget.Event, require map[gobudget.BudgetName]int64) {
	begin := time.Now().Unix()
	packCount := 0
	for {
		if claim(begin, event, "reward", require) {
			packCount++
			fmt.Println("\trelease package:", packCount)
		}
		time.Sleep(time.Second)
	}
}
