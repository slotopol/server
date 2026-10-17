package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/slotopol/server/api"
	"github.com/slotopol/server/cmd"
	cfg "github.com/slotopol/server/config"

	"github.com/gin-gonic/gin"
)

const (
	// JWT-token signing keys
	testAccessKey  = "skJgM4NsbP3fs4k7vh0gfdkgGl8dJTszdLxZ1sQ9ksFnxbgvw2RsGH8xxddUV479"
	testRefreshKey = "zxK4dUnuq3Lhd1Gzhpr3usI5lAzgvy2t3fmxld2spzz7a5nfv0hsksm9cheyutie"

	// admin credentials
	testAdminEmail  = "admin@example.org"
	testAdminSecret = "0YBoaT"

	// test user credentials
	testUserEmail  = "player@example.org"
	testUserSecret = "iVI05M"
)

// set of bets per line
var betset = []float64{0.1, 0.2, 0.5, 1, 2, 5, 10}

// set of sums to add to wallet
var sumset = []float64{
	50, 50, 50, 50,
	100, 100, 100, 100, 100, 100,
	200, 200, 200, 200,
	250, 250, 250, 250, 250, 250,
	300, 300,
	500, 500, 500, 500, 500, 500, 500, 500,
	600,
	700,
	800,
	1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000,
	1500, 1500,
	2000, 2000,
	5000, 5000,
	10000,
}

var (
	ErrHttpNotOk   = errors.New("http response is not ok")
	ErrAlienServer = errors.New("alien server")
)

func ping(t *testing.T, r *gin.Engine) (err error) {
	var req = httptest.NewRequest("GET", "/ping", nil)
	var w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		err = fmt.Errorf("%w, ping failed, code %d", ErrHttpNotOk, w.Code)
		t.Error(err)
	}
	var resp = w.Result()
	if !strings.HasPrefix(resp.Header.Get("Server"), "slotopol/") {
		err = fmt.Errorf("%w, header: %s", ErrAlienServer, resp.Header.Get("Server"))
		t.Error(err)
	}
	return
}

func post(t *testing.T, r *gin.Engine, path string, token string, arg any) (ret gin.H, err error) {
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
		err = fmt.Errorf("%w, path '%s', code %d", ErrHttpNotOk, path, w.Code)
		t.Error(err)
	}
	var resp = w.Result()
	if !strings.HasPrefix(resp.Header.Get("Server"), "slotopol/") {
		err = fmt.Errorf("%w, header: %s", ErrAlienServer, resp.Header.Get("Server"))
		t.Error(err)
	}
	if b, err = io.ReadAll(resp.Body); err != nil {
		t.Error(err)
		return
	}
	if err = json.Unmarshal(b, &ret); err != nil {
		t.Error(err)
		return
	}
	if w.Code != http.StatusOK {
		t.Logf("message: %v, code: %v", ret["what"], ret["code"])
	}
	return
}

