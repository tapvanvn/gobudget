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

func recover(event *gobudget.Event, routeName gobudget.RouteName, recoveredPacks map[gobudget.BudgetName]*gobudget.Recover) {

	fmt.Println("recover")

	if err := event.Recover(routeName, recoveredPacks); err != nil {

		fmt.Println("\tfail:", err.Error())

	} else {

		fmt.Println("\tsuccess")
	}
}

//Total budget is 50 coin
//each second take 5 coin=> need 10 second
//each 5 second recoverd 5 coin => after 5 second when out of budget can issue 1 package

func main() {
	engines.InitEngineFunc = startEngine
	eng := engines.GetEngine()

	err := gobudget.InitGoBudget(eng)
	if err != nil {
		panic(err)
	}
	event := gobudget.NewEvent("test")
	//event.AddNewBudget("10s", 10, 5)
	//event.AddNewBudget("20s", 20, 9)
	event.AddNewBudget("reward", gobudget.NoPeriod, 50)
	event.AddNewRoute("reward", []gobudget.BudgetName{
		"reward",
	})

	rewardRequire := map[gobudget.BudgetName]int64{
		"reward": 5,
	}
	claimLoop(event, rewardRequire)
}

func claimLoop(event *gobudget.Event, require map[gobudget.BudgetName]int64) {
	begin := time.Now().Unix()
	packCount := 0
	recoverSecondCount := 0
	recoverPacks := map[gobudget.BudgetName]*gobudget.Recover{}

	for budgetName, _ := range require {

		recoverPacks[budgetName] = &gobudget.Recover{
			IssuedTime: time.Now().Unix(),
			Total:      5,
		}
	}
	for {
		if claim(begin, event, "reward", require) {
			packCount++
			fmt.Println("\trelease package:", packCount)
		}
		if recoverSecondCount == 5 {

			recover(event, "reward", recoverPacks)

			for budgetName, _ := range require {

				recoverPacks[budgetName].IssuedTime = time.Now().Unix()
				recoverPacks[budgetName].Total += 5
			}
			recoverSecondCount = 0
		}
		recoverSecondCount++
		if report, err := event.GetReport(); err == nil {
			fmt.Printf("report: %v\n", report)
		}

		time.Sleep(time.Second)
	}
}
