package filex

import "testing"

func TestReadFile1(t *testing.T) {
	file1 := "D:\\tech\\repo2Company\\other\\word\\config\\paoku\\Desc_EN.txt"
	cont := ReadFile(file1)
	t.Logf("res: %s", cont)
}
