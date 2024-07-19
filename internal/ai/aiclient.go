package ai

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"cloud.google.com/go/vertexai/genai"
)

type Client struct {
	client         *genai.Client
	primaryModel   *genai.GenerativeModel
	secondaryModel *genai.GenerativeModel

	template *template.Template
}

const region = "asia-northeast1"
const primaryModel = "gemini-1.5-pro-001"
const secondaryModel = "gemini-1.0-pro-002"
const primaryTemp = 1.0
const secondaryTemp = 1.0
const templatePath = "template/prompt.txt"

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	client, err := genai.NewClient(ctx, projectID, region)
	if err != nil {
		return nil, err
	}
	newModel := client.GenerativeModel(primaryModel)
	newModel.SetTemperature(primaryTemp)

	oldModel := client.GenerativeModel(secondaryModel)
	oldModel.SetTemperature(secondaryTemp)

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &Client{client: client,
		primaryModel: newModel, secondaryModel: oldModel, template: t}, nil
}

// Ask returns hot and cold answers using theme from the LLM
func (c *Client) Ask(ctx context.Context, theme string) (string, string, error) {

	newAns, err := c.requestLLM(ctx, c.primaryModel, theme)
	if err != nil {
		return "", "", err
	}

	oldAns, err := c.requestLLM(ctx, c.secondaryModel, theme)
	if err != nil {
		return "", "", err
	}

	return newAns, oldAns, nil
}

func (c *Client) requestLLM(ctx context.Context, model *genai.GenerativeModel, theme string) (string, error) {
	var b bytes.Buffer
	c.template.Execute(&b, theme)

	prompt := genai.Text(b.String())

	resp, err := model.GenerateContent(ctx, prompt)
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
