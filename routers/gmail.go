package routers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	messagesUrl string = "https://www.googleapis.com/gmail/v1/users/me/messages"
)

type GmailRouter struct{}

type GmailResult struct {
	MessagesDeleted int32  `json:"messagesDeleted"`
	Success         bool   `json:"success"`
	ResultMessage   string `json:"resultMessage"`
}

type GmailRequest struct {
	Keywords        []string `json:"keywords"`
	Token           string   `json:"token"`
	ClearPromotions *bool    `json:"clearPromotions"`
}

type MessagesList struct {
	Messages      []map[string]string `json:"messages"`
	NextPageToken *string             `json:"nextPageToken"`
}

type MessageResponse struct {
	Id       string                 `json:"id"`
	ThreadId string                 `json:"threadId"`
	LabelIds []string               `json:"labelIds"`
	Snippet  string                 `json:"snippet"`
	Payload  map[string]interface{} `json:"payload"`
}

func (gmail GmailRouter) FlushMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Deletion result")
	nextPageToken := ""

	client := &http.Client{}
	var request GmailRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	deletionCounter := 0

	if err != nil {
		returnError(err)
	} else {
		for do := true; do; do = len(nextPageToken) > 0 {
			messageListUrl := messagesUrl
			if len(nextPageToken) > 0 {
				messageListUrl = messagesUrl + "?pageToken=" + nextPageToken
			}

			// First get the message list
			messageListRequest := buildRequest(request.Token, messageListUrl, "GET")
			messageListResponse, err := client.Do(messageListRequest)
			var messageList MessagesList
			if err != nil {
				returnError(err)
			}

			err = json.NewDecoder(messageListResponse.Body).Decode(&messageList)
			if err != nil {
				returnError(err)
			}

			if messageList.NextPageToken != nil && len(*messageList.NextPageToken) > 0 {
				nextPageToken = *messageList.NextPageToken
			} else {
				nextPageToken = ""
			}

			// Now delete the messages based on keywords
			for _, m := range messageList.Messages {
				if val, ok := m["id"]; ok {
					url := fmt.Sprintf("%s/%s", messagesUrl, val)
					messageRequest := buildRequest(request.Token, url, "GET")
					messageResponse, err := client.Do(messageRequest)
					if err != nil {
						returnError(err)
					}

					var message MessageResponse
					err = json.NewDecoder(messageResponse.Body).Decode(&message)
					if err != nil {
						returnError(err)
					}

					isUnread := contains(message.LabelIds, "UNREAD")

					isPromotion := *request.ClearPromotions && contains(message.LabelIds, "CATEGORY_PROMOTIONS")

					if !isUnread {
						continue
					}

					if badEmail, _ := isBadEmail(message.Payload, request.Keywords); (badEmail != nil && *badEmail) || isPromotion {
						deleteRequest := buildRequest(request.Token, url, "DELETE")
						_, err := client.Do(deleteRequest)
						if err != nil {
							returnError(err)
						}
						deletionCounter++
					}
				}
			}
		}
	}

	result, err := json.Marshal(GmailResult{Success: true, MessagesDeleted: int32(deletionCounter)})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(result))
	}
}

func buildRequest(token string, url string, method string) *http.Request {
	request, err := http.NewRequest(method, url, nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		return nil
	}
	return request
}

func isBadEmail(emailContents map[string]interface{}, keywords []string) (*bool, error) {
	success := false

	headers, ok := emailContents["headers"].([]interface{})
	if !ok {
		return nil, errors.New("error with is bad email getting from")
	}

	from := ""
	for _, header := range headers {
		h := header.(map[string]interface{})
		if name, ok := h["name"]; ok && name == "From" {
			if val, ok := h["value"]; ok {
				from = val.(string)
			}
		}
	}

	if len(from) > 0 {
		from = strings.ToUpper(from)
		for _, key := range keywords {
			key = strings.ToUpper(key)
			if strings.Contains(from, key) {
				success = true
				return &success, nil
			}
		}
	}

	return &success, nil
}

func returnError(err error) {
	result, err := json.Marshal(GmailResult{
		Success:       false,
		ResultMessage: err.Error(),
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(result))
	}
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
