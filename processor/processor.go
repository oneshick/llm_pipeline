package processor

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"pipeline/llm"
	"pipeline/orm"
)

type Processor struct {
	client  *llm.GroqClient
	delayMs int
}

func New(client *llm.GroqClient, delayMs int) *Processor {
	return &Processor{
		client:  client,
		delayMs: delayMs,
	}
}

func buildPrompt(p models.Product) string {
	return fmt.Sprintf(`Ты — ассистент по извлечению данных о товарах.

Верни только JSON без markdown и пояснений.

Формат:
{
  "name": "полное название товара",
  "brand": "бренд",
  "category": "категория",
  "price": 12345,
  "currency": "RUB",
  "key_specs": ["характеристика 1", "характеристика 2"]
}

Правила:
- price — только число
- currency всегда RUB
- key_specs: 3-5 характеристик
- если данных нет — пустая строка или 0

Описание товара:
%s`, p.Description)
}

func cleanJSON(raw string) string {
	raw = strings.TrimSpace(raw)

	if i := strings.Index(raw, "{"); i > 0 {
		raw = raw[i:]
	}

	if i := strings.LastIndex(raw, "}"); i >= 0 {
		raw = raw[:i+1]
	}

	return raw
}

func (p *Processor) Process(product models.Product) (models.ProductFeatures, error) {
	const retries = 3

	var raw string
	var err error

	for attempt := 1; attempt <= retries; attempt++ {
		raw, err = p.client.Complete(buildPrompt(product))

		if err == nil {
			break
		}

		if attempt < retries {
			wait := time.Duration(attempt*attempt) * time.Second

			fmt.Printf(
				"повтор %d/%d через %v (%v)\n",
				attempt,
				retries,
				wait,
				err,
			)

			time.Sleep(wait)
		}
	}

	if err != nil {
		return models.ProductFeatures{},
			fmt.Errorf("все попытки исчерпаны: %w", err)
	}

	var out models.ProductFeatures

	err = json.Unmarshal([]byte(cleanJSON(raw)), &out)
	if err != nil {
		return models.ProductFeatures{},
			fmt.Errorf("парсинг JSON: %w\nraw: %s", err, raw)
	}

	out.ID = product.ID

	return out, nil
}

func (p *Processor) ProcessAll(
	products []models.Product,
) ([]models.ProductFeatures, int) {

	var results []models.ProductFeatures
	failed := 0

	for i, product := range products {
		fmt.Printf(
			"[%d/%d] Товар #%s\n",
			i+1,
			len(products),
			product.ID,
		)

		features, err := p.Process(product)

		if err != nil {
			fmt.Printf(" %v\n", err)
			failed++
		} else {
			fmt.Printf(
				" %s — %s (%.0f ₽)\n",
				features.Brand,
				features.Name,
				features.Price,
			)

			results = append(results, features)
		}

		if i < len(products)-1 {
			time.Sleep(time.Duration(p.delayMs) * time.Millisecond)
		}
	}

	return results, failed
}
