package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const key = "f80464c47d16849da62561a10151d1b7"

func main() {
	argNum := len(os.Args)
	checkDir()
	fil, err := os.OpenFile("./data/city.dat", os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println(err.Error())
	}
	buf, err := ioutil.ReadAll(fil)
	_ = fil.Close()
	switch argNum {
	case 1:
		{
			if len(string(buf)) == 0 {
				fmt.Println("please use 'tq set <cityname>' to init.")
				return
			}
			InQuire(string(buf))
		}
	case 2:
		{
			date := os.Args[1]
			find := strings.Count(date, ".")
			if find == 1 {
				if len(string(buf)) == 0 {
					fmt.Println("please use 'tq set <cityname>' to init.")
					return
				}
				forecastDay(string(buf), date)
			} else if find > 1 {
				if len(string(buf)) == 0 {
					fmt.Println("please use 'tq set <cityname>' to init.")
					return
				}
				fmt.Println("para error, use 'tq h'")
				return
			} else {
				if date == "h" {
					printHelp()
					return
				} else {
					forecastNum(string(buf), date)
				}
			}
		}
	case 3:
		{
			para := os.Args[1]
			if para != "set" {
				fmt.Println("para error, use 'tq h'")
				return
			}
			cityName := os.Args[2]
			adcode, ok := InQuireForecase(cityName)
			if adcode == nil {
				return
			}
			if ok {
				fil, _ := os.OpenFile("./data/city.dat", os.O_TRUNC|os.O_RDWR, 0777)
				cityNum := len(adcode)
				if cityNum == 1 {
					for code, _ := range adcode {
						_, _ = fil.WriteString(code)
					}
					return
				} else {
					fmt.Println("出现重名，请选择:")
					var tmp []string
					index := 0
					for code, name := range adcode {
						tmp = append(tmp, code)
						fmt.Printf("%d %s", index, name)
						index++
					}
					for {
						_, _ = fmt.Scanf("%d", index)
						if index > cityNum-1 || index < 0 {
							fmt.Println("输入错误，请重试")
						} else {
							_, _ = fil.WriteString(tmp[index])
							break
						}
					}
					return
				}
			} else {
				fmt.Println(adcode["msg"])
				return
			}
		}
	default:
		{
			fmt.Println("para error, use 'tq h'")
			return
		}
	}

}

func checkDir() {
	dir := "./data"
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		_ = os.Mkdir(dir, os.ModePerm)
	}
}

func printHelp() {
	fmt.Println("hello,world")
}

func forecastDay(adcode, date string) {

}

func forecastNum(adcode, date string) {

}

func InQuireForecase(cityName string) (map[string]string, bool) {
	rlt, err := doWeatherNow("https://restapi.amap.com/v3/geocode/geo?key=" + key + "&address=" + cityName)
	if err != nil {
		fmt.Println("net req error")
		return nil, false
	} else {
		fmt.Println(rlt)
		var geoCode GeoCode
		_ = json.Unmarshal([]byte(rlt), &geoCode)
		if geoCode.Count == "0" {
			ret := map[string]string{
				"msg": "城市名错误",
			}
			return ret, false
		} else {
			ret := make(map[string]string)
			for _, v := range geoCode.Geocodes {
				ret[v.Adcode] = v.FormattedAddress
			}
			return ret, true
		}
	}
}

func InQuire(adcode string) {
	rlt, err := doWeatherNow("https://restapi.amap.com/v3/weather/weatherInfo?key=" + key + "&city=" + adcode)
	if err != nil {
		fmt.Println("net req error")
	} else {
		fmt.Println(rlt)
	}
}

func doWeatherNow(url string) (rlt string, err error) {

	resp, err := http.Get(url)

	if err != nil {
		return "", err
	} else {

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		} else {
			return string(body), err
		}
	}
}