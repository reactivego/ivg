package mdicons

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/encode"
	"golang.org/x/image/math/f32"
)

var skippedPaths = map[string]string{
	// hardware/svg/production/ic_scanner_48px.svg contains a filled white
	// rectangle that is overwritten by the subsequent path.
	//
	// See https://github.com/google/material-design-icons/issues/490
	//
	// Matches <path fill="#fff" d="M16 34h22v4H16z"/>
	"M16 34h22v4H16z": "#fff",

	// device/svg/production/ic_airplanemode_active_48px.svg and
	// maps/svg/production/ic_flight_48px.svg contain a degenerate path that
	// contains only one moveTo op.
	//
	// See https://github.com/google/material-design-icons/issues/491
	//
	// Matches <path d="M20.36 18"/>
	"M20.36 18": "",
}

// ErrSkip is returned to deliberately skip generating an icon.
//
// When manually debugging one particular icon, it can be useful to add
// something like:
//
//	if baseName != "check_box" { return errSkip }
//
// at the top of func ParseFile.
var ErrSkip = errors.New("skipping SVG to IconVG conversion")

func ParseFile(fqSVGName, dirName, baseName string, size float32, outSize float32, out *bytes.Buffer) (Statistics, error) {
	stat := Statistics{}

	svgData, err := os.ReadFile(fqSVGName)
	if err != nil {
		return stat, err
	}

	varName := upperCase(dirName)
	for _, s := range strings.Split(baseName, "_") {
		varName += upperCase(s)
	}
	fmt.Fprintf(out, "var %s = []byte{", varName)
	defer fmt.Fprintf(out, "\n}\n\n")
	stat.VarNames = []string{varName}

	var enc encode.Encoder
	enc.Reset(
		ivg.ViewBox{
			MinX: -24, MinY: -24,
			MaxX: +24, MaxY: +24},
		ivg.DefaultPalette)

	g := &SVG{}
	if err := xml.Unmarshal(svgData, g); err != nil {
		return stat, err
	}

	var vbx, vby float32
	for i, v := range strings.Split(g.ViewBox, " ") {
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return stat, err
		}
		switch i {
		case 0:
			vbx = float32(f)
		case 1:
			vby = float32(f)
		}
	}
	offset := f32.Vec2{
		vbx * outSize / size,
		vby * outSize / size,
	}

	// adjs maps from opacity to a cReg adj value.
	adjs := map[float32]uint8{}

	for _, p := range g.Paths {
		if fill, ok := skippedPaths[p.D]; ok && fill == p.Fill {
			continue
		}
		if err := ParsePath(&enc, &p, adjs, size, offset, outSize, g.Circles); err != nil {
			return stat, err
		}
		g.Circles = nil
	}

	if len(g.Circles) != 0 {
		if err := ParsePath(&enc, &Path{}, adjs, size, offset, outSize, g.Circles); err != nil {
			return stat, err
		}
		g.Circles = nil
	}

	ivgData, err := enc.Bytes()
	if err != nil {
		return stat, err
	}
	for i, x := range ivgData {
		if i&0x0f == 0x00 {
			out.WriteByte('\n')
		}
		fmt.Fprintf(out, "%#02x, ", x)
	}

	stat.TotalFiles += 1
	stat.TotalSVGBytes += len(svgData)
	stat.TotalIVGBytes += len(ivgData)
	return stat, nil
}

func upperCase(s string) string {
	if a, ok := acronyms[s]; ok {
		return a
	}
	if c := s[0]; 'a' <= c && c <= 'z' {
		return string(c-0x20) + s[1:]
	}
	return s
}

var acronyms = map[string]string{
	"3d":            "3D",
	"ac":            "AC",
	"adb":           "ADB",
	"airplanemode":  "AirplaneMode",
	"atm":           "ATM",
	"av":            "AV",
	"ccw":           "CCW",
	"cw":            "CW",
	"din":           "DIN",
	"dns":           "DNS",
	"dvr":           "DVR",
	"eta":           "ETA",
	"ev":            "EV",
	"gif":           "GIF",
	"gps":           "GPS",
	"hd":            "HD",
	"hdmi":          "HDMI",
	"hdr":           "HDR",
	"http":          "HTTP",
	"https":         "HTTPS",
	"iphone":        "IPhone",
	"iso":           "ISO",
	"jpeg":          "JPEG",
	"markunread":    "MarkUnread",
	"mms":           "MMS",
	"nfc":           "NFC",
	"ondemand":      "OnDemand",
	"pdf":           "PDF",
	"phonelink":     "PhoneLink",
	"png":           "PNG",
	"rss":           "RSS",
	"rv":            "RV",
	"sd":            "SD",
	"sim":           "SIM",
	"sip":           "SIP",
	"sms":           "SMS",
	"streetview":    "StreetView",
	"svideo":        "SVideo",
	"textdirection": "TextDirection",
	"textsms":       "TextSMS",
	"timelapse":     "TimeLapse",
	"toc":           "TOC",
	"tv":            "TV",
	"usb":           "USB",
	"vpn":           "VPN",
	"wb":            "WB",
	"wc":            "WC",
	"whatshot":      "WhatsHot",
	"wifi":          "WiFi",
}
