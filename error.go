package gobudget

import "errors"

var ErrInvalidDBEngine = errors.New("Invalid database engine")
var ErrBudgetNotExisted = errors.New("Budget is not existed")
var ErrRouteNotExisted = errors.New("Route is not existed")
var ErrOutOfBudget = errors.New("Out of budget")
