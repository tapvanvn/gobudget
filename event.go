package gobudget

import (
	"fmt"
)

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

		claimed, err := getBudgetClaimed(event.prefixAndEvent, budget)
		if err != nil {

			return false, err

		} else if claimed+amount > budget.Total {

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

		_, err := claimBudget(event.prefixAndEvent, budget, amount)

		if err != nil {

			return false, err
		}
	}
	return true, nil
}

func (event *Event) Reset() error {

	for _, budget := range event.budgets {
		if __eng == nil {

			return ErrInvalidDBEngine
		}

		memPool := __eng.GetMemPool()
		err := memPool.SetInt(budget.GetClaimedKey(event.prefixAndEvent), 0)
		if err != nil {

			return err
		}
	}
	return nil
}

func (event *Event) GetReport() (map[BudgetName]int64, error) {

	report := map[BudgetName]int64{}

	for _, budget := range event.budgets {
		if __eng == nil {

			return nil, ErrInvalidDBEngine
		}

		memPool := __eng.GetMemPool()
		claimed, err := memPool.GetInt(budget.GetClaimedKey(event.prefixAndEvent))
		if err != nil {

			return nil, err
		}
		report[budget.Name] = claimed
	}
	return report, nil
}
