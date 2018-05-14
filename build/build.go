package build

import "bytes"

var (
	// Date 编译时间
	Date string
	// Commit git提交ID
	Commit string
	// git分支
	Branch string
)

// Version 生成版本信息
func Version(prefix string) string {
	var buf bytes.Buffer
	if prefix != "" {
		buf.WriteString(prefix)
	}
	if Date != "" {
		buf.WriteByte('\n')
		buf.WriteString("date: ")
		buf.WriteString(Date)
	}
	if Commit != "" {
		buf.WriteByte('\n')
		buf.WriteString("commit: ")
		buf.WriteString(Commit)
	}
	if Branch != "" {
		buf.WriteByte('\n')
		buf.WriteString("branch: ")
		buf.WriteString(Branch)
	}
	return buf.String()
}
