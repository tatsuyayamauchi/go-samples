package main

import (
	"fmt"
	"log"
	"os"

	"github.com/slack-go/slack"
)

// WEBHOOK_URL=//hooks.slack.com/services/xxxx/xxxx/xxxx go run main.go
func main() {
	if err := Run(os.Getenv("WEBHOOK_URL")); err != nil {
		log.Fatal(err)
	}
}

func Run(webhookURL string) error {
	if len(webhookURL) == 0 {
		return fmt.Errorf("empty webhook URL")
	}

	/* Block kit builder上でのコード
	{
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": "Blocks text :tada:",
					"emoji": true
				}
			},
			{
				"type": "rich_text",
				"elements": [
					{
						"type": "rich_text_list",
						"style": "bullet",
						"indent": 0,
						"elements": [
							{
								"type": "rich_text_section",
								"elements": [
									{
										"type": "text",
										"text": "field1: "
									},
									{
										"type": "text",
										"text": "11111",
										"style": {
											"code": true
										}
									}
								]
							},
							{
								"type": "rich_text_section",
								"elements": [
									{
										"type": "text",
										"text": "field2: 22222"
									}
								]
							},
							{
								"type": "rich_text_section",
								"elements": [
									{
										"type": "text",
										"text": "field3: 11111",
										"style": {
											"bold": true,
											"italic": true,
											"strike": true,
											"code": true
										}
									}
								]
							}
						]
					}
				]
			}
		]
	}
	*/

	bm := &slack.WebhookMessage{
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewHeaderBlock(&slack.TextBlockObject{
					Type:  slack.PlainTextType,
					Text:  "Blocks text :tada:",
					Emoji: true,
				}),
				slack.NewRichTextBlock("message",
					slack.NewRichTextList(slack.RTEListBullet, 0,
						slack.NewRichTextSection(
							slack.NewRichTextSectionTextElement("field1: ", nil),
							slack.NewRichTextSectionTextElement("11111", &slack.RichTextSectionTextStyle{Code: true}),
						),
						slack.NewRichTextSection(
							slack.NewRichTextSectionTextElement("field2: 22222", &slack.RichTextSectionTextStyle{Bold: true}),
						),
						slack.NewRichTextSection(
							slack.NewRichTextSectionTextElement("field3: 33333", &slack.RichTextSectionTextStyle{Bold: true, Italic: true, Strike: true, Code: true}),
						),
					),
				),
			},
		},
	}

	am := &slack.WebhookMessage{
		Text: "*Attachments text* :tada:",
		Attachments: []slack.Attachment{
			{
				Text:       "・field1: `11111`",
				MarkdownIn: []string{"text"},
			},
			{
				Text:       "・ *field2: 22222*",
				MarkdownIn: []string{"text"},
			},
			{
				Text:       "・ ~`field3: 33333`~",
				MarkdownIn: []string{"text"},
			},
		},
	}

	if err := slack.PostWebhook(webhookURL, bm); err != nil {
		return err
	}
	if err := slack.PostWebhook(webhookURL, am); err != nil {
		return err
	}

	return nil
}
