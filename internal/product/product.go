package product

import (
	"encoding/json"
	"fmt"
	"github.com/run-bigpig/jb-active/internal/model"
	"github.com/run-bigpig/jb-active/internal/utils"
	"os"
	"path/filepath"
	"sync"
)

const (
	pluginBaseUrl = "https://plugins.jetbrains.com"
	ideBaseUrl    = "https://data.services.jetbrains.com/products"
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"
	referer       = "https://www.jetbrains.com/"
	ideName       = "ide.json"
	pluginName    = "plugin.json"
)

var header map[string]string

func init() {
	header = make(map[string]string)
	header["User-Agent"] = userAgent
	header["Referer"] = referer
}

func getIdeList() ([]model.Ide, error) {
	var response []model.Ide
	ideApi := fmt.Sprintf("%s?%s", ideBaseUrl, "fields=code,name,description")
	result, err := utils.SendRequest(ideApi, header, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(result, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return response, nil
}

func getPaidPluginList() ([]model.Plugin, error) {
	var (
		response       model.PluginList
		paidPluginList []model.Plugin
	)
	pluginApi := fmt.Sprintf("%s%s", pluginBaseUrl, "/api/searchPlugins?max=10000&offset=0")
	result, err := utils.SendRequest(pluginApi, header, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = json.Unmarshal(result, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, plugin := range response.Plugins {
		if plugin.PricingModel == "PAID" {
			plugin.Icon = fmt.Sprintf("%s%s", pluginBaseUrl, plugin.Icon)
			paidPluginList = append(paidPluginList, plugin)
		}
	}
	return paidPluginList, nil
}

func getPaidPluginCode(id int64) (*model.PluginCode, error) {
	var response model.PluginCode
	pluginApi := fmt.Sprintf("%s%s%d", pluginBaseUrl, "/api/plugins/", id)
	result, err := utils.SendRequest(pluginApi, header, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func savePlugin() error {
	var (
		wg   sync.WaitGroup
		sem  = make(chan struct{}, 10)
		lock sync.Mutex
	)
	codeList := make([]*model.Ide, 0)
	pluginList, err := getPaidPluginList()
	if err != nil {
		return err
	}
	for _, plugin := range pluginList {
		wg.Add(1)
		sem <- struct{}{}
		go func(plugin model.Plugin) {
			defer func() {
				wg.Done()
				<-sem
			}()
			code, err := getPaidPluginCode(plugin.ID)
			if err != nil {
				return
			}
			lock.Lock()
			codeList = append(codeList, &model.Ide{
				Code: code.PurchaseInfo.ProductCode,
				Name: code.Name,
			})
			lock.Unlock()
		}(plugin)
	}
	wg.Wait()
	//填满信号量 保证所有goroutine执行完毕
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}
	jsonData, err := json.Marshal(codeList)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(utils.GetStaticPath(), "json", pluginName), jsonData, 0666)
}

func saveIde() error {
	ideList, err := getIdeList()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(ideList)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(utils.GetStaticPath(), "json", ideName), jsonData, 0666)
}

func SyncProduct() error {
	err := saveIde()
	if err != nil {
		return err
	}
	return savePlugin()
}
