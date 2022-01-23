package gobudget

import (
	"fmt"
)

type Recover struct {
	Total      int64
	IssuedTime int64
}

func NewEvent(name EventName) *Event {

	event := &Event{

		Name:           name,
		budgets:        map[BudgetName]*Budget{},
		routes:         map[RouteName]*Route{},
		prefixAndEvent: string(name),
	}

	return event
}

func NewEventWithPrefix(name EventName, prefix string) *Event {

	event := &Event{

		Name:           name,
		budgets:        map[BudgetName]*Budget{},
		routes:         map[RouteName]*Route{},
		prefixAndEvent: fmt.Sprintf("%s_%s", prefix, string(name)),
	}

	return event
}

type Event struct {
	Name           EventName
	budgets        map[BudgetName]*Budget
	routes         map[RouteName]*Route
	prefixAndEvent string
	timeAdjust     int64
}

func (event *Event) SetTimeAdjust(timeAdjust int64) {

	event.timeAdjust = timeAdjust
}

func (event *Event) AddBudget(budget *Budget) {

	event.budgets[budget.Name] = budget
}

func (event *Event) AddNewBudget(budgetName BudgetName, period int64, total int64) {

	event.budgets[budgetName] = &Budget{
		Name:   budgetName,
		Period: period,
		Total:  total,
	}
}

func (event *Event) AddRoute(name RouteName, route *Route) error {

	for _, budName := range route.Trail {

		if _, ok := event.budgets[budName]; !ok {

			return ErrBudgetNotExisted
		}
	}
	event.routes[name] = route
	return nil
}
func (event *Event) AddNewRoute(name RouteName, trail []BudgetName) error {

	for _, budName := range trail {

		if _, ok := event.budgets[budName]; !ok {

			return ErrBudgetNotExisted
		}
	}
	event.routes[name] = &Route{
		Trail: trail,
	}
	return nil
}

func (event *Event) Measure(routeName RouteName, element map[BudgetName]int64) (bool, error) {

	route, hasRoute := event.routes[routeName]

	if !hasRoute {

		return false, ErrRouteNotExisted
	}
	for _, budgetName := range route.Trail {

		budget, hasBudget := event.budgets[budgetName]
		amount, hasRequireBudget := element[budgetName]

		if !hasBudget || !hasRequireBudget || budget == nil {

			return false, ErrBudgetNotExisted
		}

		claimed, err := getBudgetClaimed(event.prefixAndEvent, event.timeAdjust, budget)
		if err != nil {

			return false, err

		}
		recovered, err := getBudgetRecoverd(event.prefixAndEvent, event.timeAdjust, budget)
		if err != nil {

			return false, err

		}
		ajustedClaimed := claimed - recovered
		if ajustedClaimed < 0 {
			ajustedClaimed = 0
		}
		if ajustedClaimed+amount > budget.Total {

			return false, ErrOutOfBudget
		}
	}
	return true, nil
}

func (event *Event) Claim(routeName RouteName, element map[BudgetName]int64) (bool, error) {

	if mersure, err := event.Measure(routeName, element); err != nil || !mersure {

		return false, err
	}
	route, _ := event.routes[routeName]

	for _, budgetName := range route.Trail {

		budget, _ := event.budgets[budgetName]
		amount, _ := element[budgetName]

		_, err := claimBudget(event.prefixAndEvent, event.timeAdjust, budget, amount)

		if err != nil {

			return false, err
		}
	}
	return true, nil
}

func (event *Event) Recover(routeName RouteName, element map[BudgetName]*Recover) error {

	route, hasRoute := event.routes[routeName]

	if !hasRoute {

		return ErrRouteNotExisted
	}

	for _, budgetName := range route.Trail {

		budget, has := event.budgets[budgetName]
		if !has {
			return ErrBudgetNotExisted
		}
		if recover, has := element[budgetName]; has {

			err := recoveredBudget(event.prefixAndEvent, event.timeAdjust, budget, recover.IssuedTime, recover.Total)

			if err != nil {

				return err
			}
		}
	}
	return nil
}

func (event *Event) Reset() error {

	for _, budget := range event.budgets {

		if __eng == nil {

			return ErrInvalidDBEngine
		}

		memPool := __eng.GetMemPool()
		err := memPool.SetInt(budget.GetClaimedKey(event.prefixAndEvent, event.timeAdjust), 0)
		if err != nil {

			return err
		}
		err = memPool.SetInt(budget.GetRecoveredKey(event.prefixAndEvent, event.timeAdjust), 0)
		if err != nil {

			return err
		}
	}
	return nil
}

func (event *Event) GetReport() (map[BudgetName][]int64, error) {

	report := map[BudgetName][]int64{}

	for _, budget := range event.budgets {
		if __eng == nil {

			return nil, ErrInvalidDBEngine
		}

		memPool := __eng.GetMemPool()
		claimed, err := memPool.GetInt(budget.GetClaimedKey(event.prefixAndEvent, event.timeAdjust))
		recovered, err := memPool.GetInt(budget.GetRecoveredKey(event.prefixAndEvent, event.timeAdjust))

		if err != nil {

			return nil, err
		}
		report[budget.Name] = []int64{
			budget.Total,
			claimed,
			recovered,
		}
	}
	return report, nil
}
