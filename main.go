package main

import (
	"fmt"
	"log"
	"pipeline/config"
	"pipeline/llm"
	"pipeline/orm"
	"pipeline/processor"
	"pipeline/reader"
	"pipeline/writer"
)

func main() {
	cfg := config.Load()

	fmt.Println("")
	fmt.Println("LLM Pipeline | Groq +", cfg.Model)
	fmt.Println("")
	fmt.Printf("CSV:     %s\n", cfg.CSVPath)
	fmt.Printf("Output:  %s\n", cfg.OutPath)
	fmt.Printf("Delay:   %dms\n", cfg.DelayMs)
	fmt.Println("")

	fmt.Printf("\nЧитаем %s...\n", cfg.CSVPath)

	products, err := reader.New(cfg.CSVPath).Read()
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("Загружено %d товаров\n\n", len(products))

	client := llm.New(cfg.GroqAPIKey, cfg.Model, cfg.MaxTokens)
	proc := processor.New(client, cfg.DelayMs)

	features, failed := proc.ProcessAll(products)

	result := models.PipelineResult{
		Model:      cfg.Model,
		TotalItems: len(products),
		Failed:     failed,
		Products:   features,
	}

	fmt.Printf("\nСохраняем в %s...\n", cfg.OutPath)

	if err := writer.New(cfg.OutPath).Write(result); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("Успешно: %d товаров\n", len(features))
	fmt.Printf("Ошибок:  %d товаров\n", failed)
	fmt.Printf("Файл:    %s\n", cfg.OutPath)
}
