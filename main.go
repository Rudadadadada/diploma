package main

import (
	admin "diploma/services/admin/pkg/app"
	authorization "diploma/services/authorization/pkg/app"
	courier "diploma/services/courier/pkg/app"
	customer "diploma/services/customer/pkg/app"
	distribution "diploma/services/distribution/pkg/app"
	order "diploma/services/order/pkg/app"
)

func main() {
	go admin.Launch()
	go customer.Launch()
	go authorization.Launch()
	go courier.Launch()
	go order.Launch()
	go distribution.Launch()

	select {}
}
