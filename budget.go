package gobudget

import (
	"fmt"
	"math"
	"time"
)

type Budget struct {
	Name   BudgetName
	Period int64 //-1 mean have no preiod
	Total  int64 //total budget
}

func (budget *Budget) CurrentPeriodID(timeAdjust int64) int64 {

	return budget.PeriodID(timeAdjust, time.Now().Unix())
}

func (budget *Budget) PeriodID(timeAdjust int64, timestamp int64) int64 {

	if budget.Period == -1 {

		return 0
	}
	moment := timestamp + timeAdjust

	return moment / budget.Period
}

func (budget *Budget) GetRecoveredKey(prefixAndEvent string, timeAdjust int64) string {

	return fmt.Sprintf("%s.%s.%d_rcv", prefixAndEvent, budget.Name, budget.CurrentPeriodID(timeAdjust))
}
func (budget *Budget) GetRecoveredKeyAtMoment(prefixAndEvent string, timeAdjust int64, timestamp int64) string {

	return fmt.Sprintf("%s.%s.%d_rcv", prefixAndEvent, budget.Name, budget.PeriodID(timeAdjust, timestamp))
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

func getBudgetRecoverd(prefixAndEvent string, timeAdjust int64, budget *Budget) (int64, error) {

	if __eng == nil {

		return math.MaxInt64, ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()

	return memPool.GetInt(budget.GetRecoveredKey(prefixAndEvent, timeAdjust))
}

func claimBudget(prefixAndEvent string, timeAdjust int64, budget *Budget, amount int64) (int64, error) {

	if __eng == nil {

		return math.MaxInt64, ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()

	return memPool.IncrIntBy(budget.GetClaimedKey(prefixAndEvent, timeAdjust), amount)
}

func recoveredBudget(prefixAndEvent string, timeAdjust int64, budget *Budget, issuedTime int64, totalRecovered int64) error {

	if __eng == nil {

		return ErrInvalidDBEngine
	}

	memPool := __eng.GetMemPool()

	return memPool.SetInt(budget.GetRecoveredKeyAtMoment(prefixAndEvent, timeAdjust, issuedTime), totalRecovered)
}
