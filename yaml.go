package SimpleReverseProxy

import (
	"errors"
	"github.com/helmutkemper/yaml"
	"io/ioutil"
	"reflect"
	"strconv"
)

const KCodeVersion = "0.1 alpha"
const KVersionMinimum = 1.0
const KVersionMaximum = 1.0
const KVersionMinimumString = "1.0"
const KVersionMaximumString = "1.0"
const KSiteErrorInformation = " please, see manual at kemper.com.br for more information"

type configProxyServerNameAndHost struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type configProxy struct {
	Host              string                         `yaml:"host"`
	Path              string                         `yaml:"path"`
	PathExpReg        string                         `yaml:"pathExpReg"`
	QueryStringEnable bool                           `yaml:"queryStringEnable"`
	Server            []configProxyServerNameAndHost `yaml:"server"`
}

type configStaticFolder struct {
	Folder     string `yaml:"folder"`
	ServerPath string `yaml:"serverPath"`
}

type configServer struct {
	ListenAndServer string               `yaml:"listenAndServer"`
	OutputConfig    bool                 `yaml:"outputConfig"`
	StaticServer    bool                 `yaml:"staticServer"`
	StaticFolder    []configStaticFolder `yaml:"staticFolder"`
}

type configReverseProxy struct {
	Config configServer           `yaml:"config"`
	Proxy  map[string]configProxy `yaml:"proxy"`
}

type Config struct {
	Version      string             `yaml:"version"`
	ReverseProxy configReverseProxy `yaml:"reverseProxy"`
}

func (el *Config) Unmarshal(filePath string) error {
	var fileContent []byte
	var err error
	var version float64

	fileContent, err = ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContent, el)
	if err != nil {
		return err
	}

	version, err = strconv.ParseFloat(el.Version, 64)

	if version == 0.0 {
		return errors.New("you must inform the version of config file as numeric value. example: version: '1.0'" + KSiteErrorInformation)
	}

	if version < KVersionMinimum || version > KVersionMaximum {
		return errors.New("this project version accept only configs between versions " + KVersionMaximumString + " and " + KVersionMinimumString + "." + KSiteErrorInformation)
	}

	if reflect.DeepEqual(el.ReverseProxy, configReverseProxy{}) {
		return errors.New("reverse proxy config not found. " + KSiteErrorInformation)
	}

	if el.ReverseProxy.Config.ListenAndServer == "" {
		return errors.New("reverseProxy > config > listenAndServer config not found. " + KSiteErrorInformation)
	}

	if el.ReverseProxy.Config.StaticServer == true && len(el.ReverseProxy.Config.StaticFolder) == 0 {
		return errors.New("reverseProxy > config > staticFolder config not found. " + KSiteErrorInformation)
	}

	for _, folder := range el.ReverseProxy.Config.StaticFolder {
		_, err = ioutil.ReadDir(folder.Folder)
		if err != nil {
			return errors.New("reverseProxy > config > staticFolder error: " + err.Error())
		}
	}

	if len(el.ReverseProxy.Proxy) == 0 {
		return errors.New("reverseProxy > proxy not found. " + KSiteErrorInformation)
	}

	for proxyName, proxyConfig := range el.ReverseProxy.Proxy {
		//if proxyConfig.Host == "" {
		//	return errors.New("reverseProxy > proxy > " + proxyName + " > host not found. " + KSiteErrorInformation)
		//}

		if len(proxyConfig.Server) == 0 {
			return errors.New("reverseProxy > proxy > " + proxyName + " > server not found. " + KSiteErrorInformation)
		}

		for _, proxyServerConfig := range proxyConfig.Server {
			if proxyServerConfig.Host == "" {
				return errors.New("reverseProxy > proxy > " + proxyName + " > server > host not found. " + KSiteErrorInformation)
			}

			if proxyServerConfig.Name == "" {
				return errors.New("reverseProxy > proxy > " + proxyName + " > server > name not found. " + KSiteErrorInformation)
			}
		}
	}

	return nil
}
