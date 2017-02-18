package main

import (
	"image/png"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jakobvarmose/go-qidenticon"
)

func handleIdenticon(w http.ResponseWriter, req *http.Request, address string) {
	w.Header().Set("Cache-Control", "max-age=604800")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	code := qidenticon.Code(address)
	image := qidenticon.Render(code, 30, qidenticon.DefaultSettings())
	err := png.Encode(w, image)
	if err != nil {
		logrus.Error(err)
	}
}
