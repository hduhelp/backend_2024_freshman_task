package utils

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func GenerateAIAnswer(question string) (string, error) {
	client := arkruntime.NewClientWithApiKey(
		"ec989009-e9d3-4600-bb9a-68457a8f5e1b",
	)

	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := model.ChatCompletionRequest{
		Model: "ep-20241003011435-spbz4",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("请根据以下问题生成回答："),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(question),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return "", err
	}
	answer := *resp.Choices[0].Message.Content.StringValue
	return answer, nil
}
