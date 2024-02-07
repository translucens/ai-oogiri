package ai

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"cloud.google.com/go/vertexai/genai"
)

type Client struct {
	client    *genai.Client
	hotModel  *genai.GenerativeModel
	coldModel *genai.GenerativeModel

	template *template.Template
}

const region = "asia-northeast1"
const model = "gemini-pro"
const highTemp = 0.8
const lowTemp = 0.2
const templatePath = "template/prompt.txt"

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	client, err := genai.NewClient(ctx, projectID, region)
	if err != nil {
		return nil, err
	}
	hotModel := client.GenerativeModel(model)
	hotModel.SetTemperature(highTemp)

	coldModel := client.GenerativeModel(model)
	coldModel.SetTemperature(lowTemp)

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &Client{client: client,
		hotModel: hotModel, coldModel: coldModel, template: t}, nil
}

// Ask returns hot and cold answers using theme from the LLM
func (c *Client) Ask(ctx context.Context, theme string) (string, string, error) {

	hotAns, err := c.requestLLM(ctx, c.hotModel, theme)
	if err != nil {
		return "", "", err
	}

	coldAns, err := c.requestLLM(ctx, c.coldModel, theme)
	if err != nil {
		return "", "", err
	}

	return hotAns, coldAns, nil
}

func (c *Client) requestLLM(ctx context.Context, model *genai.GenerativeModel, theme string) (string, error) {
	var b bytes.Buffer
	c.template.Execute(&b, theme)

	prompt := genai.Text(b.String())

	resp, err := c.hotModel.GenerateContent(ctx, prompt)
	if err != nil {
		return "", err
	}

	if resp.PromptFeedback != nil {
		return resp.PromptFeedback.BlockReasonMessage, nil
	}

	resStr := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		switch v := part.(type) {
		case genai.Text:
			resStr += string(v)
		default:
			return "", fmt.Errorf("unexpected type: %T", v)
		}
	}
	return resStr, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
