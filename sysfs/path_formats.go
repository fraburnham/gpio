package sysfs

import "fmt"

func exportFmt(baseDir string) string {
	return fmt.Sprintf("%s/export", baseDir)
}

func unexportFmt(baseDir string) string {
	return fmt.Sprintf("%s/unexport", baseDir)
}

func directionFmt(baseDir string, pin int) string {
	return fmt.Sprintf("%s/gpio%d/direction", baseDir, pin)
}

func valueFmt(baseDir string, pin int) string {
	return fmt.Sprintf("%s/gpio%d/value", baseDir, pin)
}

func edgeFmt(baseDir string, pin int) string {
	return fmt.Sprintf("%s/gpio%d/edge", baseDir, pin)
}
