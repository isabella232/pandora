//go:generate ffjson $GOFILE

package ammo

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/yandex/pandora/config"
	"golang.org/x/net/context"
)

// ffjson: noencoder
type Http struct {
	Host    string
	Method  string
	Uri     string
	Headers map[string]string
	Tag     string
}

func (h *Http) Request() (*http.Request, error) {
	// FIXME: something wrong here with https
	req, err := http.NewRequest(h.Method, "http://"+h.Host+h.Uri, nil)
	if err == nil {
		for k, v := range h.Headers {
			req.Header.Set(k, v)
		}
	}
	return req, err
}

// HttpJSONDecoder implements ammo.Decoder interface
type HttpJSONDecoder struct {
	pool sync.Pool
}

func (d *HttpJSONDecoder) Decode(jsonDoc []byte) (Ammo, error) {
	a := d.pool.Get().(*Http)
	err := a.UnmarshalJSON(jsonDoc)
	return a, err
}

// be polite and return unused Ammo to the pool
// be shure that you return Http because we don't make any checks here
func (d *HttpJSONDecoder) Release(a Ammo) {
	d.pool.Put(a)
}

func NewHttpJSONDecoder() Decoder {
	return &HttpJSONDecoder{
		pool: sync.Pool{
			New: func() interface{} {
				return &Http{}
			},
		},
	}
}

// ffjson: skip
type HttpProvider struct {
	*BaseProvider

	sink         chan<- Ammo
	ammoFileName string
	ammoLimit    int
	passes       int
}

func (ap *HttpProvider) Start(ctx context.Context) error {
	defer close(ap.sink)
	ammoFile, err := os.Open(ap.ammoFileName)
	if err != nil {
		return fmt.Errorf("failed to open ammo source: %v", err)
	}
	defer ammoFile.Close()
	ammoNumber := 0
	passNum := 0
loop:
	for {
		passNum++
		scanner := bufio.NewScanner(ammoFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() && (ap.ammoLimit == 0 || ammoNumber < ap.ammoLimit) {
			data := scanner.Bytes()
			if a, err := ap.Decode(data); err != nil {
				return fmt.Errorf("failed to decode ammo: %v", err)
			} else {
				ammoNumber++
				select {
				case ap.sink <- a:
				case <-ctx.Done():
					break loop
				}
			}
		}
		if ap.passes != 0 && passNum >= ap.passes {
			break
		}
		ammoFile.Seek(0, 0)
		if ap.passes == 0 {
			log.Printf("Restarted ammo from the beginning. Infinite passes.\n")
		} else {
			log.Printf("Restarted ammo from the beginning. Passes left: %d\n", ap.passes-passNum)
		}
	}
	log.Println("Ran out of ammo")
	return nil
}

func NewHttpProvider(c *config.AmmoProvider) (Provider, error) {
	ammoCh := make(chan Ammo)
	ap := &HttpProvider{
		ammoLimit:    c.AmmoLimit,
		passes:       c.Passes,
		ammoFileName: c.AmmoSource,
		sink:         ammoCh,
		BaseProvider: NewBaseProvider(
			ammoCh,
			NewHttpJSONDecoder(),
		),
	}
	return ap, nil
}
