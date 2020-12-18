package gamestate

import (
	"strings"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/pkg/fileutils"
)

func TestRegisteredActionTypes(t *testing.T) {
	currDirFiles, err := fileutils.GetAllFilesInCurrDir()
	if err != nil {
		t.Errorf("Can't get files in current directory: %v", err)
		return
	}

	actionNamesFromFiles := map[string]bool{}
	for _, file := range currDirFiles {
		name := file.Name()
		if len(name) > 7 {
			if name[:7] == "action_" && name[len(name)-7:] != "test.go" {
				lowercaseActionName := name[7 : len(name)-3]
				actionNamesFromFiles[lowercaseActionName] = false
			}
		}
	}

	actionNamesFromRegister := map[string]bool{}
	for name := range actionTypeStringMap {
		actionNamesFromRegister[strings.ToLower(name.String())] = false
	}

	for name := range actionNamesFromFiles {
		if _, ok := actionNamesFromRegister[name]; !ok {
			t.Errorf("Did not register action: %v", name)
		} else {
			delete(actionNamesFromRegister, name)
		}
	}

	for name := range actionNamesFromRegister {
		t.Errorf("Each action should be in a file with naming convention 'action_<action-name>.go'. "+
			"Failed action: %v", name)
	}
}

func TestRegisteredActionNumber(t *testing.T) {
	if len(actionTypeStringMap) != int(actionTypeEnd) {
		t.Errorf("Not all ActionTypes registered!")
	}
}
