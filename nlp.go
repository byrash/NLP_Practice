package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/jdkato/prose.v2"
)

type ProdigyOutput struct {
	Text   string
	Spans  []prose.LabeledEntity
	Answer string
}

func readProdigy(jsonLines []byte) []ProdigyOutput {
	decoder := json.NewDecoder(bytes.NewReader(jsonLines))
	entries := []ProdigyOutput{}
	for {
		ent := ProdigyOutput{}
		err := decoder.Decode(&ent)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(errors.Wrap(err, "Error parsing JSONL line"))
		}
		entries = append(entries, ent)
	}
	return entries
}

func split(data []ProdigyOutput) ([]prose.EntityContext, []ProdigyOutput) {
	cutoff := int(float64(len(data)) * 0.8)
	trainingData, testData := []prose.EntityContext{}, []ProdigyOutput{}
	for i, entry := range data {
		if i < cutoff {
			trainingData = append(trainingData, prose.EntityContext{
				Text:   entry.Text,
				Spans:  entry.Spans,
				Accept: strings.EqualFold(entry.Answer, "accept"),
			})
		} else {
			testData = append(testData, entry)
		}
	}
	return trainingData, testData
}

func main() {
	// data, err := ioutil.ReadFile("reddit_product.jsonl")
	// if err != nil {
	// 	panic(errors.Wrap(err, "Unable to opne JOSNL file"))
	// }
	// trainingData, testData := split(readProdigy(data))
	// log.Printf("Training with %d and testing with %d entries.\n",
	// 	len(trainingData), len(testData))
	// model := prose.ModelFromData("PRODUCT", prose.UsingEntities(trainingData))
	// correct := 0.0
	// for _, entry := range testData {
	// 	proseDoc, err := prose.NewDocument(
	// 		entry.Text,
	// 		prose.WithSegmentation(false),
	// 		prose.UsingModel(model))
	// 	if err != nil {
	// 		panic(errors.Wrap(err, "Problem with creating prose document"))
	// 	}
	// 	entities := proseDoc.Entities()
	// 	if entry.Answer != "accept" && len(entities) == 0 {
	// 		correct++
	// 	} else {
	// 		expected := []string{}
	// 		for _, span := range entry.Spans {
	// 			expected = append(expected, entry.Text[span.Start:span.End])
	// 		}
	// 		if reflect.DeepEqual(expected, entities) {
	// 			correct++
	// 		}
	// 	}
	// }
	// log.Printf("Correct (%%): %f\n", correct/float64(len(testData)))
	// model.Write("PRODUCT")

	model := prose.ModelFromDisk("PRODUCT")

	for {
		log.Println("Enter the string to identify the product:")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		log.Printf("String you supplied %v", text)
		doc, err := prose.NewDocument(
			text,
			prose.WithSegmentation(false),
			prose.UsingModel(model))

		if err != nil {
			panic(errors.Wrap(err, "Error with prose"))
		}
		log.Println("Results --")
		for _, ent := range doc.Entities() {
			log.Println(ent.Text, ent.Label)
		}
		log.Println()
	}

}
