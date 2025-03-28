package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"

	"github.com/Extaleus/selenium-project/common"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LikesCountRequest struct {
	LikesNeeded int `json:"likesneeded"`
}

func main() {
	randNum, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	userDataDir := filepath.Join(os.TempDir(), fmt.Sprintf("chrome-data-%d", randNum))
	defer os.RemoveAll(userDataDir)

	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Path: "./chrome-linux64/chrome",
		Prefs: map[string]interface{}{
			"intl.accept_languages": "ru,ru-RU",
		},
		Args: []string{
			"--lang=ru",
			"--accept-lang=ru-RU",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--user-data-dir=" + userDataDir,
			"disable-gpu",
			"--headless=new",
		}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer driver.Quit()

	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
	}

	http.HandleFunc("/getcookies", func(w http.ResponseWriter, r *http.Request) {
		GetCookies(w, r, driver)
	})

	http.HandleFunc("/getposts", func(w http.ResponseWriter, r *http.Request) {
		GetPosts(w, r, driver)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func GetCookies(w http.ResponseWriter, req *http.Request, driver selenium.WebDriver) {
	w.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds AuthRequest
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	common.AuthFlow(driver, creds.Username, creds.Password)

	allCookies, err := driver.GetCookies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to get cookies",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cookies": allCookies,
	})
}

func GetPosts(w http.ResponseWriter, req *http.Request, driver selenium.WebDriver) {
	w.Header().Set("Content-Type", "application/json")

	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input LikesCountRequest
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	results := common.CollectPosts(driver, input.LikesNeeded)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

//
//
//
//
//

//logout
// logOut(driver)

// time.Sleep(3 * time.Second)
// pageScreenshot(driver, "screen11")

//clean up cookies
// cleanUpAllCookies(driver)

//
//
//
//
//

// func getAllCookies(driver selenium.WebDriver) {
// 	allCookies, err := driver.GetCookies()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fileAllCookies, err := os.Create("allCookies.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer fileAllCookies.Close()
// 	encoder := json.NewEncoder(fileAllCookies)
// 	encoder.SetIndent("", "  ")
// 	err = encoder.Encode(allCookies)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Успешно сохранили Cookies в allCookies.json")
// }

// func cleanUpAllCookies(driver selenium.WebDriver) {
// 	err := driver.DeleteAllCookies()
// 	if err != nil {
// 		log.Printf("Не удалось удалить все cookies: %v", err)
// 	}
// }

// func logOut(driver selenium.WebDriver) {
// 	//find with waiting
// 	var elemMenu selenium.WebElement
// 	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
// 		foundElem, err := driver.FindElement(selenium.ByXPATH,
// 			"//div[@role='button']"+
// 				"[.//*[local-name()='svg' and @aria-label='Ещё']]"+
// 				"[.//*[local-name()='title' and text()='Ещё']]"+
// 				"[.//*[local-name()='rect']]"+
// 				"[.//*[local-name()='rect']]")
// 		if err != nil {
// 			panic(fmt.Errorf("не удалось найти кнопку 'Дополнительное меню': %v", err))
// 		}
// 		elemMenu = foundElem
// 		visible, err := foundElem.IsDisplayed()
// 		return visible, err
// 	}, 10*time.Second)
// 	if err != nil {
// 		panic(fmt.Errorf("не удалось найти элемент: %v", err))
// 	}

// 	//click
// 	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
// 	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemMenu})
// 	if err != nil {
// 		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Дополнительное меню': %v", err))
// 	}
// 	fmt.Println("Успешно нажали на 'Дополнительное меню'")

// 	time.Sleep(2 * time.Second)
// 	// pageScreenshot(driver, "screen8")

// 	//tab to exit
// 	for i := 0; i < 5; i++ {
// 		driver.KeyDown(selenium.TabKey)
// 		randDelay := cryptoRandom(50, 100)
// 		time.Sleep(time.Duration(randDelay) * time.Millisecond)
// 		driver.KeyUp(selenium.TabKey)
// 		randDelay = cryptoRandom(200, 400)
// 		time.Sleep(time.Duration(randDelay) * time.Millisecond)
// 	}

// 	// pageScreenshot(driver, "screen9")

// 	//exit
// 	driver.KeyDown(selenium.EnterKey)
// 	randDelay := cryptoRandom(50, 100)
// 	time.Sleep(time.Duration(randDelay) * time.Millisecond)
// 	driver.KeyUp(selenium.EnterKey)
// 	fmt.Println("Успешно нажали на 'Выход'")
// }
