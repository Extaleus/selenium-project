package common

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tebeka/selenium"
)

func AuthFlow(c *gin.Context, driver selenium.WebDriver, username, password string) {
	err := driver.Get("https://www.threads.net/login/")
	if err != nil {
		log.Fatal("Error:", err)
	}

	driver.SetPageLoadTimeout(30 * time.Second)

	err = WaitForPageLoad(driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Page not load"})
		return
	}

	time.Sleep(time.Duration(CryptoRandom(1000, 5000)) * time.Millisecond)

	// // PageScreenshot(driver, "1")
	// time.Sleep(10 * time.Second)
	// // PageScreenshot(driver, "2")
	// time.Sleep(5 * time.Second)

	_, errRus := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	_, errEng := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Allow all cookies')]]")
	if errRus == nil || errEng == nil {
		AcceptAllCookies(driver)
	}

	// AcceptAllCookies(driver)

	time.Sleep(time.Duration(CryptoRandom(500, 2000)) * time.Millisecond)

	// time.Sleep(5 * time.Second)
	// PageScreenshot(driver, "5")

	continueWithInstagram(c, driver)

	driver.SetPageLoadTimeout(30 * time.Second)

	err = WaitForPageLoad(driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Page not load"})
		return
	}

	time.Sleep(time.Duration(CryptoRandom(1000, 5000)) * time.Millisecond)

	_, errRus = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	_, errEng = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Allow all cookies')]]")
	if errRus == nil || errEng == nil {
		AcceptAllCookies(driver)
	}

	time.Sleep(time.Duration(CryptoRandom(1000, 2000)) * time.Millisecond)

	// time.Sleep(5 * time.Second)
	// PageScreenshot(driver, "6")

	// foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	// if err == nil {
	// 	fmt.Printf(foundElem.Text())
	// 	AcceptAllCookies(driver)
	// }

	// time.Sleep(5 * time.Second)
	// PageScreenshot(driver, "6")

	fillCredsAndLogin(c, driver, username, password)

	// time.Sleep(5 * time.Second)
	// PageScreenshot(driver, "6.1")

	driver.SetPageLoadTimeout(30 * time.Second)

	err = WaitForPageLoad(driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Page not load"})
		return
	}

	time.Sleep(time.Duration(CryptoRandom(1000, 5000)) * time.Millisecond)

	_, errRus = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	_, errEng = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Allow all cookies')]]")
	if errRus == nil || errEng == nil {
		AcceptAllCookies(driver)
	}

	// foundElem, err = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Разрешить все cookie')]]")
	// if err == nil {
	// 	fmt.Printf(foundElem.Text())
	// 	AcceptAllCookies(driver)
	// }

	// time.Sleep(10 * time.Second)
	// PageScreenshot(driver, "10")

	//get cookies
	// getAllCookies(driver)
}

func continueWithInstagram(c *gin.Context, driver selenium.WebDriver) {
	//find with waiting
	var elemContinueWithInstagram selenium.WebElement
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := wd.FindElement(selenium.ByXPATH, "//a[.//span[contains(text(), 'Продолжить с аккаунтом Instagram')]]")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByXPATH, "//a[.//span[contains(text(), 'Continue with Instagram')]]")
			if err != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return false, err
			}
		}
		elemContinueWithInstagram = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
		// panic(fmt.Errorf("не удалось найти элемент: %v", err))
	}

	//scroll to element
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
	driver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{elemContinueWithInstagram})

	//click
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemContinueWithInstagram})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
		// panic(fmt.Errorf("не удалось кликнуть по кнопке 'Продолжить с аккаунтом Instagram': %v", err))
	}

	fmt.Println("Успешно нажали на 'Продолжить с аккаунтом Instagram'")
}

func fillCredsAndLogin(c *gin.Context, driver selenium.WebDriver, username, password string) {
	var elemUsername selenium.WebElement
	err := driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Имя пользователя, номер телефона или электронный адрес']")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Username, phone or email']")
			if err != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				// panic(fmt.Errorf("не удалось найти кнопку 'Имя пользователя, номер телефона или электронный адрес': %v", err))
				return false, err
			}
		}
		elemUsername = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// panic(fmt.Errorf("не удалось найти элемент: %v", err))
		return
	}

	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)

	err = elemUsername.SendKeys(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// panic(fmt.Errorf("не удалось ввести 'username': %v", err))
		return
	}

	// time.Sleep(1 * time.Second)
	// PageScreenshot(driver, "8")

	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)

	var elemPassword selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Пароль']")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByCSSSelector, "input[placeholder='Password']")
			if err != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				// panic(fmt.Errorf("не удалось найти кнопку 'Пароль': %v", err))
				return false, err
			}
		}
		elemPassword = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// panic(fmt.Errorf("не удалось найти элемент: %v", err))
		return
	}

	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)

	err = elemPassword.SendKeys(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// panic(fmt.Errorf("не удалось ввести 'password': %v", err))
		return
	}

	// time.Sleep(1 * time.Second)
	// PageScreenshot(driver, "9")
	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)

	var elemSignInButton selenium.WebElement
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		foundElem, err := driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Войти')]]")
		if err != nil {
			foundElem, err = driver.FindElement(selenium.ByXPATH, "//div[@role='button' and .//div[contains(text(), 'Log in')]]")
			if err != nil {
				// c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				// panic(fmt.Errorf("не удалось найти кнопку 'Войти': %v", err))
				return false, err
			}
		}
		elemSignInButton = foundElem
		visible, err := foundElem.IsDisplayed()
		return visible, err
	}, 10*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// panic(fmt.Errorf("не удалось найти элемент: %v", err))
		return
	}

	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)

	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{elemSignInButton})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		// panic(fmt.Errorf("не удалось кликнуть по кнопке 'Войти': %v", err))
		return
	}

	fmt.Println("Успешно нажали на 'Вход'")
}
