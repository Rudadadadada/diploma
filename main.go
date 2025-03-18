package main

import (
	admin 	 "diploma/services/admin/pkg/app"
	customer "diploma/services/customer/pkg/app"
)

func main() {
    go admin.Launch()
	go customer.Launch()

	select {}
}