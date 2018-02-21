package api

import(
	"github.com/olivere/elastic"
	"net/http"
	"golang.org/x/net/context"
	"encoding/json"
	"io"
	"strings"
	"github.com/Open-Android/webserver/utils"
	"log"
)

func Query(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	args := req.URL.Query()
        toArg, _ := utils.GetArg("to", args)
        fromArg, _ := utils.GetArg("from", args)
	flds, _ := utils.GetArg("fields", args)
        //http.Error(w, InfoMsg, http.StatusBadRequest)

	fields := []string{}
	if len(flds) > 0 {
		fields = strings.Split(flds, ",")
	}

	rangeQuery := elastic.NewRangeQuery("Date")
	if len(fromArg) == 0 {
		rangeQuery = rangeQuery.From(nil)
	} else {
		rangeQuery = rangeQuery.From(fromArg)
	}
	if len(toArg) == 0 {
		rangeQuery = rangeQuery.To(nil)
	} else {
		rangeQuery = rangeQuery.To(toArg)
	}

	log.Println("===============================")
	log.Println("From: " + fromArg)
	log.Println("To: " + toArg)
	log.Println("Fields: ")
	log.Println(fields)
	log.Println("===============================")

	scroll := client.Scroll("apks").Type("apk").Query(rangeQuery).Size(100)
	io.WriteString(w, "{\"data\":[")
	for {
		results, err := scroll.Do(ctx)
		if err == io.EOF {
			break// all results retrieved
		}
		if err != nil {
			log.Println(err.Error())
			break
		}
			// Send the hits to the hits channel
		for _, hit := range results.Hits.Hits {
			if len(fields) == 0 {
				w.Write(*hit.Source)
				io.WriteString(w, ",")
				continue
			}
			data := make(map[string]interface{})
                        out := make(map[string]interface{})
			err := json.Unmarshal(*hit.Source, &data)
                        if err != nil {
				panic(err)
                        }
			for key, val := range data {
				if utils.StringInSlice(key, fields) {
					out[key] = val
				}
			}
			writer := []byte{}
                        writer, err = json.Marshal(out)
                        if err != nil {
				panic(err)
                        }
                        w.Write(writer)
			io.WriteString(w, ",")
		}
	}
	io.WriteString(w, "]}")

	// This example illustrates how to use goroutines to iterate
	// through a result set via ScrollService.
	//
	// It uses the excellent golang.org/x/sync/errgroup package to do so.
	//
	// The first goroutine will Scroll through the result set and send
	// individual documents to a channel.
	//
	// The second cluster of goroutines will receive documents from the channel and
	// deserialize them.
	//
	// Feel free to add a third goroutine to do something with the
	// deserialized results.
	//
	// Let's go.

	// 1st goroutine sends individual hits to channel.
	/*hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(hits)
		// Initialize scroller. Just don't call Do yet.
		scroll := client.Scroll("apks").Type("apk").Size(100)
		for {
			results, err := scroll.Do(ctx)
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				return err // something went wrong
			}

			// Send the hits to the hits channel
			for _, hit := range results.Hits.Hits {
				select {
				case hits <- *hit.Source:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	})

	// 2nd goroutine receives hits and deserializes them.
	//
	// If you want, setup a number of goroutines handling deserialization in parallel.
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for hit := range hits {
				data := &ApkData{}
				err := json.Unmarshal(hit, &data)
				if err != nil {
					return err
				}
				writer := []byte{}
				writer, err = json.Marshal(data)
				if err != nil {
					return err
				}
				w.Write(writer)
				select {
				default:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		panic(err)
	}*/
}
