package main

import (
	"fmt"

	"k8s.io/klog"
)

func main() {
	fmt.Println("Hello Web!")
	klog.Info("Hello Web!")
}
