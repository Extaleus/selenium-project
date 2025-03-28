package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/tebeka/selenium"
)

func CryptoRandom(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
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
	fmt.Printf("Screen save: %s\n", fileName)
}

// func SaveElementHTML(element selenium.WebElement, filename string) error {
// 	html, err := element.GetAttribute("outerHTML")
// 	if err != nil {
// 		return fmt.Errorf("ошибка получения HTML элемента: %v", err)
// 	}

// 	return

// 	err = os.WriteFile(filename, []byte(html), 0644)
// 	if err != nil {
// 		return fmt.Errorf("ошибка записи в файл: %v", err)
// 	}

// 	return nil
// }

// func getLikesFromHTML(wd selenium.WebDriver, html string) (int, error) {
// 	// Создаем и заполняем временный элемент
// 	script := `
//     var temp = document.createElement('div');
//     temp.innerHTML = arguments[0];
//     document.body.appendChild(temp);

//     // Находим элемент с лайками и возвращаем его текст
//     var likeElement = temp.querySelector('div svg[aria-label="Нравится"]').parentNode.nextElementSibling.querySelector('div span');
//     var result = likeElement ? likeElement.innerText.trim() : '';

//     temp.remove();
//     return result;
//     `

// 	likeText, err := wd.ExecuteScript(script, []interface{}{html})
// 	if err != nil {
// 		return 0, fmt.Errorf("ошибка выполнения скрипта: %v", err)
// 	}

// 	text, ok := likeText.(string)
// 	if !ok || text == "" {
// 		return 0, fmt.Errorf("элемент с лайками не найден")
// 	}

// 	likes, err := strconv.Atoi(text)
// 	if err != nil {
// 		return 0, fmt.Errorf("не удалось преобразовать '%s' в число: %v", text, err)
// 	}

// 	return likes, nil
// }

