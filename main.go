package main

import (
	"github.com/jinzhu/gorm"
	"github.com/maiguangyang/gqlgen/cmd"
)

var db *gorm.DB

func main() {
	cmd.Execute()
}
