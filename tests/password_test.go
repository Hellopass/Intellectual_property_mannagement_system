package tests

import (
	"intellectual_property/pkg/utils"
	"testing"
)

func Test_GetSlat(t *testing.T) {
	//测试数据
	s := []string{"VDLRTJADPNNMUOCU",
		"VXHRKKZVQBNQTCQN",
		"GUBSEBISTIFTGRRL",
		"GXVUAZWDOFLGWXLR",
	}
	for i := 0; i < len(s); i++ {
		t.Log(utils.GetSlat(s[i]))
	}
}
