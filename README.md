# gobudget

### General budget managerment

- event : the event in which budgets had been declare using for.

- budget : the limit resource in (int64) that need to be manage.


### Requirement
- any memory database that supported by godbengine.


### Some scenarior
- The total money can be spend on entire event 

- The total money can be spend on each day of event.

- The total package can be create each day

### Budget Period 
- budget can be classify to 2 main type; the budget that has and doesnt have period. 

- If has period the period id should be ```period_id = floor( timestamp/period ) ```
- The none period budget has period id always be ```period_id = 0```

### Route
- One budget can be independent or depend on others budgets.

- In case the budget is depend on others pudgets. Those dependency budgets will be measure before calc the base budget. And will be also claimed before claiming on base budget.

- We define a budget route type for implementing this. The route object containt an array of budget should be mersure and claim on order from begin to end of array.

### Event Init
- Budget 