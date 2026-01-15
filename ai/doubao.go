package ai

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model/responses"
)

const (
	doubaoBaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	doubaoModel   = "ep-20260115103448-9764k"
)

func Ask(ctx context.Context, prompt string) (string, error) {
	apiKey := os.Getenv("ARK_API_KEY")
	if apiKey == "" {
		return "", errors.New("missing ARK_API_KEY")
	}

	client := arkruntime.NewClientWithApiKey(
		apiKey,
		arkruntime.WithBaseUrl(doubaoBaseURL),
	)

	inputMessage := &responses.ItemInputMessage{
		Role: responses.MessageRole_user,
		Content: []*responses.ContentItem{
			{
				Union: &responses.ContentItem_Text{
					Text: &responses.ContentItemText{
						Type: responses.ContentItemType_input_text,
						Text: prompt,
					},
				},
			},
		},
	}

	resp, err := client.CreateResponses(ctx, &responses.ResponsesRequest{
		Model: doubaoModel,
		Input: &responses.ResponsesInput{
			Union: &responses.ResponsesInput_ListValue{
				ListValue: &responses.InputItemList{ListValue: []*responses.InputItem{{
					Union: &responses.InputItem_InputMessage{
						InputMessage: inputMessage,
					},
				}}},
			},
		},
	})
	if err != nil {
		return "", err
	}

	return extractText(resp), nil
}

func extractText(resp *responses.ResponseObject) string {
	if resp == nil {
		return ""
	}

	var parts []string
	for _, item := range resp.Output {
		output := item.GetOutputMessage()
		if output == nil {
			continue
		}
		for _, content := range output.Content {
			text := content.GetText()
			if text == nil || text.Text == "" {
				continue
			}
			parts = append(parts, text.Text)
		}
	}

	return strings.Join(parts, "")
}
