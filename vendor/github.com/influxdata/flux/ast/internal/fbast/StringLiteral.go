// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbast

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type StringLiteral struct {
	_tab flatbuffers.Table
}

func GetRootAsStringLiteral(buf []byte, offset flatbuffers.UOffsetT) *StringLiteral {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &StringLiteral{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *StringLiteral) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StringLiteral) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *StringLiteral) BaseNode(obj *BaseNode) *BaseNode {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(BaseNode)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *StringLiteral) Value() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func StringLiteralStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func StringLiteralAddBaseNode(builder *flatbuffers.Builder, baseNode flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(baseNode), 0)
}
func StringLiteralAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(value), 0)
}
func StringLiteralEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
