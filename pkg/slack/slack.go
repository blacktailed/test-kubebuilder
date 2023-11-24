package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blacktailed/test-kubebuilder.git/pkg/common"
)

type SlackAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

func TestMess(content *common.SlackMsg) {
	webhookURL := content.WebhookURL //os.Getenv("SLACK_WEBHOOK_URL")

	// Webhook URL이 비어있는지 확인합니다.
	if webhookURL == "" {
		fmt.Println("SLACK_WEBHOOK_URL 환경 변수를 설정해야 합니다.")
		return
	}

	// 보낼 attachment를 설정합니다.
	attachment := SlackAttachment{
		Title: "Warning",
		Text:  string(content.Text),
		Color: "#DF0101",
	}

	// Slack으로 보낼 데이터를 준비합니다.
	payload, err := json.Marshal(map[string]interface{}{
		"text":        "경고",
		"username":    "Alert Alarm",
		"attachments": []SlackAttachment{attachment},
		// "channel":     "#test-alert", // 토큰 사용 시 채널 지정 가능
	})
	if err != nil {
		fmt.Printf("JSON 데이터 생성 중 오류 발생: %v\n", err)
		return
	}

	// Slack Webhook URL로 POST 요청을 보냅니다.
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("메시지를 보내는 중 오류 발생: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("메시지가 성공적으로 전송되었습니다.")
}

/**
* slack 라이브러리 이용하는 방법
 */
// func SendSlack(token string) {
// 	api := slack.New(token)
// 	attachment := slack.Attachment{
// 		Pretext: "내가 잘해야!",
// 		Text:    "모두가 편하다 삐빕",
// 	}

// 	channelID, timestamp, err := api.PostMessage(
// 		"test-alert",
// 		slack.MsgOptionText("", false),
// 		slack.MsgOptionAttachments(attachment),
// 		slack.MsgOptionAsUser(false), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
// 	)
// 	if err != nil {
// 		fmt.Printf("%s\n", err)
// 		return
// 	}
// 	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
// }
