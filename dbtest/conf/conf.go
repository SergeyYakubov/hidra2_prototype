package conf

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
				 Server          string
				 EnsureDiskWrite bool
				 Name            string
			 }
	Nthreads int
	Consumer struct {
				 GroupID int
				 NImagesPerRead int
			 }
}

func ReadConfig(fname string, config interface{}) error {
	yamlFile, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		return err
	}
	return nil
}

