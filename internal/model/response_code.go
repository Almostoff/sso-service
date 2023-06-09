package model

import (
	"AuthService/internal/cConstants"
	"AuthService/pkg/loggerService"
	"encoding/json"
	"fmt"
)

type OmitZeroInt int

func (oi OmitZeroInt) MarshalJSON() ([]byte, error) {
	if oi == 0 {
		return []byte("null"), nil
	}
	return json.Marshal(int(oi))
}

type CodeModel struct {
	IsError         bool        `json:"-"`
	InternalCode    OmitZeroInt `json:"internal_code,omitempty"`
	StandardCode    int64       `json:"code,omitempty"`
	InternalMessage string      `json:"internal_message,omitempty"`
	StandardMessage string      `json:"message,omitempty"`
}

type Success struct {
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	RedirectUrl string `json:"redirect_url"`
}

func GetSuccess(suc bool, mes string) *Success {
	return &Success{Success: suc, Message: mes}
}

func GetSuccessByCM(cm *CodeModel) *Success {
	var ok bool
	if !cm.IsError {
		ok = true
	}
	return GetSuccess(ok, "")
}

func GetError(interCode OmitZeroInt, sdrdCode int64, interMes string, sdrdMes string) *CodeModel {
	return &CodeModel{IsError: true, InternalMessage: interMes, InternalCode: interCode, StandardCode: sdrdCode,
		StandardMessage: sdrdMes}
}

func GetResponseCode(cm *CodeModel, data interface{}, params interface{}) *CodeModel {
	if cm.IsError || cm.InternalCode != 0 {
		errMes := fmt.Sprintf("Произошла ошибка!\nERROR:\n %+v \nPARAMS:\n %+v,\n "+
			"DATA:\n %+v", cm, params, data)
		go loggerService.GetInstance().DevLog(errMes, 10)
		fmt.Println(fmt.Errorf(errMes))
		if !cConstants.Debug {
			cm.InternalMessage = ""
			cm.InternalCode = 0
		}
		fmt.Println("here")
		return cm
	} else {
		fmt.Println("here1")
		return &CodeModel{}
	}

	return cm
}
