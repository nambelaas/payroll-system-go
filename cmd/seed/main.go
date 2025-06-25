package main

import (
	"github.com/nambelaas/payroll-system-go/internal/seeder"
	"github.com/nambelaas/payroll-system-go/pkg"
)

func main() {
	pkg.ConnectDB()
	seeder.SeedUser()
}
