package dashboard

import (
	"sync"
	"testing"
)

func TestBaseDashboard_LaunchAndTerminate(t *testing.T) {
	group := &sync.WaitGroup{}
	dashboard := BaseDashboard{}
	dashboard.Launch(group)
	if (dashboard.terminated) {
		t.Error("Dashboard has just launched. It should not have been terminated.")
	}
	dashboard.Terminate()
	group.Wait()
	if (!dashboard.terminated) {
		t.Error("Dashboard should have been terminated after the call to Terminate()")
	}
}
