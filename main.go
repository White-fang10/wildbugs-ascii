package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	cols  = 47
	rows  = 32
	fps   = 30
	start = 0.5 // seconds before the first scan
	pass  = 2.0 // seconds for the scan line to travel down the logo
	head  = 2.0 // rows scrambled by the scan line
	trail = 5.0 // rows of afterglow behind the scan line
	noise = "#$%&*+<=>?@[]{}"
)

var palette = []string{
	"#b7ccb1", "#141818", "#459529", "#384138", "#96a791", "#fcfdfc",
	"#5ba831", "#686d6a", "#d0decc", "#666969", "#86ba74", "#7e847e",
	"#bbc6b8", "#fdfefd", "#94c679", "#9da09e", "#e9f0e8", "#b9baba",
	"#c7dfbf", "#c3c6c3", "#e0e5de", "#fefefe", "#cee5c1", "#d2d3d2",
	"#b7ccb1", "#e8ebef", "#459529", "#565656", "#96a791", "#fcfdfc",
	"#5ba831", "#686d6a", "#d0decc", "#f0f2f5", "#86ba74", "#919191",
	"#bbc6b8", "#fdfefd", "#94c679", "#9da09e", "#e9f0e8", "#f8f9fa",
	"#c7dfbf", "#cccccc", "#e0e5de", "#fefefe", "#cee5c1", "#d2d3d2",
}

const monoArt = `
             .88qq_         _pq88:
               '"88q,     .p8P"'
                  )8bppqqqp8P
        qq,     _p88888888888q,     \pp
        d8b    _88P""88888""888_    d8P
        '88_, _8888q8888888p8888p  _88'
         "888q8888888888888888888p88P"
            _"""\888888888888888""_
         _p888 _8888888888888888|888q_
     __p888P" \88888888888888888b "Y888q_,
  _pp888P"    d888888888888888888    "Y888qq,
  8888q_      8888888888888888888      _p888P
   '"8888qq_  888888888888888888P  _pp888P"'
       ""8888p"88888888888888888'q888P""
         __"""_88888888888888888_"""_,
        _88P  "88888888888888888P  d88,
       .d8P    888888888888888"    88b,
       d88'      888888888888P"      '88b
                    '"""""'
     .___  ___  ___,___,.__,    ______,
      888 .888,.888 d88|:88b    d888888q,
      )88|q888p|88P d88|:88b    d88  '888
       88b88"88d88  d88|:88b    d88   888
       )888P d888(  d88|:88bqqp,d88pp888"
        ^^^' 'YY^   "YY''Y^YYY^'"^^^^""
       \_____, \__  .__   __pq_,  _pqq_
       |88"Y88,d88  |88 _88P^^YP 88P"Y8
       |88q88P d88  |88 88P qqqp O88qqq,
       |88''88pd88, )88 88b_')88 _,""88p
       )88888P' Y8888P"  "88888P 88888P'

`

const colorArt = `
             .88qq_         _pq88:
               '"88q,     _p8P"'
                  )8bppqqqp8P
        qq,     _p88888888888q,     \pp
        d8b    _888888888888888_    d8P
        '88_, _88888888888888888p  _88'
         "888q8888888888888888888p88P"
            _"""\888888888888888""_
         _p888 _8888888888888888p888q_
     __p888P" \88888888888888888b "Y888q_,
  _pp888P"    d888888888888888888    "Y888qq,
  8888q_      8888888888888888888      _p888P
   '"8888qq_  888888888888888888P  _pp888P"'
       ""8888p"88888888888888888'q888P""
         __"""_88888888888888888_"""_,
        _88P  "88888888888888888P  d88,
       .d8P    d888888888888888"    88b,
       d88'      "888888888888P"      '88b
                    '"""""'
     .___  ___  ___,___,.__,    ______,
      888 .888,.888 d88|:88b    d888888q,
      )88|q888p|88P d88|:88b    d88888888
       88b88"88d88  d88|:88b    d88888888
       )888P d888(  d88|:88bqqp,d8888888"
        ^^^' 'YY^   "YY''Y^YYY^'"^^^^""
       \_____, \__  .__   __pq_,  _pqq_
       |888888,d88  |88 _88P^^YP 88P"Y8
       |88888P d88  |88 88P qqqp O88qqq,
       |888888pd88, )88 88b_')88 _,""88p
       )888888P' Y8888P"  "88888P 88888P'

`

const inkArt = `
             011111         111113
               111111     111111
                  11111111111
        111     111111111111111     111
        111    11115511111551111    111
        11110 1111111111111111111  1111
         11111111111111111111111111111
            11111666661116666666111
         11111 66666661116666666611111
     11111111 6666666611166666666 11111111
  11111111    6666666611166666666    11111111
  111111      6666662211166666662      111111
   111111111  2222222211122222222  111111111
       111111122222221111122222201111111
         11111112222211111222221111111
        1111  1111222111112221111  1111
       0111    11111111111111111    1111
       1111      1111111111111      1111
                    0111111
     1111  111  111111100111    1111111
      111 111110111 11100111    111111111
      1111111111111 11100111    111555111
       11111111111  11100111    111555111
       11111 11111  111031111111111111111
        1111 1111   1110011111111111111
       2222222 222  222   222222  22222
       22252220222  222 22222222 222222
       2222222 222  222 222 2222 2222220
       222552222220 222 22222222 2022222
       22222222 2222222  2222222 2222222

`

