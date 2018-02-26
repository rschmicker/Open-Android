package main

import(
        "github.com/olivere/elastic"
        "golang.org/x/net/context"
        "encoding/json"
        "io"
        "os"
        "strings"
        "log"
        "time"
	"flag"
)

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func main() {
	fieldsFlag := flag.String("fields", "", "fields to query")
	toFlag := flag.String("to", "", "to timestamp")
	fromFlag := flag.String("from", "", "from timestamp")
	flag.Parse()

	client, err := elastic.NewClient()
        if err != nil {
                panic(err)
        }

	ctx := context.Background()

        toArg := *toFlag
        fromArg := *fromFlag
        flds := *fieldsFlag

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

        filename := time.Now().Format(time.RFC850) + ".json"
        f, err := os.Create("/iscsi/queries/" + filename)
	if err != nil {
                log.Println(err)
                return
        }
        defer f.Close()

        io.WriteString(f, "{\"data\":[")
        first := true
        for {
                results, err := scroll.Do(ctx)
                if err == io.EOF {
                        break// all results retrieved
                } else if !first {
                        io.WriteString(f, ",")
                }
                first = false
                if err != nil {
                        log.Println(err.Error())
                        break
                }
                        // Send the hits to the hits channel
                for _, hit := range results.Hits.Hits[:len(results.Hits.Hits)-1] {
                        if len(fields) == 0 {
                                f.Write(*hit.Source)
                                io.WriteString(f, ",")
                                continue
                        }
                        data := make(map[string]interface{})
                        out := make(map[string]interface{})
                        err := json.Unmarshal(*hit.Source, &data)
                        if err != nil {
                                panic(err)
                        }
                        for key, val := range data {
                                if StringInSlice(key, fields) {
                                        out[key] = val
                                }
                        }
                        writer := []byte{}
                        writer, err = json.Marshal(out)
			if err != nil {
                                panic(err)
                        }
                        f.Write(writer)
                        io.WriteString(f, ",")
                }
                data := make(map[string]interface{})
                out := make(map[string]interface{})
                err = json.Unmarshal(*results.Hits.Hits[len(results.Hits.Hits)-1].Source, &data)
                if err != nil {
                        panic(err)
                }
                for key, val := range data {
                        if StringInSlice(key, fields) {
                                out[key] = val
                        }
                }
                writer := []byte{}
                writer, err = json.Marshal(out)
                if err != nil {
                        panic(err)
                }
                f.Write(writer)
	}
        io.WriteString(f, "]}")
}
