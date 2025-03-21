package main

import (
	admin "diploma/services/admin/pkg/app"
	authorization "diploma/services/authorization/pkg/app"
	customer "diploma/services/customer/pkg/app"
	courier "diploma/services/courier/pkg/app"
)

func main() {
	go admin.Launch()
	go customer.Launch()
	go authorization.Launch()
	go courier.Launch()

	select {}
}