// Simulate a player playing a game with random bets and lines selected,
// and sometimes making double-ups on wins.
func TestPlay(t *testing.T) {
	const cid, uid uint64 = 1, 3
	var err error
	var arg, ret gin.H
	var admtoken, usrtoken string
	var gid uint64
	var wallet, gain float64
	var fsr int

	var exitctx = context.Background()

	// Load yaml-files
	if _, err := cmd.LoadInternalYaml(exitctx); err != nil {
		t.Fatalf("can not load internal yaml files: %s", err.Error())
		return
	}
	cmd.UpdateAlgList()
	cmd.CheckAlgList()

	// Prepare in-memory database
	cfg.CfgPath = "../appdata"
	if err := cmd.InitXORM(); err != nil {
		t.Fatalf("can not initialize xorm: %s", err.Error())
	}

	gin.SetMode(gin.TestMode)
	var r = gin.New()
	r.HandleMethodNotAllowed = true
	api.SetupRouter(r)

	// Setup test JWT-token signing keys
	cfg.Cfg.AccessKey = testAccessKey
	cfg.Cfg.RefreshKey = testRefreshKey

	// Send "ping" and check that server is our
	if err = ping(t, r); err != nil {
		return
	}

	// Sign-in admin in order to top up player account
	arg = gin.H{
		"email":  testAdminEmail,
		"secret": testAdminSecret,
	}
	if ret, err = post(t, r, "/signin", "", arg); err != nil {
		return
	}
	admtoken = ret["access"].(string)
	t.Logf("[sign-in] %s", arg["email"])

	// Sign-in player with predefined credentials
	arg = gin.H{
		"email":  testUserEmail,
		"secret": testUserSecret,
	}
	if ret, err = post(t, r, "/signin", "", arg); err != nil {
		return
	}
	usrtoken = ret["access"].(string)
	t.Logf("[sign-in] %s", arg["email"])

	// Join game
	arg = gin.H{
		"cid":   cid, // 'virtual' club
		"uid":   uid, // player ID
		"alias": "Novomatic / Dolphins Pearl",
	}
	if ret, err = post(t, r, "/game/new", usrtoken, arg); err != nil {
		return
	}
	gid = uint64(ret["gid"].(float64))
	wallet = ret["wallet"].(float64)
	t.Logf("[new] gid: %d, wallet: %.2f", gid, wallet)

	var bet, sel = 1., 5
	arg = gin.H{
		"gid": gid,
		"bet": bet,
	}
	post(t, r, "/slot/bet/set", usrtoken, arg)
	t.Logf("[bet/set] gid: %d, bet: %g", gid, bet)
	arg = gin.H{
		"gid": gid,
		"sel": sel,
	}
	post(t, r, "/slot/sel/set", usrtoken, arg)
	t.Logf("[sel/set] gid: %d, sel: %d", gid, sel)

	// Play the game with 100 spins
	for range 100 {
		// check money at wallet
		if wallet < bet*float64(sel) {
			var sum float64
			for wallet+sum < bet*float64(sel) {
				sum = sumset[rand.N(len(sumset))]
			}
			arg = gin.H{
				"cid": cid,
				"uid": uid,
				"sum": sum,
			}
			if ret, err = post(t, r, "/prop/wallet/add", admtoken, arg); err != nil {
				return
			}
			wallet = ret["wallet"].(float64)
			t.Logf("[wallet/add] gid: %d, wallet: %.2f, sum: %g", gid, wallet, sum)
		}

		// make spin
		arg = gin.H{
			"gid": gid,
		}
		if ret, err = post(t, r, "/slot/spin", usrtoken, arg); err != nil {
			return
		}
		var game = ret["game"].(map[string]any)
		if v, ok := game["gain"]; ok {
			gain = v.(float64)
		} else {
			gain = 0
		}
		if v, ok := game["fsr"]; ok {
			fsr = int(v.(float64))
		} else {
			fsr = 0
		}
		wallet = ret["wallet"].(float64)
		t.Logf("[spin] gid: %d, slid: %d, fsr: %d, wallet: %.2f, gain: %g", gid, uint64(ret["slid"].(float64)), fsr, wallet, gain)

		// no any more actions on free spins
		if fsr > 0 {
			continue
		}

		// if there has a win, make double-ups sometime
		if gain > 0 && rand.Float64() < 0.3 {
			for {
				arg = gin.H{
					"gid":  gid,
					"mult": 2,
					"half": rand.Float64() < 0.25,
				}
				if ret, err = post(t, r, "/slot/doubleup", usrtoken, arg); err != nil {
					return
				}
				var gain = ret["gain"].(float64)
				wallet = ret["wallet"].(float64)
				t.Logf("[doubleup] gid: %d, mlid: %d, wallet: %.2f, gain: %g", gid, uint64(ret["mlid"].(float64)), wallet, gain)
				if gain == 0 {
					break
				}
				if rand.Float64() < 0.5 {
					arg = gin.H{
						"gid": gid,
					}
					post(t, r, "/slot/collect", usrtoken, arg)
					t.Logf("[collect] gid: %d", gid)
					break
				}
			}
		}

		// change bet value sometimes
		if rand.Float64() < 1./25. {
			bet = betset[rand.N(len(betset))]
			arg = gin.H{
				"gid": gid,
				"bet": bet,
			}
			post(t, r, "/slot/bet/set", usrtoken, arg)
			t.Logf("[bet/set] gid: %d, bet: %g", gid, bet)
		}

		// change selected bet lines sometimes
		if rand.Float64() < 1./25. {
			sel = 3 + rand.N(8)
			arg = gin.H{
				"gid": gid,
				"sel": sel,
			}
			post(t, r, "/slot/sel/set", usrtoken, arg)
			t.Logf("[sel/set] gid: %d, sel: %d", gid, sel)
		}
	}

	// Close database connection
	if err := cmd.DoneXORM(); err != nil {
		t.Fatal(err)
	}
}
