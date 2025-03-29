package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Extaleus/selenium-project/common"
	"github.com/gin-gonic/gin"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LikesCountRequest struct {
	LikesNeeded int    `json:"likesneeded"`
	CallbackURL string `json:"callback_url"`
}

func main() {
	r := gin.Default()

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

	//clean up cookies
	CleanUpAllCookies(driver)

	// Маршруты
	r.POST("/getcookies", func(c *gin.Context) {
		GetCookies(c, driver)
	})

	// r.GET("/getposts", func(c *gin.Context) {
	// 	GetPosts(c, driver)
	// })

	r.POST("/getposts", func(c *gin.Context) {
		var input LikesCountRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		// Валидация callback URL
		if input.CallbackURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "callback_url is required"})
			return
		}

		// Запускаем сбор постов в фоне
		go func() {
			results := common.CollectPosts(driver, input.LikesNeeded)
			sendCallback(input.CallbackURL, results)
		}()

		c.JSON(http.StatusAccepted, gin.H{
			"status":  "processing",
			"message": "Request accepted. Results will be sent to the callback URL",
		})
	})

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}

func GetCookies(c *gin.Context, driver selenium.WebDriver) {
	// Проверяем метод запроса
	if c.Request.Method != "POST" {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	// Парсим тело запроса
	var creds AuthRequest
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	common.AuthFlow(driver, creds.Username, creds.Password)

	allCookies, err := driver.GetCookies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cookies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cookies": allCookies})
}

// Улучшенная функция отправки callback
func sendCallback(url string, data interface{}) {
	client := &http.Client{Timeout: 10 * time.Second}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal callback data: %v", err)
		return
	}

	// Добавляем логирование для отладки
	log.Printf("Sending callback to: %s", url)
	log.Printf("Payload: %s", string(jsonData))

	// Пытаемся отправить несколько раз в случае ошибки
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to create request (attempt %d): %v", i+1, err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Selenium-Webhook-Service")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Callback attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("Callback successfully delivered")
			resp.Body.Close()
			return
		}

		log.Printf("Callback attempt %d failed with status: %d", i+1, resp.StatusCode)
		resp.Body.Close()
		time.Sleep(2 * time.Second)
	}

	log.Printf("Failed to send callback after %d attempts", maxRetries)
}

// func GetPosts(c *gin.Context, driver selenium.WebDriver) {
// 	if c.Request.Method != "POST" {
// 		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
// 		return
// 	}

// 	var input LikesCountRequest
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
// 		return
// 	}

// 	results := common.CollectPosts(driver, input.LikesNeeded)

// 	c.JSON(http.StatusOK, results)
// }

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

func CleanUpAllCookies(driver selenium.WebDriver) {
	err := driver.DeleteAllCookies()
	if err != nil {
		log.Printf("Не удалось удалить все cookies: %v", err)
	}
}

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
