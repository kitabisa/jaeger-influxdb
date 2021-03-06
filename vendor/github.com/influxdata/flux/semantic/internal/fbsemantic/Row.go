// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbsemantic

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Row struct {
	_tab flatbuffers.Table
}

func GetRootAsRow(buf []byte, offset flatbuffers.UOffsetT) *Row {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Row{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Row) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Row) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Row) Props(obj *Prop, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Row) PropsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Row) Extends(obj *Var) *Var {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Var)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func RowStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func RowAddProps(builder *flatbuffers.Builder, props flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(props), 0)
}
func RowStartPropsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RowAddExtends(builder *flatbuffers.Builder, extends flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(extends), 0)
}
func RowEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
