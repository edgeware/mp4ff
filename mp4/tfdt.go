package mp4

import (
	"io"
	"io/ioutil"
)

// TfdtBox - Track Fragment Decode Time (tfdt)
//
// Contained in : Track Fragment box (traf)
type TfdtBox struct {
	Version             byte
	Flags               uint32
	BaseMediaDecodeTime uint64
}

// DecodeTfdt - box-specific decode
func DecodeTfdt(r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	s := NewSliceReader(data)
	versionAndFlags := s.ReadUint32()
	version := byte(versionAndFlags >> 24)
	flags := versionAndFlags & 0xffffff
	var baseMediaDecodeTime uint64
	if version == 0 {
		baseMediaDecodeTime = uint64(s.ReadUint32())
	} else {
		baseMediaDecodeTime = s.ReadUint64()
	}

	b := &TfdtBox{
		Version:             version,
		Flags:               flags,
		BaseMediaDecodeTime: baseMediaDecodeTime,
	}
	return b, nil
}

// Type - returns box type
func (t *TfdtBox) Type() string {
	return "tfdt"
}

// Size - returns calculated size
func (t *TfdtBox) Size() int {
	return BoxHeaderSize + 8 + 4*int(t.Version)
}

// Encode - write box to w
func (t *TfdtBox) Encode(w io.Writer) error {
	err := EncodeHeader(t, w)
	if err != nil {
		return err
	}
	buf := makebuf(t)
	bw := NewBufferWrapper(buf)
	versionAndFlags := (uint32(t.Version) << 24) + t.Flags
	bw.WriteUint32(versionAndFlags)
	if t.Version == 0 {
		bw.WriteUint32(uint32(t.BaseMediaDecodeTime))
	} else {
		bw.WriteUint64(t.BaseMediaDecodeTime)
	}
	_, err = w.Write(buf)
	return err
}
