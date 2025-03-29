package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sync"

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
	LikesNeeded int `json:"likesneeded"`
}

var (
	sseClients   = make(map[string]chan string) // taskID -> канал
	sseClientsMu sync.Mutex
)

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

	// Измененный GET /getposts с SSE
	r.POST("/getposts", func(c *gin.Context) {
		var input LikesCountRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		taskID := fmt.Sprintf("task_%d", randNum.Int64())

		// Создаем канал для SSE
		sseClientsMu.Lock()
		sseClients[taskID] = make(chan string)
		sseClientsMu.Unlock()

		// Запускаем сбор постов в фоне
		go func() {
			results := common.CollectPosts(driver, input.LikesNeeded)
			jsonData, _ := json.Marshal(results)

			sseClientsMu.Lock()
			if ch, exists := sseClients[taskID]; exists {
				ch <- string(jsonData)
				close(ch)
			}
			sseClientsMu.Unlock()
		}()

		c.JSON(http.StatusOK, gin.H{"task_id": taskID})
	})

	// SSE эндпоинт для n8n
	r.GET("/sse/:taskID", func(c *gin.Context) {
		taskID := c.Param("taskID")

		sseClientsMu.Lock()
		eventChan, exists := sseClients[taskID]
		sseClientsMu.Unlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.Stream(func(w io.Writer) bool {
			select {
			case msg, ok := <-eventChan:
				if !ok {
					return false
				}
				c.SSEvent("message", msg)
				return false // Закрываем соединение после отправки
			case <-c.Writer.CloseNotify():
				return false
			}
		})

		// Очистка
		sseClientsMu.Lock()
		delete(sseClients, taskID)
		sseClientsMu.Unlock()
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
