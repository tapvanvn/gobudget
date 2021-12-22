package gobudget

import (
	"fmt"
	"math"
	"time"
)

type Budget struct {
	Name   BudgetName
	Period int64 //-1 mean have no preiod
	Total  int64
}

func (budget *Budget) CurrentPeriodID() int64 {

	if budget.Period == -1 {

		return 0
	}
	return time.Now().Unix() / budget.Period
}

func (budget *Budget) GetClaimedKey(prefixAndEvent string) string {

	return fmt.Sprintf("%s.%s.%d", prefixAndEvent, budget.Name, budget.CurrentPeriodID())
}

func getBudgetClaimed(prefixAndEvent string, budget *Budget) (int64, error) {

	if __eng == nil {

		return math.MaxInt64, ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()

	return memPool.GetInt(budget.GetClaimedKey(prefixAndEvent))
}

func claimBudget(prefixAndEvent string, budget *Budget, amount int64) (int64, error) {

	if __eng == nil {

		return math.MaxInt64, ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()
	return memPool.IncrIntBy(budget.GetClaimedKey(prefixAndEvent), amount)
}