func splitArt(art string) []string {
	art = strings.TrimSuffix(art, "\n")
	lines := strings.Split(art, "\n")
	for i := range lines {
		if len(lines[i]) < cols {
			lines[i] += strings.Repeat(" ", cols-len(lines[i]))
		}
	}
	for len(lines) < rows {
		lines = append(lines, strings.Repeat(" ", cols))
	}
	return lines
}

var monoLines = splitArt(monoArt)
var artLines = splitArt(colorArt)
var inkLines = splitArt(inkArt)

func hash(x, y, tick int) uint32 {
	h := uint32(x)*73856093 ^ uint32(y)*19349663 ^ uint32(tick)*83492791
	h = (h ^ (h >> 13)) * 1274126177
	return h ^ (h >> 16)
}

func parseHex(s string) (int, int, int) {
	s = strings.TrimPrefix(s, "#")
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return 255, 255, 255
	}
	return int((n >> 16) & 0xff), int((n >> 8) & 0xff), int(n & 0xff)
}

func fg(hex string) string {
	r, g, b := parseHex(hex)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func firstLast() (int, int) {
	first, last := rows, -1
	for _, pic := range [][]string{monoLines, artLines} {
		for y, line := range pic {
			if strings.TrimSpace(line) != "" {
				if y < first {
					first = y
				}
				if y > last {
					last = y
				}
			}
		}
	}
	if last < first {
		return 0, rows - 1
	}
	return first, last
}

var firstRow, lastRow = firstLast()

// frame builds one full terminal frame. "paper=false" selects the dark-page
// palette from the supplied ascii.rest script.
func frame(t float64) string {
	every := 3.0
	span := float64(lastRow-firstRow+1) + head + trail
	at := math.Inf(-1)
	since := t - start
	if since >= 0 {
		at = float64(firstRow) + span*math.Mod(since, every)/pass
	}
	tick := int(math.Floor((t * fps) / 2))
	var out strings.Builder
	out.Grow(rows * (cols*12 + 16))

	for y := 0; y < rows; y++ {
		d := at - float64(y)
		level := 0
		if d >= 0 && d < head {
			level = 2
		} else if d >= head && d < head+trail {
			level = 1
		}
		line := artLines[y]
		ink := inkLines[y]
		for x := 0; x < cols; x++ {
			ch := byte(' ')
			if x < len(line) {
				ch = line[x]
			}
			if ch != ' ' {
				if level == 2 {
					h := hash(x, y, tick)
					if h%4 != 0 {
						ch = noise[int((h>>2)%uint32(len(noise)))]
					}
				}
				inkValue := byte('0')
				if x < len(ink) {
					inkValue = ink[x]
				}
				index, err := strconv.ParseInt(string(inkValue), 36, 32)
				if err != nil {
					index = 0
				}
				// The original renderer uses palette offset 24 for dark backgrounds.
				paletteIndex := 24 + level*8 + int(index)
				if paletteIndex >= 0 && paletteIndex < len(palette) {
					out.WriteString(fg(palette[paletteIndex]))
				}
				out.WriteByte(ch)
			} else {
				out.WriteByte(' ')
			}
		}
		out.WriteString("\x1b[0m")
		if y < rows-1 {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func wildbugsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming is not supported by this server", http.StatusInternalServerError)
		return
	}

	// A request may use ?once=1 to receive a single frame for diagnostics.
	if r.URL.Query().Get("once") == "1" {
		_, _ = fmt.Fprint(w, "\x1b[H\x1b[2J", frame(0), "\n")
		flusher.Flush()
		return
	}

	ticker := time.NewTicker(time.Second / fps)
	defer ticker.Stop()
	started := time.Now()

	for {
		select {
		case <-r.Context().Done():
			return
		case now := <-ticker.C:
			elapsed := now.Sub(started).Seconds()
			if _, err := fmt.Fprint(w, "\x1b[H\x1b[2J", frame(elapsed), "\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/wildbugs", wildbugsHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintln(w, "Wild Bugs ASCII animation server")
		_, _ = fmt.Fprintln(w, "")
		_, _ = fmt.Fprintln(w, "Play the animation:")
		_, _ = fmt.Fprintln(w, "  curl -N https://YOUR-SERVICE.onrender.com/wildbugs")
		_, _ = fmt.Fprintln(w, "")
		_, _ = fmt.Fprintln(w, "Health check: /healthz")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Wild Bugs ASCII server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
