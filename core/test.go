package core

import "os"

func TestSetup() {
	os.Setenv("PORT", "8888")
	os.Setenv("DB", "mongodb://learnt:learnt@localhost:27017/learnt_test")
}
