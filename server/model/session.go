package model

import (
	"fmt"
	"time"
	"encoding/json"
	"appengine/memcache"
	. "server/lib"
)

// セッションを開始する。サーバの memcache 内にセッションIDとユーザキーの対応を保存している。
// セッションは最後のページアクセスから24時間の有効で、24時間経過したものは cron で定期的に削除される。
// 関数はセッションIDを返す。
func (this *Model) StartSession(userKey string) string {
	sessionId := ""
	for i := 0; i < 4; i++ {
		sessionId = fmt.Sprintf("%s%s", sessionId, GetRandomizedString())
	}
	expire := time.Now().Add(time.Hour * 24)
	
	data := make(map  [string]string, 2)
	data["u"] = userKey
	data["e"] = expire.String()
	
	encodedData, err := json.Marshal(data)
	item := &memcache.Item {
		Key: sessionId,
		Value: encodedData,
	}
	err = memcache.Set(this.c, item)
	Check(this.c, err)
	
	return sessionId
}

// memcache から指定されたセッションIDのセッション情報を削除する
func (this *Model) RemoveSession(sessionId string) {
	err := memcache.Delete(this.c, sessionId)
	Check(this.c, err)
}

// memcache から指定されたセッションIDに関連づいているユーザキーを返す
// 該当するユーザーキーが見つからなかったら空文字を返す
func (this *Model) GetUserKeyFromSession(sessionId string) string {
	item, err := memcache.Get(this.c, sessionId)
	Check(this.c, err)
	if err != nil {
		this.c.Warningf("セッションIDに関連付けられたユーザIDが存在しませんでした: sessionId:%s", sessionId)
		return ""
	}
	data := make(map[string]string, 2)
	err = json.Unmarshal(item.Value, &data)
	Check(this.c, err)
	return data["u"]
}

// セッションのクリア
func (this *Model) ResetSession() {
	err := memcache.Flush(this.c)
	Check(this.c, err)
}