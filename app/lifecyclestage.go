package app

//go:generate go run golang.org/x/tools/cmd/stringer -type=LifecycleStage -linecomment

// Represents a stage in the Developer Journey lifecycle.
type LifecycleStage int

const (
	LifecycleStageCode    LifecycleStage = iota + 1 // code
	LifecycleStageBuild                             // build
	LifecycleStageTest                              // test
	LifecycleStageRelease                           // release
	LifecycleStageDeploy                            // deploy
	LifecycleStageOperate                           // operate
	LifecycleStageMonitor                           // monitor
	LifecycleStageOther                             // other
)
