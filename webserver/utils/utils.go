package utils

import(
	"fmt"
	"os"
)

const InfoMsg string = `
No GET parameters found!

Here's some brief info on how to query:

q=
     - Query search parameter (Required)
     - If you'd like to download everything, use q=*:*
     - Search using regex or for a particular text with q=<field name>:<regex or literal>
     - Example: q=Malicious:true
        - Returns all APKs that are malicious

fl=
     - Return only these fields
     - Example: fl=Permissions+Apis+PackageName
        - Returns only the fields Permissions, Apis, and PackageName

omitHeader=
     - Remove extra query info
     - Example: omitHeader=true

rows=
     - Number of entries to return
     - Example: rows=7000
         - Returns top 7000 entries of given query

wt=
     - Writer Type
     - Type of file that is returned from the query
     - json, csv, xml, etc.
     - Default of json
     - Example: wt=json

More detailed descriptions and examples can be seen on our wiki: https://github.com/rschmicker/Open-Android/wiki

`

func GetArg(name string, aMap map[string][]string) (arg string, err error) {
        var ok bool
        var arglist []string
        if arglist, ok = aMap[name]; !ok {
                return "", fmt.Errorf("unable to obtain arg from map with key %s", name)
        }
        return arglist[0], nil
}

func PrintUsage() {
        fmt.Println(`
Syntax:
        >webserver -key <Directory to HTTPS key> -cert <Directory to HTTPS certificate>

Example:
        >webserver -key ./keys/server.key -cert ./keys/server.crt
`)
        os.Exit(1)
}

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


