package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"unicode"

	"github.com/tebeka/selenium"
)

func CryptoRandom(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

// func PageScreenshot(driver selenium.WebDriver, fileName string) {
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
// 	fmt.Printf("Screen save: %s\n", fileName)
// }

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

type Post struct {
	Username    string `json:"username"`
	Description string `json:"description"`
	Likes       string `json:"likes"`
}

type Result struct {
	// MainPost Post   `json:"main_post"`
	Answers []Post `json:"answers"`
}

func ParsePostEntities(driver selenium.WebDriver) []byte {
	// fullHTML, err := driver.PageSource()
	// if err != nil {
	// 	log.Fatal("Не удалось получить HTML страницы:", err)
	// }

	// // 1. Переключаемся на контекст <head>
	// headElem, err := driver.FindElement(selenium.ByTagName, "head")
	// if err != nil {
	// 	log.Fatal("Не удалось найти <head>:", err)
	// }

	// headHtml, err := headElem.GetAttribute("innerHTML")
	// if err != nil {
	// 	fmt.Printf("ошибка получения HTML элемента: %v", err)
	// }

	// fmt.Println(headHtml)

	// Поиск всех постов

	// parent, err := driver.FindElement(selenium.ByXPATH, "//div[@aria-label='Содержимое столбца']")
	// if err != nil {
	// 	log.Printf("Ошибка поиска постов1: %v", err)
	// 	// continue
	// }

	// postEntity, err := driver.FindElements(selenium.ByXPATH, "(//div[@data-pressable-container='true'])[1]/*/*/*/div[1]")
	// postEntity, err := driver.FindElements(selenium.ByXPATH, ".(//div[@data-pressable-container='true'])[1]/*/*/*/div[1]")
	// postEntity, err := driver.FindElements(selenium.ByXPATH, "//div[@data-pressable-container='true']/*/*/*/div")
	postEntity, err := driver.FindElements(selenium.ByXPATH, "//div[@data-pressable-container='true']/*/*/*/div")
	if err != nil {
		log.Printf("Ошибка поиска постов2: %v", err)
		// continue
	}

	var data []string
	foundNonEmpty := false

	for i, pe := range postEntity {
		val, err := pe.Text()
		if err != nil {
			fmt.Printf("Ошибка получения текста для элемента %d: %v\n", i, err)
			continue
		}

		if val != "" || foundNonEmpty {
			data = append(data, val)
			foundNonEmpty = true
		}
	}

	result := parseData(data)
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	fmt.Println(string(jsonData))

	return jsonData

	// for i, pe := range postEntity {
	// 	val, err := pe.Text()
	// 	if err != nil {
	// 		fmt.Printf("postEntityChild1: HEUTA")
	// 	}
	// 	if val != "" {
	// 		fmt.Printf("postEntityChild %d: %v\n", i, val)
	// 		fmt.Println()
	// 	}
	// }

	// postEntityChild1
	// if err != nil {
	// 	log.Printf("Ошибка поиска постов: %v", err)
	// 	// continue
	// }
	// fmt.Printf("postEntityChild1: %v\n", postEntityChild1)

	// postEntity[0]

	// // 2. Ищем meta[og:title] именно в <head>
	// titleElems, err := headElem.FindElements(selenium.ByCSSSelector, "meta[property='og:title']")
	// if err != nil {
	// 	log.Fatal("Meta og:title не найден:", err)
	// }

	// // 3. Обрабатываем результаты
	// for _, metaTag := range titleElems {
	// 	title, err := metaTag.GetAttribute("content")
	// 	if err != nil {
	// 		log.Println("Не удалось получить content:", err)
	// 		continue
	// 	}

	// 	// title = strings.ReplaceAll(title, "auf Threads)", "")
	// 	fmt.Println("og:title:", strings.TrimSpace(title))
	// }

	// // 2. Получаем атрибут 'content'
	// title, err := titleElem.GetAttribute("content")
	// if err != nil {
	// 	log.Fatal("Не удалось получить атрибут:", err)
	// }
	// title = strings.ReplaceAll(title, "auf Threads)", "")
	// fmt.Println("Значение атрибута content:", title)

	// descriptionElem, err := headElem.FindElement(selenium.ByCSSSelector, "meta[property='og:description']")
	// if err != nil {
	// 	log.Fatal("Элемент не найден:", err)
	// }

	// // 2. Получаем атрибут 'content'
	// description, err := descriptionElem.GetAttribute("content")
	// if err != nil {
	// 	log.Fatal("Не удалось получить атрибут:", err)
	// }

	// fmt.Println("Значение атрибута content:", description)

	// scripts, err := driver.FindElements(selenium.ByCSSSelector, `script[type="application/json"]`)
	// if err != nil {
	// 	log.Fatal("Не удалось найти скрипты:", err)
	// }

	// // 2. Фильтруем скрипты, содержащие нужный паттерн
	// for _, script := range scripts {
	// 	content, err := script.GetAttribute("innerHTML")
	// 	if err != nil {
	// 		continue
	// 	}

	// 	// Удаляем лишние пробелы и переносы
	// 	content = strings.TrimSpace(content)

	// 	// Проверяем, начинается ли содержимое с требуемого JSON
	// 	if strings.Contains(content, `"require": [["ScheduledServe"`) {
	// 		fmt.Printf("Найден подходящий скрипт:\n%s\n", content)

	// 		// Дополнительно: парсим JSON
	// 		var data map[string]interface{}
	// 		if err := json.Unmarshal([]byte(content), &data); err == nil {
	// 			fmt.Println("Распарсенный require:", data["require"])
	// 		}
	// 	}
	// }

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

func parseData(data []string) Result {
	var result Result
	var i int

	// // Парсинг главного поста
	// mainPostEnd := findMainPostEnd(data)
	// if mainPostEnd > 0 {
	// 	result.MainPost = parseMainPost(data[:mainPostEnd])
	// 	i = mainPostEnd + 1
	// }

	// // Парсинг ответов
	// for i < len(data) {
	// 	// Пропускаем пустые строки
	// 	if data[i] == "" {
	// 		i++
	// 		continue
	// 	}

	// 	// Находим начало поста (username)
	// 	username := data[i]
	// 	i++

	// 	// Пропускаем пустые строки между username и description
	// 	for i < len(data) && data[i] == "" {
	// 		i++
	// 	}

	// 	// Получаем description
	// 	var description string
	// 	if i < len(data) {
	// 		description = data[i]
	// 		i++
	// 	}

	// 	// Пропускаем пустые строки между description и likes
	// 	for i < len(data) && data[i] == "" {
	// 		i++
	// 	}

	// 	// Получаем likes
	// 	var likes string
	// 	if i < len(data) && !strings.Contains(data[i], "Смотреть действия") {
	// 		likes = data[i]
	// 		i++
	// 	}

	// 	// Создаем пост только если есть все необходимые данные
	// 	if username != "" && description != "" {
	// 		post := Post{
	// 			Username:    strings.Fields(username)[0], // Берем только первую часть (без времени)
	// 			Description: strings.TrimSpace(description),
	// 			Likes:       strings.TrimSpace(likes),
	// 		}
	// 		result.Answers = append(result.Answers, post)
	// 	}

	// 	// Пропускаем пустые строки между постами
	// 	for i < len(data) && data[i] == "" {
	// 		i++
	// 	}
	// }

	// Парсинг ответов
	for i < len(data) {
		// Пропускаем пустые строки
		if data[i] == "" {
			i++
			continue
		}

		// 1. Получаем username (первая непустая строка)
		username := strings.Fields(data[i])[0]
		i++

		// Пропускаем пустые строки
		for i < len(data) && data[i] == "" {
			i++
		}

		// 2. Получаем description (следующая непустая строка)
		if i >= len(data) {
			break
		}
		description := data[i]
		i++

		// Пропускаем пустые строки
		for i < len(data) && data[i] == "" {
			i++
		}

		// 3. Проверяем, есть ли likes (должны быть цифры или \n)
		var likes string
		if i < len(data) && isLikelyLikes(data[i]) {
			likes = data[i]
			i++
		} else {
			likes = "0" // Если лайков нет, ставим 0
		}

		// Добавляем пост
		post := Post{
			Username:    username,
			Description: strings.TrimSpace(description),
			Likes:       strings.TrimSpace(likes),
		}
		result.Answers = append(result.Answers, post)

		// Пропускаем пустые строки между постами
		for i < len(data) && data[i] == "" {
			i++
		}
	}

	return result
}

// func findMainPostEnd(data []string) int {
// 	fmt.Printf("\n\n\n\ndata: %v\n\n\n\n\n\n", data)

// 	for i, s := range data {
// 		// if strings.Contains(s, "Смотреть действия") {
// 		if strings.HasSuffix(s, "Share") {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func parseMainPost(postData []string) Post {
// 	var post Post

// 	newPostData := []string{}
// 	for _, row := range postData {
// 		if row != "" {
// 			newPostData = append(newPostData, row)
// 		}
// 	}

// 	if len(newPostData) > 0 {
// 		parts := strings.Fields(newPostData[0])
// 		if len(parts) > 0 {
// 			post.Username = parts[0]
// 		}
// 	}

// 	if len(newPostData) > 1 {
// 		post.Description = strings.TrimSpace(newPostData[1])
// 	}

// 	if len(newPostData) > 2 {
// 		post.Likes = strings.TrimSpace(newPostData[2])
// 	}

// 	return post
// }

// Функция проверяет, является ли строка вероятными лайками (содержит цифры)
func isLikelyLikes(s string) bool {
	if s == "" {
		return false
	}

	// Проверяем, содержит ли строка хотя бы одну цифру
	hasDigit := false
	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if r == '\n' {
			continue // переносы строки допускаем
		}
		if !unicode.IsDigit(r) && r != '\n' {
			return false // если есть не-цифра и не перенос строки
		}
	}

	return hasDigit
}
