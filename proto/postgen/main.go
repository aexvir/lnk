package main

import (
	"bytes"
	"os"
)

// post gen hooks or smth like that
func main() {
	fixgopbfile()
	renamegrpcpbfile()
}

// hacky function to fix broken import caused by experimental gnostic buf repository
// context: https://github.com/google/gnostic/issues/337
func fixgopbfile() {
	target := "lnk.pb.go"
	input, err := os.ReadFile(target)
	if err != nil {
		panic(err)
	}

	stat, err := os.Stat(target)
	if err != nil {
		panic(err)
	}

	output := bytes.Replace(
		input,
		[]byte(`"./openapiv3"`),
		[]byte(`"github.com/google/gnostic/openapiv3"`),
		1,
	)

	err = os.WriteFile(target, output, stat.Mode())
	if err != nil {
		panic(err)
	}
}

// renames the lnk grpc file so it matches the rest of the generated files
// no real value, I just like it more :D
func renamegrpcpbfile() {
	old := "lnk_grpc.pb.go"
	new := "lnk.pb.grpc.go"

	stat, err := os.Stat(old)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(old)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(new, content, stat.Mode())
	if err != nil {
		panic(err)
	}

	err = os.Remove(old)
	if err != nil {
		panic(err)
	}
}