// func FindMoreButton(driver selenium.WebDriver, postText string) {
// 	// Находим пост по тексту (часть текста поста)
// 	// postText := "Сегодня видела живой пример, как родители могут напрочь убить самооценку детей"
// 	// post, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf("//*[contains(text(), '%s')]", postText))
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// moreButton, err := driver.FindElement(selenium.ByXPATH,
// 	// 	fmt.Sprintf("//*[contains(text(), '%s')]/ancestor::div//*[@aria-haspopup='menu']", postText))
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	moreButton, err := driver.FindElement(selenium.ByXPATH,
// 		// fmt.Sprintf("//*[contains(text(), '%s')]/ancestor::div//*[@aria-haspopup='menu']", postText))
// 		fmt.Sprintf("//*[contains(text(), '%s')]/ancestor::*[4]//*[@aria-haspopup='menu']", postText))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// // Кликаем по кнопке
// 	// err = moreButton.Click()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	//click
// 	time.Sleep(time.Duration(CryptoRandom(300, 500)) * time.Millisecond)
// 	_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{moreButton})
// 	if err != nil {
// 		panic(fmt.Errorf("не удалось кликнуть по кнопке 'Ещё': %v", err))
// 	}
// 	fmt.Println("Успешно нажали на 'Ещё'")
// }

// type PostEntity struct {
// 	Username         *string
// 	Description      *string
// 	LikeCount        *string
// 	DirectReplyCount *string
// 	RepostCount      *string
// 	QuoteCount       *string
// }

func ParsePostEntities(driver selenium.WebDriver) {
	// fullHTML, err := driver.PageSource()
	// if err != nil {
	// 	log.Fatal("Не удалось получить HTML страницы:", err)
	// }

	titleElem, err := driver.FindElement(selenium.ByCSSSelector, "meta[property='og:title']")
	if err != nil {
		log.Fatal("Элемент не найден:", err)
	}

	// 2. Получаем атрибут 'content'
	title, err := titleElem.GetAttribute("content")
	if err != nil {
		log.Fatal("Не удалось получить атрибут:", err)
	}
	title = strings.ReplaceAll(title, "auf Threads)", "")
	fmt.Println("Значение атрибута content:", title)

	descriptionElem, err := driver.FindElement(selenium.ByCSSSelector, "meta[property='og:description']")
	if err != nil {
		log.Fatal("Элемент не найден:", err)
	}

	// 2. Получаем атрибут 'content'
	description, err := descriptionElem.GetAttribute("content")
	if err != nil {
		log.Fatal("Не удалось получить атрибут:", err)
	}

	fmt.Println("Значение атрибута content:", description)

	scripts, err := driver.FindElements(selenium.ByCSSSelector, `script[type="application/json"]`)
	if err != nil {
		log.Fatal("Не удалось найти скрипты:", err)
	}

	// 2. Фильтруем скрипты, содержащие нужный паттерн
	for _, script := range scripts {
		content, err := script.GetAttribute("innerHTML")
		if err != nil {
			continue
		}

		// Удаляем лишние пробелы и переносы
		content = strings.TrimSpace(content)

		// Проверяем, начинается ли содержимое с требуемого JSON
		if strings.Contains(content, `"require": [["ScheduledServe"`) {
			fmt.Printf("Найден подходящий скрипт:\n%s\n", content)

			// Дополнительно: парсим JSON
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(content), &data); err == nil {
				fmt.Println("Распарсенный require:", data["require"])
			}
		}
	}

	// targetElement, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf("//span[text()='%s']/ancestor::*[5]", postText))
	// if err != nil {
	// 	log.Fatal("Не удалось найти элемент:", err)
	// }

	// time.Sleep(time.Duration(CryptoRandom(300, 600)) * time.Millisecond)
	// _, err = driver.ExecuteScript("arguments[0].click();", []interface{}{targetElement})
	// if err != nil {
	// 	panic(fmt.Errorf("не удалось кликнуть по кнопке 'Копировать ссылку': %v", err))
	// }

	// time.Sleep(1 * time.Second)
	// PageScreenshot(driver, "clicked")

	// time.Sleep(10 * time.Second)

	// // Поиск всех постов (обновите селектор)
	// postEntities, err := driver.FindElements(selenium.ByCSSSelector, "div[data-pressable-container=\"true\"]")
	// if err != nil {
	// 	log.Printf("Ошибка поиска постов: %v", err)
	// 	// continue
	// }

	//
	//
	//
	//
	//

	// for _, postEntity := range postEntities {
	// var entity PostEntity

	// fmt.Println(postEntity.TagName())

	// // Парсим username
	// if username, err := postEntity.FindElement(selenium.ByCSSSelector, "span[dir='auto']"); err == nil {
	// 	entity.Username, _ = username.Text()
	// }

	// // Парсим описание
	// if desc, err := post.FindElement(selenium.ByCSSSelector, "span[style*='line-height']"); err == nil {
	// 	entity.Description, _ = desc.Text()
	// }

	// // Парсим количество лайков
	// if likes, err := post.FindElement(selenium.ByXPATH, ".//*[contains(@aria-label, 'Нравится')]/following-sibling::span"); err == nil {
	// 	entity.LikeCount, _ = likes.Text()
	// }

	// // Парсим количество репостов
	// if reposts, err := post.FindElement(selenium.ByXPATH, ".//*[contains(@aria-label, 'Репост')]/following-sibling::span"); err == nil {
	// 	entity.ReshareCount, _ = reposts.Text()
	// }

	// // Аналогично для других полей...
	// // entity.DirectReplyCount = ...
	// // entity.RepostCount = ...
	// // entity.QuoteCount = ...

	// postEntriesArr = append(postEntriesArr, entity)
	// }

	//
	//
	//
	//
	//

	// postEntriesArr := []PostEntity{}

	// fmt.Println(len(postEntities))
	// for j := range postEntities {
	// 	postEntity, err := postEntities[j].Text()
	// 	if err == nil {
	// 		postEntityTextArr := strings.Split(postEntity, "\n")
	// 		fmt.Printf("\npost entity array length: %d", len(postEntityTextArr))
	// 		fmt.Println("\npost entity array:")
	// 		for i := range postEntityTextArr {
	// 			fmt.Printf("\t%s\n", postEntityTextArr[i])
	// 		}
	// 		fmt.Println("\tfull entity post:")
	// 		fmt.Printf("postEntityTextArr: %v\n", postEntityTextArr)

	// 		postEntriesArr = append(postEntriesArr, PostEntity{
	// 			Username:         &postEntityTextArr[0],
	// 			RepostCount:      &postEntityTextArr[len(postEntityTextArr)],
	// 			DirectReplyCount: &postEntityTextArr[len(postEntityTextArr)-1],
	// 			LikeCount:        &postEntityTextArr[len(postEntityTextArr)-2],
	// 			Description:      &postEntityTextArr[1],
	// 		})

	// 		fmt.Println(postEntriesArr)
	// 	}
	// }
}
