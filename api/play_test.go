package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/slotopol/server/api"
	"github.com/slotopol/server/cmd"
	cfg "github.com/slotopol/server/config"

	"github.com/gin-gonic/gin"
)

const email, secret = "player@example.org", "iVI05M"

func ping(t *testing.T, r *gin.Engine) {
	var req = httptest.NewRequest("GET", "/ping", nil)
	var w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Error("ping is not ok")
	}
	var resp = w.Result()
	if !strings.HasPrefix(resp.Header.Get("Server"), "slotopol/") {
		t.Error("alien server")
	}
}

func post(t *testing.T, r *gin.Engine, path string, token string, arg any) (ret gin.H) {
	var err error
	var b []byte

	if b, err = json.Marshal(arg); err != nil {
		t.Error(err)
	}
	var req = httptest.NewRequest("POST", path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	var w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusNoContent {
		return
	}
	if w.Code != http.StatusOK {
		t.Errorf("'%s' is not ok, status %d", path, w.Code)
	}
	var resp = w.Result()
	if b, err = io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	}
	if err = json.Unmarshal(b, &ret); err != nil {
		t.Error(err)
	}
	if w.Code != http.StatusOK {
		t.Logf("message: %v, code: %v", ret["what"], ret["code"])
	}
	return
}

func TestPlay(t *testing.T) {
	var err error
	var arg, ret gin.H
	var token string
	var gid any

	// Prepare in-memory database
	cfg.CfgPath = "../appdata"
	if err = cmd.Init(); err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	var r = gin.New()
	r.HandleMethodNotAllowed = true
	api.SetupRouter(r)

	// Send "ping" and check that server is our
	ping(t, r)

	// Sign-in player
	arg = gin.H{
		"email":  email,
		"secret": secret,
	}
	ret = post(t, r, "/signin", "", arg)
	token = ret["access"].(string)

	// Join game
	arg = gin.H{
		"cid":   1, // 'virtual' club
		"uid":   3, // player ID
		"alias": "Novomatic / Dolphins Pearl",
	}
	ret = post(t, r, "/game/join", token, arg)
	gid = ret["gid"]

	arg = gin.H{
		"gid": gid,
		"bet": 1,
	}
	post(t, r, "/slot/bet/set", token, arg)
	arg = gin.H{
		"gid": gid,
		"sel": 5,
	}
	post(t, r, "/slot/sel/set", token, arg)

	// Play the game with 50 spins
	for range 50 {
		// make spin
		arg = gin.H{
			"gid": gid,
		}
		ret = post(t, r, "/slot/spin", token, arg)
		t.Logf("[spin] gid: %v, sid: %v, wallet: %v", gid, ret["sid"], ret["wallet"])
		var game = ret["game"].(map[string]any)

		// no any more actions on free spins
		if _, ok := game["fsr"]; ok {
			continue
		}
		if _, ok := game["gain"]; !ok {
			continue
		}
		t.Logf("gain: %v", game["gain"])

		// if there has a win, make double-ups sometime
		if rand.Float64() < 0.3 {
			for {
				arg = gin.H{
					"gid":  gid,
					"mult": 2,
				}
				ret = post(t, r, "/slot/doubleup", token, arg)
				var gain = ret["gain"].(float64)
				t.Logf("[doubleup] gid: %v, id: %v, wallet: %v, gain: %g", gid, ret["id"], ret["wallet"], gain)
				if gain == 0 {
					break
				}
				if rand.Float64() < 0.5 {
					arg = gin.H{
						"gid": gid,
					}
					post(t, r, "/slot/collect", token, arg)
					t.Logf("[collect] gid: %v", gid)
					break
				}
			}
		}
	}

	// Part game
	arg = gin.H{
		"gid": gid,
	}
	post(t, r, "/game/part", token, arg)
}
