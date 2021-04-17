package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const key = "f80464c47d16849da62561a10151d1b7"

var a = flag.Bool("a", false, "查看设定城市未来三天天气")
var s = flag.String("s", "", "查询城市名并进行绑定")
var h = flag.Bool("h", false, "查看帮助菜单")
var t = flag.String("t", "", "查询城市名并进行绑定")

func main() {
	flag.Parse()
	//dirwd,_:=os.Getwd()
	home, err := Home()
	if err != nil {
		fmt.Println("目录寻找失败")
		return
	}
	dir := conf(home)
	//fmt.Println(dir)
	buf, ok := getBuf(dir)
	if !ok {
		return
	}
	isBufExist := checkBuf(buf)
	if len(os.Args) == 1 {
		if !isBufExist {
			return
		}
		msg, ok := getWeatherNow(buf)
		if !ok {
			fmt.Println(msg)
		}
		return
	}
	if *a == true {
		if !isBufExist {
			return
		}
		msg, ok := getWeatherFea(buf)
		if !ok {
			fmt.Println(msg)
			return
		}
	}
	if *h == true {
		printHelp()
		return
	}
	if len(*s) > 0 {
		cityName := os.Args[2]
		adcode, ok := getCityCode(cityName)
		if adcode == nil {
			return
		}
		if ok {
			fil, _ := os.OpenFile(dir, os.O_TRUNC|os.O_RDWR, 0777)
			defer func() {
				_ = fil.Close()
			}()
			cityNum := len(adcode)
			if cityNum == 1 {
				for code, name := range adcode {
					_, _ = fil.WriteString(code)
					fmt.Println(name, "设置成功")
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
						fmt.Println(adcode[tmp[index]], "设置成功")
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
	if len(*t) > 0 {
		cityName := os.Args[2]
		adcode, ok := getCityCode(cityName)
		if adcode == nil {
			return
		}
		if ok {
			cityNum := len(adcode)
			if cityNum == 1 {
				for code := range adcode {
					msg, ok := getWeatherNow(code)
					if !ok {
						fmt.Println(msg)
					}
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
						msg, ok := getWeatherNow(tmp[index])
						if !ok {
							fmt.Println(msg)
						}
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
}

func getBuf(dir string) (string, bool) {
	fil, err := os.OpenFile(dir, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	buf, err := ioutil.ReadAll(fil)
	_ = fil.Close()
	return string(buf), true
}

func checkBuf(buf string) bool {
	if len(buf) == 0 {
		fmt.Println("please use 'tq set <cityname>' to init.")
		return false
	}
	return true
}

func getCityCode(cityName string) (map[string]string, bool) {
	rlt, err := doWeatherNow("https://restapi.amap.com/v3/geocode/geo?key=" + key + "&address=" + cityName)
	if err != nil {
		ret := map[string]string{
			"msg": "net req error",
		}
		return ret, false
	} else {
		//fmt.Println(rlt)
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
		str, _ := json.MarshalIndent(ret.Forecasts, "", "    ")
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
		str, _ := json.MarshalIndent(ret.Lives, "", "    ")
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

func printHelp() {
	fmt.Println("tq 查看所设定的城市天气\n" +
		"------------------------\n" +
		"-a 查看未来三天天气\n" +
		"-h 查看帮助\n" +
		"-s 设定城市\n" +
		"------------------------\n" +
		"                 by lzy")
}
