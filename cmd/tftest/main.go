package main

import (
	"flag"
	"fmt"

	tf "github.com/wamuir/graft/tensorflow"
)

func main() {
	path := flag.String("path", "", "model path")
	flag.Parse()

	model, err := tf.LoadSavedModel(*path, []string{"serve"}, nil)
	if err != nil {
		panic(err)
	}

	out := model.Graph.Operation("StatefulPartitionedCall_1").Output(0)

	in := map[string]int{
		"first":     1,
		"skill1":    -1,
		"skill2":    -1,
		"skill3":    -1,
		"skill4":    -1,
		"skill5":    -1,
		"skill6":    -1,
		"skill7":    -1,
		"skill8":    -1,
		"skill9":    -1,
		"oppskill1": -1,
		"oppskill2": -1,
		"oppskill3": -1,
		"oppskill4": -1,
		"oppskill5": -1,
		"oppskill6": -1,
		"oppskill7": -1,
		"oppskill8": -1,
		"oppskill9": -1,
	}

	for range 100 {
		feeds, err := calcFeeds(model.Graph, in)
		if err != nil {
			panic(err)
		}
		fetches := []tf.Output{out}

		res, err := model.Session.Run(feeds, fetches, nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", res[0].Value())
	}
}

func calcFeeds(g *tf.Graph, in map[string]int) (map[tf.Output]*tf.Tensor, error) {
	res := make(map[tf.Output]*tf.Tensor)

	for key, val := range in {
		name := "serving_default_" + key
		output := g.Operation(name).Output(0)
		tensor, err := tf.NewTensor([1][1]int64{{int64(val)}})
		if err != nil {
			return nil, err
		}

		res[output] = tensor
	}

	return res, nil
}
