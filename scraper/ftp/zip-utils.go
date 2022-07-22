package ftp

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

func gzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gw, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
	defer gw.Close()
	gw.Write(data)
	return err
}

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
