package common

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

func AuthFlow(driver selenium.WebDriver, username, password string) {
	err := driver.Get("https://www.threads.net/login/")
	if err != nil {
		log.Fatal("Error:", err)
	}

	driver.SetPageLoadTimeout(100 * time.Second)

	PageScreenshot(driver, "1")
	time.Sleep(10 * time.Second)
	PageScreenshot(driver, "2")
	time.Sleep(5 * time.Second)

	acceptAllCookies(driver)

	time.Sleep(5 * time.Second)
	PageScreenshot(driver, "5")

	continueWithInstagram(driver)

	time.Sleep(5 * time.Second)
	PageScreenshot(driver, "6")

	foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	if err == nil {
		fmt.Printf(foundElem.Text())
		acceptAllCookies(driver)
	}

	time.Sleep(5 * time.Second)
	PageScreenshot(driver, "6")

	fillCredsAndLogin(driver, username, password)

	time.Sleep(5 * time.Second)
	PageScreenshot(driver, "6.1")

	foundElem, err = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	if err == nil {
		fmt.Printf(foundElem.Text())
		acceptAllCookies(driver)
	}

	time.Sleep(10 * time.Second)
	PageScreenshot(driver, "10")

	//get cookies
	// getAllCookies(driver)
}

func acceptAllCookies(driver selenium.WebDriver) {
	var elemCookieAccept selenium.WebElement
	//find with waiting
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Allow all cookies')]]")
			if err != nil {
				// return
				fmt.Printf("не удалось найти кнопку 'Разрешить все cookie': %v", err)
			}
			// elemCookieAccept = foundElem
		}
		elemCookieAccept = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		fmt.Printf("не удалось найти элемент: %v", err)
	}

	time.Sleep(2 * time.Second)
	PageScreenshot(driver, "3")
	// scroll to element

	if err == nil {
		driver.ExecuteScript("arguments[0].scrollIntoView({block: 'center'});", []interface{}{elemCookieAccept})

		//click
		time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
		_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemCookieAccept})
		if err != nil {
			fmt.Printf("не удалось кликнуть по кнопке 'Разрешить все cookie': %v", err)
		}
	}

	time.Sleep(2 * time.Second)
	PageScreenshot(driver, "4")

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
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
	driver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{elemContinueWithInstagram})

	//click
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
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
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
	err = elemUsername.SendKeys(username)
	if err != nil {
		panic(fmt.Errorf("не удалось ввести 'username': %v", err))
	}

	time.Sleep(1 * time.Second)
	PageScreenshot(driver, "8")

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
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
	err = elemPassword.SendKeys(password)
	if err != nil {
		panic(fmt.Errorf("не удалось ввести 'password': %v", err))
	}

	time.Sleep(1 * time.Second)
	PageScreenshot(driver, "9")

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
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemSignInButton})
	if err != nil {
		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Войти': %v", err))
	}
	fmt.Println("Успешно нажали на 'Вход'")
}
