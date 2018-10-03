package frugal

import "git.apache.org/thrift.git/lib/go/thrift"

// ReadString reads a string from p and assigns it to obj.
func ReadString(p thrift.TProtocol, obj *string, msg string) error {
	if v, err := p.ReadString(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadBool reads a bool from p and assigns it to obj.
func ReadBool(p thrift.TProtocol, obj *bool, msg string) error {
	if v, err := p.ReadBool(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadByte reads a byte from p and assigns it to obj.
func ReadByte(p thrift.TProtocol, obj *int8, msg string) error {
	if v, err := p.ReadByte(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = int8(v)
	}
	return nil
}

// ReadDouble reads a float64 from p and assigns it to obj.
func ReadDouble(p thrift.TProtocol, obj *float64, msg string) error {
	if v, err := p.ReadDouble(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadI16 reads a int16 from p and assigns it to obj.
func ReadI16(p thrift.TProtocol, obj *int16, msg string) error {
	if v, err := p.ReadI16(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadI32 reads a int32 from p and assigns it to obj.
func ReadI32(p thrift.TProtocol, obj *int32, msg string) error {
	if v, err := p.ReadI32(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadI64 reads a int64 from p and assigns it to obj.
func ReadI64(p thrift.TProtocol, obj *int64, msg string) error {
	if v, err := p.ReadI64(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadBinary reads a []byte from p and assigns it to obj.
func ReadBinary(p thrift.TProtocol, obj *[]byte, msg string) error {
	if v, err := p.ReadBinary(); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	} else {
		*obj = v
	}
	return nil
}

// ReadStruct reads a thrift.TStruct from p and assigns it to obj.
func ReadStruct(p thrift.TProtocol, obj thrift.TStruct, msg string) error {
	if err := obj.Read(p); err != nil {
		return thrift.PrependError("error reading "+msg+":", err)
	}
	return nil
}

// WriteString writes string `value` of field name and id `name` and `field` respectively into `p`.
func WriteString(p thrift.TProtocol, value, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.STRING, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteString(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteBool writes bool `value` of field name and id `name` and `field` respectively into `p`.
func WriteBool(p thrift.TProtocol, value bool, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.BOOL, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteBool(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteByte writes byte `value` of field name and id `name` and `field` respectively into `p`.
func WriteByte(p thrift.TProtocol, value int8, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.BYTE, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteByte(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteDouble writes float64 `value` of field name and id `name` and `field` respectively into `p`.
func WriteDouble(p thrift.TProtocol, value float64, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.DOUBLE, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteDouble(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteI16 writes int16 `value` of field name and id `name` and `field` respectively into `p`.
func WriteI16(p thrift.TProtocol, value int16, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.I16, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteI16(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteI32 writes int32 `value` of field name and id `name` and `field` respectively into `p`.
func WriteI32(p thrift.TProtocol, value int32, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.I32, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteI32(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteI64 writes int64 `value` of field name and id `name` and `field` respectively into `p`.
func WriteI64(p thrift.TProtocol, value int64, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.I64, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteI64(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteBinary writes []byte `value` of field name and id `name` and `field` respectively into `p`.
func WriteBinary(p thrift.TProtocol, value []byte, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.BINARY, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := p.WriteBinary(value); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}

// WriteStruct writes thrift.Struct of filed and id `name` and `field` respectively into `p`.
func WriteStruct(p thrift.TProtocol, value thrift.TStruct, name string, field int16) error {
	if err := p.WriteFieldBegin(name, thrift.STRUCT, field); err != nil {
		return thrift.PrependError("write field begin error: ", err)
	}
	if err := value.Write(p); err != nil {
		return thrift.PrependError("field write error: ", err)
	}
	if err := p.WriteFieldEnd(); err != nil {
		return thrift.PrependError("write field end error: ", err)
	}
	return nil
}
