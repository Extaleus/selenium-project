package common

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

func CollectPosts(driver selenium.WebDriver, likesNeeded int) []Result {
	PageScreenshot(driver, "screen21")

	err := driver.Get("https://www.threads.net/for_you")
	if err != nil {
		log.Fatal("Error:", err)
	}

	time.Sleep(10 * time.Second)
	PageScreenshot(driver, "screen22")

	time.Sleep(4 * time.Second)
	PageScreenshot(driver, "screen23")

	driver.SetPageLoadTimeout(100 * time.Second)

	// // pageScreenshot(driver, "screen1")
	// // time.Sleep(4 * time.Second)
	// // pageScreenshot(driver, "screen2")
	// // time.Sleep(2 * time.Second)

	// // acceptAllCookies(driver)

	// // time.Sleep(2 * time.Second)
	// // pageScreenshot(driver, "screen3")

	// // continueWithInstagram(driver)

	// // time.Sleep(2 * time.Second)
	// // pageScreenshot(driver, "screen4")

	// // pageScreenshot(driver, "screen7")

	// //get cookies
	// // getAllCookies(driver)

	// // lastHeight := 0
	// // newHeight := 0

	// // Прокрутка страницы вниз
	// PageScreenshot(driver, "screen3")

	// script := "window.scrollTo(0, document.body.scrollHeight);"
	// if _, err := driver.ExecuteScript(script, nil); err != nil {
	// 	log.Printf("Ошибка прокрутки: %v", err)
	// }

	// time.Sleep(2 * time.Second)
	// PageScreenshot(driver, "screen4")

	// // Ожидание загрузки новых постов
	// time.Sleep(10 * time.Second)

	// time.Sleep(2 * time.Second)
	// PageScreenshot(driver, "screen5")

	// // Проверка высоты страницы с обработкой типа
	// height, err := driver.ExecuteScript("return document.body.scrollHeight", nil)
	// if err != nil {
	// 	log.Printf("Ошибка получения высоты страницы: %v", err)
	// 	// continue
	// }

	// time.Sleep(2 * time.Second)
	// PageScreenshot(driver, "screen6")

	// newHeight, ok := height.(float64)
	// if !ok {
	// 	log.Printf("Неожиданный тип возвращаемого значения: %T", height)
	// 	// continue
	// }
	// fmt.Printf("n\newHeight: %f\n", newHeight)

	// if int(newHeight) == lastHeight {
	// break // Если прокрутка больше не работает
	// }
	// lastHeight = int(newHeight)
	postsData := []Result{}

	containedData := []string{}

	re := regexp.MustCompile(`<title>Нравится<\/title>[\s\S]*?<div[^>]*>\s*<span[^>]*>(\d+)<\/span>`)
	// re := regexp.MustCompile(`<title>Likes<\/title>[\s\S]*?<div[^>]*>\s*<span[^>]*>(\d+)<\/span>`)

	for len(postsData) < 10 {
		// Поиск всех постов
		posts, err := driver.FindElements(selenium.ByCSSSelector, "div[data-pressable-container=\"true\"]")
		if err != nil {
			log.Printf("Ошибка поиска постов: %v", err)
			// continue
		}
		PageScreenshot(driver, "screen23")

		fmt.Println(len(posts))
		for j := range posts {
			PageScreenshot(driver, fmt.Sprintf("screen%d", j))
			postText, err := posts[j].Text()
			if err == nil {
				html, err := posts[j].GetAttribute("innerHTML")
				if err != nil {
					fmt.Printf("ошибка получения HTML элемента: %v", err)
				}

				if !strings.Contains(html, "<img class") {
					postTextArr := strings.Split(postText, "\n")
					fmt.Printf("\npost array length: %d", len(postTextArr))
					fmt.Println("\npost array:")
					for i := range postTextArr {
						fmt.Printf("\t%s\n", postTextArr[i])
					}
					fmt.Println("\tfull post:")
					fmt.Printf("postTextArr: %v\n", postTextArr)

					// html, err := posts[j].GetAttribute("outerHTML")
					// if err != nil {
					// 	fmt.Printf("ошибка получения HTML элемента: %v", err)
					// }

					for _, contData := range containedData {
						if contData == postTextArr[2] {
							continue
						} else {
							containedData = append(containedData, postTextArr[2])
						}
					}

					// Ищем совпадения
					matches := re.FindStringSubmatch(html)
					if len(matches) > 1 {
						fmt.Println("Найдено число:", matches[1])
						likesNumber, err := strconv.Atoi(matches[1])
						if err != nil {
							likesNumber = 0
						}
						if likesNumber > likesNeeded {
							// FindMoreButton(driver, postTextArr[2])

							time.Sleep(1 * time.Second)
							// PageScreenshot(driver, "hello")

							// targetElement, err := driver.FindElement(selenium.ByXPATH, "//span[text()='Копировать ссылку']/ancestor::*[2]")
							// if err != nil {
							// 	log.Fatal("Не удалось найти элемент:", err)
							// }

							// time.Sleep(time.Duration(CryptoRandom(300, 600)) * time.Millisecond)
							// _, err = driver.ExecuteScript("arguments[0].click();", []interface{}{targetElement})
							// if err != nil {
							// 	panic(fmt.Errorf("не удалось кликнуть по кнопке 'Копировать ссылку': %v", err))
							// }
							// fmt.Println("Успешно нажали на 'Копировать ссылку'")
							// time.Sleep(time.Duration(CryptoRandom(300, 600)) * time.Millisecond)
							// PageScreenshot(driver, "clicked")

							// // Команда для получения текста из буфера
							// cmd := exec.Command("xclip", "-o", "-selection", "clipboard")

							// // Запускаем и получаем вывод
							// output, err := cmd.Output()
							// if err != nil {
							// 	fmt.Println("Ошибка:", err)
							// 	return
							// }

							// // Выводим текст
							// fmt.Println("Текст из буфера:")
							// fmt.Println(string(output))

							//
							//
							//
							//
							//

							targetElement, err := driver.FindElement(selenium.ByXPATH, fmt.Sprintf("//span[text()='%s']/ancestor::*[5]", postTextArr[2]))
							if err != nil {
								log.Printf("не удалось найти элемент: %s", err)
								continue
							}

							PageScreenshot(driver, fmt.Sprintf("screen2-%d", j))

							driver.ExecuteScript("arguments[0].scrollIntoView({block: 'center'});", []interface{}{targetElement})

							PageScreenshot(driver, fmt.Sprintf("screen3-%d", j))

							time.Sleep(time.Duration(CryptoRandom(300, 600)) * time.Millisecond)
							_, err = driver.ExecuteScript("arguments[0].click();", []interface{}{targetElement})
							if err != nil {
								panic(fmt.Errorf("dfas': %v", err))
							}

							time.Sleep(1 * time.Second)
							// PageScreenshot(driver, "clicked")
							PageScreenshot(driver, fmt.Sprintf("screen3-%d", j))

							time.Sleep(10 * time.Second)

							// // Поиск всех постов (обновите селектор)
							// postEntries, err := driver.FindElements(selenium.ByCSSSelector, "div[data-pressable-container=\"true\"]")
							// if err != nil {
							// 	log.Printf("Ошибка поиска постов: %v", err)
							// 	// continue
							// }

							// fmt.Println(len(postEntries))
							// for j := range postEntries {
							// 	postEntity, err := postEntries[j].Text()
							// 	if err == nil {
							// 		postEntityTextArr := strings.Split(postEntity, "\n")
							// 		fmt.Printf("\npost entity array length: %d", len(postEntityTextArr))
							// 		fmt.Println("\npost entity array:")
							// 		for i := range postEntityTextArr {
							// 			fmt.Printf("\t%s\n", postEntityTextArr[i])
							// 		}
							// 		fmt.Println("\tfull entity post:")
							// 		fmt.Printf("postEntityTextArr: %v\n", postEntityTextArr)
							// 	}
							// }

							parsedPost := ParsePostEntities(driver)
							result := Result{}
							err = json.Unmarshal([]byte(parsedPost), &result)
							if err != nil {
								log.Fatal("Ошибка при распарсивании JSON:", err)
							}
							postsData = append(postsData, result)

							time.Sleep(1 * time.Second)
							// PageScreenshot(driver, "after post parsing")
							PageScreenshot(driver, fmt.Sprintf("screen4-%d", j))

							script := "window.scrollTo(0, document.body.scrollHeight);"
							if _, err := driver.ExecuteScript(script, nil); err != nil {
								log.Printf("Ошибка прокрутки: %v", err)
							}

							time.Sleep(1 * time.Second)
							// PageScreenshot(driver, "after scroll post parsing")
							// break
							PageScreenshot(driver, fmt.Sprintf("screen5-%d", j))

							driver.Back()
							time.Sleep(5 * time.Second)
							PageScreenshot(driver, fmt.Sprintf("screen6-%d", j))

							time.Sleep(1 * time.Second)
							// PageScreenshot(driver, "after back")
						}
					} else {
						fmt.Println("Число не найдено")
						// PageScreenshot(driver, "screen when number not find")
						time.Sleep(1 * time.Second)
					}
				}
			}
		}
		// можно листать к последнему элементу из posts
		script := "window.scrollTo(0, document.body.scrollHeight);"
		if _, err := driver.ExecuteScript(script, nil); err != nil {
			log.Printf("Ошибка прокрутки: %v", err)
		}

		time.Sleep(5 * time.Second)
		PageScreenshot(driver, "after scroll")

		break
	}

	// postText0, err := posts[0].Text()
	// if err == nil {
	// 	fmt.Println("\n0 post array:")
	// 	postTextArr := strings.Split(postText0, "\n")
	// 	for i := range postTextArr {
	// 		fmt.Println(postTextArr[i])
	// 	}
	// 	fmt.Println("\n0 full post:")
	// 	fmt.Printf("postTextArr: %v\n", postTextArr)
	// }

	// likesPost0, err := getLikesFromHTML(driver, html)
	// if err == nil {
	// 	fmt.Printf("\n\nlikes from html: %d\n\n", likesPost0)
	// } else {
	// 	fmt.Println("ошибка при парсе лайков из виртуального html")
	// }

	// humanReadPost := []string{}
	// for _, post := range posts {
	// 	postText, err := post.Text()
	// 	if err == nil {
	// 		humanReadPost = append(humanReadPost, postText)
	// 	}
	// }
	// fmt.Printf("\nhumanReadPost: %v\n\n", humanReadPost)

	// time.Sleep(2 * time.Second)
	// PageScreenshot(driver, "screen5")
	// // fmt.Printf("\nposts: %v\n\n", posts)

	// for _, post := range posts {
	// 	likes, err := getLikesCount(post)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	fmt.Printf("likes: %d\n", likes)

	// 	if likes > 499 {
	// 		postInfo, err := extractPostInfo(post)
	// 		if err != nil {
	// 			continue
	// 		}

	// 		fmt.Printf("postInfo: %v\n", postInfo)

	// 		// Проверка на дубликаты
	// 		if !containsPost(postsData, postInfo) {
	// 			postsData = append(postsData, postInfo)

	// 			// Прерывание, если собрали достаточно постов
	// 			if len(postsData) >= 10 {
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// fmt.Printf("postsData: %v\n", postsData)
	// // }

	// time.Sleep(2 * time.Second)
	// PageScreenshot(driver, "screen5")
	// // fmt.Printf("\nposts data: %v\n\n", postsData)

	// // Сохранение данных в JSON
	// jsonData, err := json.MarshalIndent(postsData, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Ошибка маршалинга JSON: %v", err)
	// }

	// fmt.Println(string(jsonData))

	// Или сохранение в файл:
	// if err := os.WriteFile("top_posts.json", jsonData, 0644); err != nil {
	//     log.Fatalf("Ошибка записи в файл: %v", err)
	// }

	return postsData
}

