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

func (budget *Budget) CurrentPeriodID(timeAdjust int64) int64 {

	if budget.Period == -1 {

		return 0
	}
	now := time.Now().Unix() + timeAdjust
	return now / budget.Period
}

func (budget *Budget) GetClaimedKey(prefixAndEvent string, timeAdjust int64) string {

	return fmt.Sprintf("%s.%s.%d", prefixAndEvent, budget.Name, budget.CurrentPeriodID(timeAdjust))
}

func getBudgetClaimed(prefixAndEvent string, timeAdjust int64, budget *Budget) (int64, error) {

	if __eng == nil {

		return math.MaxInt64, ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()

	return memPool.GetInt(budget.GetClaimedKey(prefixAndEvent, timeAdjust))
}

func claimBudget(prefixAndEvent string, timeAdjust int64, budget *Budget, amount int64) (int64, error) {

	if __eng == nil {

		return math.MaxInt64, ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()
	return memPool.IncrIntBy(budget.GetClaimedKey(prefixAndEvent, timeAdjust), amount)
}
