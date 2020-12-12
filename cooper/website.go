//    Copyright (C) 2020 Gary Kim <gary@garykim.dev>, All Rights Reserved
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as published
//    by the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cooper

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrorUnexpectedReturnCode = errors.New("unexpected return code")
)

type CourseInfo struct {
	Code        string
	Codes       []string
	Name        string
	Description string
	ExtraInfo   string
}

func ScrapeInfo() ([]CourseInfo, error) {
	var tr []CourseInfo
	urls := []string{"https://cooper.edu/engineering/curriculum/courses", "https://cooper.edu/humanities/curriculum/courses", "https://cooper.edu/art/curriculum/courses", "https://cooper.edu/architecture/curriculum/courses"}
	for _, url := range urls {
		t, err := scrapePage(url)
		if err != nil {
			return nil, err
		}
		tr = append(tr, t...)
	}
	return tr, nil
}

func scrapePage(url string) ([]CourseInfo, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, ErrorUnexpectedReturnCode
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var tr []CourseInfo

	doc.Find("#course-listings .content li").Each(func(i int, selection *goquery.Selection) {
		codes := splitCodes(selection.Find("h3").Text())
		for _, courseCode := range codes {
			description := selection.Find("p").First().Text()
			tr = append(tr, CourseInfo{
				Code:        strings.ToUpper(strings.TrimSpace(courseCode)),
				Codes:       codes,
				Name:        strings.TrimSpace(selection.Find("h4").Text()),
				Description: strings.TrimSpace(description),
				ExtraInfo:   strings.TrimSpace(strings.TrimPrefix(selection.Find("p").Text(), description)),
			})
		}
	})
	return tr, nil
}

func splitCodes(code string) (tr []string) {
	var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code = strings.TrimSpace(code)
	for _, c := range strings.Split(code, ",") {
		toset := []string{strings.TrimSpace(c)}
		if ran := strings.Split(c, "-"); len(ran) == 2 {
			toset = []string{}
			arr := regexp.MustCompile("[\\.\\s]+").Split(strings.TrimSpace(ran[0]), -1)
			toChange := arr[len(arr)-1]
			stayStill := ran[0][0:strings.LastIndex(ran[0], toChange)]
			toChange = strings.TrimSpace(toChange)

			arr = regexp.MustCompile("[\\.\\s]+").Split(ran[1], -1)

			target := strings.TrimSpace(arr[len(arr)-1])
			if num, err := strconv.Atoi(toChange); err == nil {
				for i := 0; i < 300; i++ {
					toset = append(toset, stayStill+strconv.Itoa(num+i))
					if strconv.Itoa(num+i) == target {
						break
					}
				}
			} else {
				for i := strings.Index(alphabet, string(toChange[len(toChange)-1])); i < len(alphabet); i++ {
					toset = append(toset, stayStill+string(alphabet[i]))
					if alphabet[i] == target[len(target)-1] {
						break
					}
				}
			}
		}
		tr = append(tr, toset...)
	}
	return
}
