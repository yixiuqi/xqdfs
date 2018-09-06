package helper

import (
	"bytes"
	"fmt"
)

/*
字符串拼接类
eg.
	str:=NewStringBuilder().Append("1").Append(2).Append("3").ToString()
	str is "123"
*/

type StringBuilder struct {
	msg *bytes.Buffer
}

func NewStringBuilder() *StringBuilder{
	msg:=new(StringBuilder)
	msg.msg=new(bytes.Buffer)
	return msg
}

func NewStringBuilderWithSeq(seq string) *StringBuilder{
	msg:=new(StringBuilder)
	msg.msg=new(bytes.Buffer)
	msg.msg.WriteString("[")
	msg.msg.WriteString(seq)
	msg.msg.WriteString("]")
	return msg
}

func (this *StringBuilder) Append(a ...interface{}) *StringBuilder{
	this.msg.WriteString(fmt.Sprint(a...))
	return this
}

func (this *StringBuilder) ToString() string {
	return this.msg.String()
}
