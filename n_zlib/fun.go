package n_zlib

import (
	"bytes"
	"compress/zlib"
	"io"
	"bbTool/n_log"
	"encoding/hex"
)

func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) ([]byte,error) {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, e := zlib.NewReader(b)
	if e != nil {
		return nil, n_log.ErroBack("zlib err %v",e)
	}

	io.Copy(&out, r)
	return out.Bytes(),nil
}



func DoZlibCompress_hex(src []byte) string {
	i := DoZlibCompress(src)
	str := hex.EncodeToString(i)
	return str
}

//进行zlib解压缩
func DoZlibUnCompress_hex(compressSrc []byte) ([]byte,error) {
	i,e := hex.DecodeString(string(compressSrc))

	if e != nil {
		return nil,e
	}

	return DoZlibUnCompress(i)
}