package statistics

import "encoding/json"

type Jsoncat struct {
	StaticInfo *LogStatic
}

func (jc *Jsoncat) ToJson() []byte{
	json,err:=json.Marshal(jc.StaticInfo)
	if err!=nil{
		logger.Info("jsonMashall statistic info failed!")
	}

	return json
}