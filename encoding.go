package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Order struct {
	ID          string    `json:"id"`
	DateOrdered time.Time `json:"date_ordered"`
	CustomerID  string    `json:"customer_id"`
	Items       []Item    `json:"items"`
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// custom MarshalJSON() for Order
// because we want to use custom logic to marshal DateOrdered field
// this is the recommended approach
func (o Order) MarshalJSON() ([]byte, error) {
	type Dup Order

	tmp := struct {
		DateOrdered string `json:"date_ordered"`
		Dup
	}{
		Dup: (Dup)(o),
	}

	tmp.DateOrdered = o.DateOrdered.Format(time.RFC822Z)
	b, err := json.Marshal(tmp)

	return b, err
}

// custom MarshalJSON() for Order
// because we want to use custom logic to marshal DateOrdered field
// this is the recommended approach
func (o *Order) UnmarshalJSON(b []byte) error {
	type Dup Order

	tmp := struct {
		DateOrdered string `json:"date_ordered"`
		*Dup
	}{
		Dup: (*Dup)(o),
	}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	o.DateOrdered, err = time.Parse(time.RFC822Z, tmp.DateOrdered)
	if err != nil {
		return err
	}

	return nil
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RFC288ZTime struct {
	time.Time
}

// this is not the recommended approach
func (rt RFC288ZTime) MarshalJSON() ([]byte, error) {
	out := rt.Time.Format(time.RFC822Z)
	return []byte(`"` + out + `"`), nil
}

// this is not the recommended approach
func (rt *RFC288ZTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}

	t, err := time.Parse(`"`+time.RFC822Z+`"`, string(b))
	if err != nil {
		return err
	}

	*rt = RFC288ZTime{t}
	return nil
}

func countLetters(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048)
	out := make(map[string]int)

	for {
		n, err := r.Read(buf)
		for _, b := range buf[:n] {
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func buildGzipReader(fileName string) (*gzip.Reader, func(), error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}

	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, nil, err
	}

	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}

func encoding() {
	streamData := `
    {"name": "Fred", "age": 40}
    {"name": "Mary", "age": 21}
    {"name": "Pat", "age": 30}
  `

	var t struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	dec := json.NewDecoder(strings.NewReader(streamData))

	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	// example of reading stream data and decoding it
	for {
		err := dec.Decode(&t)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		// process t
		fmt.Println(t.Name)

		err = enc.Encode(t)
		if err != nil {
			panic(err)
		}

	}

	out := b.String()
	fmt.Println(out)
}
