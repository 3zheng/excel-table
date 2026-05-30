package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashOutput() {
	passwords := []string{"changeme"}
	for _, p := range passwords {
		hash, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		fmt.Printf("密码: %s\n哈希: %s\n\n", p, string(hash))
	}
}
