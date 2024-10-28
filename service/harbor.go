package service

import (
	"fmt"
	"net/http"
)

func deleteHarborImages() {
	url := "http://example.com/resource/123" // 替换为你的目标URL
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}
