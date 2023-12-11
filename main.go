// convert uid and gid between mac and linux during rclone sync, using --metadata-mapper
//
// usage: rclone sync source:path dest:path --metadata-mapper /path/to/rclone-permissions-mapper
// or, to see input and output:
// rclone sync source:path dest:path --metadata-mapper /path/to/rclone-permissions-mapper -v --dump mapper
// https://rclone.org/docs/#metadata-mapper
// https://github.com/nielash/rclone-permissions-mapper

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	// Read the input
	var in map[string]any
	err := json.NewDecoder(os.Stdin).Decode(&in)
	if err != nil {
		log.Fatal(err)
	}

	// Check the input
	metadata, ok := in["Metadata"]
	if !ok {
		fmt.Fprintf(os.Stderr, "Metadata key not found\n")
		os.Exit(1)
	}

	// Map the metadata
	metadataOut := map[string]string{}
	var out = map[string]any{
		"Metadata": metadataOut,
	}

	// debug settings
	debugging := false // set true to debug
	dir := "/"
	debugfile := ""
	var df *os.File
	if debugging {
		dir, err = os.UserHomeDir()
		if err == nil {
			debugfile = filepath.Join(dir, "rclone-permissions-mapper-debug.txt")
		}
		df, err = os.Create(debugfile)
		defer df.Close()
	}
	debug := func(format string, a ...any) {
		if !debugging {
			return
		}
		fmt.Fprintf(df, format, a...)
	}

	// loop through the metadata keys
	for k, v := range metadata.(map[string]any) {
		switch k {
		case "error":
			fmt.Fprintf(os.Stderr, "Error: %s\n", v)
			os.Exit(1)
		case "uid":
			debug("uid detected! key: %s, val: %v\n", k, v)
			uidstr, ok := v.(string)
			osuid := os.Getuid()
			debug("osuid: %v\n", osuid)
			if ok {
				uid, err := strconv.Atoi(uidstr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error converting string to int: %v\n", uidstr)
				}
				// mac is 501 - 999, linux is 1000+
				if (uid >= 1000 && osuid >= 501 && osuid < 1000) || (osuid >= 1000 && uid >= 501 && uid < 1000) {
					// unset it
					debug("unsetting. key: %s, val: %v\n", k, v)
					continue
				}
			}
		case "gid":
			debug("gid detected! key: %s, val: %v\n", k, v)
			gidstr, ok := v.(string)
			if ok {
				gid, err := strconv.Atoi(gidstr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error converting string to int: %v\n", gidstr)
				}
				// mac is 20, linux is 1000
				debug("osgid: %v\n", os.Getgid())
				if gid != os.Getgid() {
					// unset it
					debug("unsetting. key: %s, val: %v\n", k, v)
					continue
				}
			}
		default:
			debug("skipping -- key: %s, val: %v\n", k, v)
		}
		metadataOut[k] = v.(string)
	}

	// Write the output
	json.NewEncoder(os.Stdout).Encode(&out)
	if err != nil {
		log.Fatal(err)
	}

	if debugging {
		debug("final: \n")
		bytes, err := json.MarshalIndent(&out, "", "\t")
		debug("%v", string(bytes))
		if err != nil {
			debug("json err: %v", err.Error())
		}
	}
}
