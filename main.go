package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	docopt "github.com/docopt/docopt-go"
	api "github.com/hashicorp/consul/api"
)

//type KVPair struct {
//    Key         string
//    CreateIndex uint64
//    ModifyIndex uint64
//    LockIndex   uint64
//    Flags       uint64
//    Value       []byte
//    Session     string
//}

type ByCreateIndex api.KVPairs

func (a ByCreateIndex) Len() int      { return len(a) }
func (a ByCreateIndex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

//Sort the KVs by createIndex
func (a ByCreateIndex) Less(i, j int) bool { return a[i].CreateIndex < a[j].CreateIndex }

func backup(ipaddress string, token string) {

	config := api.DefaultConfig()
	config.Address = ipaddress
	config.Token = token

	client, _ := api.NewClient(config)
	kv := client.KV()

	pairs, _, err := kv.List("/", nil)
	if err != nil {
		panic(err)
	}

	sort.Sort(ByCreateIndex(pairs))

	outstring := ""
	for _, element := range pairs {
		encoded_value := base64.StdEncoding.EncodeToString(element.Value)
		outstring += fmt.Sprintf("%s:%s\n", element.Key, encoded_value)
	}

	fmt.Print(outstring)
}

func backupAcls(ipaddress string, token string) {

	config := api.DefaultConfig()
	config.Address = ipaddress
	config.Token = token

	client, _ := api.NewClient(config)
	acl := client.ACL()

	tokens, _, err := acl.List(nil)
	if err != nil {
		panic(err)
	}
	// sort.Sort(ByCreateIndex(tokens))

	outstring := ""
	for _, element := range tokens {
		// outstring += fmt.Sprintf("%s:%s:%s:%s\n", element.ID, element.Name, element.Type, element.Rules)
		outstring += fmt.Sprintf("====\nID: %s\nName: %s\nType: %s\nRules:\n%s\n", element.ID, element.Name, element.Type, element.Rules)
	}

	fmt.Print(outstring)
}

/* File needs to be in the following format:
   KEY1:VALUE1
   KEY2:VALUE2
*/
func restore(ipaddress string, token string, infile string) {

	config := api.DefaultConfig()
	config.Address = ipaddress
	config.Token = token

	file := fmt.Sprintf("/restore/%s", infile)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	client, _ := api.NewClient(config)
	kv := client.KV()

	for _, element := range strings.Split(string(data), "\n") {
		kvp := strings.Split(element, ":")

		if len(kvp) > 1 {
			//log.Printf("Encoded: %s\n", kvp[1])
			decoded_value, decode_err := base64.StdEncoding.DecodeString(kvp[1])
			if decode_err != nil {
				panic(decode_err)
			}
			//log.Printf("Decoded: %s\n", decoded_value)

			p := &api.KVPair{Key: kvp[0], Value: decoded_value}

			_, err := kv.Put(p, nil)
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {

	usage := `Consul KV and ACL Backup with KV Restore tool.

Usage:
  consul-backup [-i IP:PORT] [-t TOKEN] [--aclbackup] [--restore <filename>]
  consul-backup -h | --help
  consul-backup --version

Options:
  -h --help                          Show this screen.
  --version                          Show version.
  -i, --address=IP:PORT              The HTTP endpoint of Consul [default: 127.0.0.1:8500].
  -t, --token=TOKEN                  An ACL Token with proper permissions in Consul [default: ].
  -a, --aclbackup                    Backup ACLs, does nothing in restore mode. ACL restore not available at this time.
  -b, --aclbackupfile=ACLBACKUPFILE  ACL Backup Filename [default: acl.bkp].
  -r, --restore                      Activate restore mode`

	arguments, _ := docopt.Parse(usage, nil, true, "consul-backup 1.0", false)

	if arguments["--restore"] == true {
		log.Print("Restore mode:")
		log.Print("Restoring KV from file: ", arguments["<filename>"].(string))
		restore(arguments["--address"].(string), arguments["--token"].(string), arguments["<filename>"].(string))
	} else {
		log.Print("Backup mode:")
		backup(arguments["--address"].(string), arguments["--token"].(string))
		if arguments["--aclbackup"] == true {
			log.Print("ACL Tokens will be backed up to file: ", arguments["--aclbackupfile"].(string))
			backupAcls(arguments["--address"].(string), arguments["--token"].(string))
		}
	}
}
