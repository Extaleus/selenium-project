package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Product struct {
	Name, Price string
}

func main() {
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}

	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Path: "./chrome-linux64/chrome", Args: []string{
		"--headless",
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error:", err)
	}

	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
	}

	err = driver.Get("https://www.threads.net/login/")
	if err != nil {
		log.Fatal("Error:", err)
	}

	driver.SetPageLoadTimeout(100 * time.Second)

	PageScreenshot(driver, "screen1")
	time.Sleep(2 * time.Second)
	PageScreenshot(driver, "screen2")
	time.Sleep(1 * time.Second)

	var elemCookieAccept selenium.WebElement
	//find with waiting
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
		if err != nil {
			panic(fmt.Errorf("не удалось найти кнопку 'Разрешить все cookie': %v", err))
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

	time.Sleep(2 * time.Second)
	PageScreenshot(driver, "screen3")

	//find with waiting
	var elemContinueWithInstagram selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := wd.FindElement(selenium.ByXPATH, "//a[.//span[contains(text(), 'Продолжить с аккаунтом Instagram')]]")
		if err != nil {
			panic(fmt.Errorf("не удалось найти кнопку 'Продолжить с аккаунтом Instagram': %v", err))
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

	time.Sleep(2 * time.Second)
	PageScreenshot(driver, "screen4")

	//find with waiting
	var elemUsername selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Имя пользователя, номер телефона или электронный адрес']")
		if err != nil {
			panic(fmt.Errorf("не удалось найти кнопку 'Имя пользователя, номер телефона или электронный адрес': %v", err))
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
	err = elemUsername.SendKeys("yourcreator32025")
	if err != nil {
		panic(fmt.Errorf("не удалось ввести 'username': %v", err))
	}

	time.Sleep(1 * time.Second)
	PageScreenshot(driver, "screen5")

	//find with waiting
	var elemPassword selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Пароль']")
		if err != nil {
			panic(fmt.Errorf("не удалось найти кнопку 'Пароль': %v", err))
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
	err = elemPassword.SendKeys("Gagarin2$")
	if err != nil {
		panic(fmt.Errorf("не удалось ввести 'password': %v", err))
	}

	time.Sleep(1 * time.Second)
	PageScreenshot(driver, "screen6")

	//find with waiting
	var elemSignInButton selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Войти')]]")
		if err != nil {
			panic(fmt.Errorf("не удалось найти кнопку 'Войти': %v", err))
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

	time.Sleep(20 * time.Second)
	PageScreenshot(driver, "screen7")

	//get cookies
	allCookies, err := driver.GetCookies()
	if err != nil {
		log.Fatal(err)
	}
	fileAllCookies, err := os.Create("allCookies.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fileAllCookies.Close()
	encoder := json.NewEncoder(fileAllCookies)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(allCookies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешно сохранили Cookies в allCookies.json")

	//find with waiting
	var elemMenu selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH,
			"//div[@role='button']"+
				"[.//*[local-name()='svg' and @aria-label='Ещё']]"+
				"[.//*[local-name()='title' and text()='Ещё']]"+
				"[.//*[local-name()='rect']]"+
				"[.//*[local-name()='rect']]")
		if err != nil {
			panic(fmt.Errorf("не удалось найти кнопку 'Дополнительное меню': %v", err))
		}
		elemMenu = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//click
	time.Sleep(time.Duration(cryptoRandom(300, 500)) * time.Millisecond)
	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemMenu})
	if err != nil {
		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Дополнительное меню': %v", err))
	}
	fmt.Println("Успешно нажали на 'Дополнительное меню'")

	time.Sleep(2 * time.Second)
	PageScreenshot(driver, "screen8")

	//tab to exit
	for i := 0; i < 5; i++ {
		driver.KeyDown(selenium.TabKey)
		randDelay := cryptoRandom(50, 100)
		time.Sleep(time.Duration(randDelay) * time.Millisecond)
		driver.KeyUp(selenium.TabKey)
		randDelay = cryptoRandom(200, 400)
		time.Sleep(time.Duration(randDelay) * time.Millisecond)
	}

	PageScreenshot(driver, "screen9")

	//exit
	driver.KeyDown(selenium.EnterKey)
	randDelay := cryptoRandom(50, 100)
	time.Sleep(time.Duration(randDelay) * time.Millisecond)
	driver.KeyUp(selenium.EnterKey)
	fmt.Println("Успешно нажали на 'Выход'")

	PageScreenshot(driver, "screen10")

	time.Sleep(3 * time.Second)
	PageScreenshot(driver, "screen11")

	//clean up cookies
	err = driver.DeleteAllCookies()
	if err != nil {
		log.Printf("Не удалось удалить все cookies: %v", err)
	}
}

func PageScreenshot(driver selenium.WebDriver, fileName string) {
	byteImg, err := driver.Screenshot()
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("./screenshots/" + fileName + ".png")
	if err != nil {
		fmt.Println(err)
	}
	f.Write(byteImg)
	f.Close()
}

func cryptoRandom(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}
