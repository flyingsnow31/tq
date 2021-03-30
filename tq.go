package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const key = "f80464c47d16849da62561a10151d1b7"

func main() {
	argNum := len(os.Args)
	//dirwd,_:=os.Getwd()
	home, err := Home()
	if err != nil {
		fmt.Println("目录寻找失败")
		return
	}
	checkDir(home)
	dir := home + "/AppData/Local/tq/data/city.dat"
	fmt.Println(dir)
	fil, err := os.OpenFile(dir, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return
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
			msg, ok := getWeatherNow(string(buf))
			if !ok {
				fmt.Println(msg)
				return
			}
		}
	case 2:
		{
			date := os.Args[1]
			if date == "h" { //帮助信息
				printHelp()
				return
			} else if date == "a" { //未来第1、2、3天的天气
				if len(string(buf)) == 0 {
					fmt.Println("please use 'tq set <cityname>' to init.")
					return
				}
				msg, ok := getWeatherFea(string(buf))
				if !ok {
					fmt.Println(msg)
					return
				}
			} else { //未来三天
				if len(string(buf)) == 0 {
					fmt.Println("please use 'tq set <cityname>' to init.")
					return
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
			adcode, ok := getCityCode(cityName)
			if adcode == nil {
				return
			}
			if ok {
				fil, _ := os.OpenFile(home+"/AppData/Local/tq/data/city.dat", os.O_TRUNC|os.O_RDWR, 0777)
				defer func() {
					_ = fil.Close()
				}()
				cityNum := len(adcode)
				if cityNum == 1 {
					for code := range adcode {
						_, _ = fil.WriteString(code)
					}
					return
				} else {
					fmt.Println("出现重名，请选择:")
					tmp := make([]string, cityNum)
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

func checkDir(home string) {
	dir := home + "/AppData/Local/tq/data"
	//fmt.Println(dir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
}

func printHelp() {
	fmt.Println("hello,world")
}

func getCityCode(cityName string) (map[string]string, bool) {
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

func getWeatherFea(adcode string) (string, bool) {
	wea, ok := InQuire(adcode, false)
	if !ok {
		return "net req error", false
	} else {
		var ret WeatherInfoFea
		err := json.Unmarshal([]byte(wea), &ret)
		if err != nil {
			return "格式错误", false
		}
		str, _ := json.MarshalIndent(ret, "\t", "    ")
		fmt.Println(string(str))
		return "执行成功", true
	}
}

func getWeatherNow(adcode string) (string, bool) {
	wea, ok := InQuire(adcode, true)
	if !ok {
		return "net req error", false
	} else {
		var ret WeatherInfoNow
		err := json.Unmarshal([]byte(wea), &ret)
		if err != nil {
			return "格式错误", false
		}
		str, _ := json.MarshalIndent(ret, "\t", "    ")
		fmt.Println(string(str))
		return "执行成功", true
	}
}

func InQuire(adcode string, isNow bool) (string, bool) {
	var extensions string
	if isNow {
		extensions = "base"
	} else {
		extensions = "all"
	}
	rlt, err := doWeatherNow("https://restapi.amap.com/v3/weather/weatherInfo?key=" + key + "&city=" + adcode + "&extensions=" + extensions)
	if err != nil {
		return "", false
	} else {
		return rlt, true
		//fmt.Println(rlt)
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
