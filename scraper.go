package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetCookies(w http.ResponseWriter, req *http.Request) {
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

	service, err := selenium.NewChromeDriverService("./chromedriver", 8080)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Path: "./chrome-linux64/chrome",
		Args: []string{
			// "window-size=1920x1080",
			// "--no-sandbox",
			// "--disable-dev-shm-usage",
			// "disable-gpu",
			"--headless",
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

	authFlow(driver, creds.Username, creds.Password)

	// allCookies, err := driver.GetCookies()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fileAllCookies, err := os.Create("allCookies.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer fileAllCookies.Close()
	// encoder := json.NewEncoder(fileAllCookies)
	// encoder.SetIndent("", "  ")
	// err = encoder.Encode(allCookies)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Успешно сохранили Cookies в allCookies.json")

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

func main() {
	http.HandleFunc("/getcookies", GetCookies)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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

// func pageScreenshot(driver selenium.WebDriver, fileName string) {
// 	byteImg, err := driver.Screenshot()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	f, err := os.Create("./screenshots/" + fileName + ".png")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	f.Write(byteImg)
// 	f.Close()
// 	fmt.Printf("Screen save: %s", fileName)
// }

func authFlow(driver selenium.WebDriver, username, password string) {
	err := driver.Get("https://www.threads.net/login/")
	if err != nil {
		log.Fatal("Error:", err)
	}

	driver.SetPageLoadTimeout(100 * time.Second)

	// pageScreenshot(driver, "screen1")
	time.Sleep(4 * time.Second)
	// pageScreenshot(driver, "screen2")
	time.Sleep(2 * time.Second)

	acceptAllCookies(driver)

	time.Sleep(2 * time.Second)
	// pageScreenshot(driver, "screen3")

	continueWithInstagram(driver)

	time.Sleep(2 * time.Second)
	// pageScreenshot(driver, "screen4")

	fillCredsAndLogin(driver, username, password)

	time.Sleep(20 * time.Second)
	// pageScreenshot(driver, "screen7")

	//get cookies
	// getAllCookies(driver)
}

func cryptoRandom(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func acceptAllCookies(driver selenium.WebDriver) {
	var elemCookieAccept selenium.WebElement
	//find with waiting
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Allow all cookies')]]")
			if err != nil {
				panic(fmt.Errorf("не удалось найти кнопку 'Разрешить все cookie': %v", err))
			}
			// elemCookieAccept = foundElem
		}
		elemCookieAccept = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//scroll to element
	driver.ExecuteScript("arguments[0].scrollIntoView({block: 'center'});", []interface{}{elemCookieAccept})

	//click
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemCookieAccept})
	if err != nil {
		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Разрешить все cookie': %v", err))
	}
	fmt.Println("Успешно нажали на 'Разрешить все cookie'")
}

func continueWithInstagram(driver selenium.WebDriver) {
	//find with waiting
	var elemContinueWithInstagram selenium.WebElement
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := wd.FindElement(selenium.ByXPATH, "//a[.//span[contains(text(), 'Продолжить с аккаунтом Instagram')]]")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByXPATH, "//a[.//span[contains(text(), 'Continue with Instagram')]]")
			if err != nil {
				panic(fmt.Errorf("не удалось найти кнопку 'Продолжить с аккаунтом Instagram': %v", err))
			}
			// elemContinueWithInstagram = foundElem
		}
		elemContinueWithInstagram = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//scroll to element
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	driver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{elemContinueWithInstagram})

	//click
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemContinueWithInstagram})
	if err != nil {
		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Продолжить с аккаунтом Instagram': %v", err))
	}
	fmt.Println("Успешно нажали на 'Продолжить с аккаунтом Instagram'")
}

func fillCredsAndLogin(driver selenium.WebDriver, username, password string) {
	//find with waiting
	var elemUsername selenium.WebElement
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Имя пользователя, номер телефона или электронный адрес']")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Username, phone or email']")
			if err != nil {
				panic(fmt.Errorf("не удалось найти кнопку 'Имя пользователя, номер телефона или электронный адрес': %v", err))
			}
			// elemUsername = foundElem
		}
		elemUsername = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//fill input
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	err = elemUsername.SendKeys(username)
	if err != nil {
		panic(fmt.Errorf("не удалось ввести 'username': %v", err))
	}

	time.Sleep(1 * time.Second)
	// pageScreenshot(driver, "screen5")

	//find with waiting
	var elemPassword selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Пароль']")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Password']")
			if err != nil {
				panic(fmt.Errorf("не удалось найти кнопку 'Пароль': %v", err))
			}
			// elemPassword = foundElem
		}
		elemPassword = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//fill input
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	err = elemPassword.SendKeys(password)
	if err != nil {
		panic(fmt.Errorf("не удалось ввести 'password': %v", err))
	}

	time.Sleep(1 * time.Second)
	// pageScreenshot(driver, "screen6")

	//find with waiting
	var elemSignInButton selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Войти')]]")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Log in')]]")
			if err != nil {
				panic(fmt.Errorf("не удалось найти кнопку 'Войти': %v", err))
			}
			// elemSignInButton = foundElem
		}
		elemSignInButton = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//click
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemSignInButton})
	if err != nil {
		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Войти': %v", err))
	}
	fmt.Println("Успешно нажали на 'Вход'")
}

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