// func getLikesCount(post selenium.WebElement) (int, error) {
// 	// Поиск элемента с количеством лайков (обновите селектор)
// 	likeElement, err := post.FindElement(selenium.ByXPATH,
// 		`.//div[.//svg[@aria-label="Нравится"]]/following-sibling::span//div//span`)
// 	if err != nil {
// 		return 0, fmt.Errorf("не найден элемент с лайками: %v", err)
// 	}

// 	likeText, err := likeElement.Text()
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Преобразование текста в число (может потребоваться обработка формата "1.2K" и т.д.)
// 	likes, err := strconv.Atoi(likeText)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return likes, nil
// }

// func extractPostInfo(post selenium.WebElement) (Post, error) {
// 	postInfo := Post{}

// 	// Получение текста поста (обновите селектор)
// 	if textElement, err := post.FindElement(selenium.ByCSSSelector, "div.x1iorvi4"); err == nil {
// 		if text, err := textElement.Text(); err == nil {
// 			postInfo.Text = text
// 		}
// 	}

// 	// Получение автора (обновите селектор)
// 	if authorElement, err := post.FindElement(selenium.ByCSSSelector, "span.x1lliihq"); err == nil {
// 		if author, err := authorElement.Text(); err == nil {
// 			postInfo.Author = author
// 		}
// 	}

// 	// Получение времени публикации (обновите селектор)
// 	if timeElement, err := post.FindElement(selenium.ByCSSSelector, "span.x1p4m5qa"); err == nil {
// 		if timestamp, err := timeElement.Text(); err == nil {
// 			postInfo.Timestamp = timestamp
// 		}
// 	}

// 	// Получение количества лайков
// 	if likes, err := getLikesCount(post); err == nil {
// 		postInfo.Likes = likes
// 	} else {
// 		return Post{}, err
// 	}

// 	return postInfo, nil
// }

// func containsPost(posts []Post, post Post) bool {
// 	for _, p := range posts {
// 		if p.Likes == post.Likes && p.Text == post.Text && p.Author == post.Author {
// 			return true
// 		}
// 	}
// 	return false
// }
