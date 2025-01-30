package entity

import "time"

type LogEntry struct {
	ClientAddr            string    `json:"ClientAddr"`
	ClientHost            string    `json:"ClientHost"`
	ClientPort            string    `json:"ClientPort"`
	ClientUsername        string    `json:"ClientUsername"`
	DownstreamContentSize int64     `json:"DownstreamContentSize"`
	DownstreamStatus      int       `json:"DownstreamStatus"`
	Duration              int64     `json:"Duration"`
	GzipRatio             float64   `json:"GzipRatio"`
	OriginContentSize     int64     `json:"OriginContentSize"`
	OriginDuration        int64     `json:"OriginDuration"`
	OriginStatus          int       `json:"OriginStatus"`
	Overhead              int64     `json:"Overhead"`
	RequestAddr           string    `json:"RequestAddr"`
	RequestContentSize    int64     `json:"RequestContentSize"`
	RequestCount          int       `json:"RequestCount"`
	RequestHost           string    `json:"RequestHost"`
	RequestMethod         string    `json:"RequestMethod"`
	RequestPath           string    `json:"RequestPath"`
	RequestPort           string    `json:"RequestPort"`
	RequestProtocol       string    `json:"RequestProtocol"`
	RequestScheme         string    `json:"RequestScheme"`
	RetryAttempts         int       `json:"RetryAttempts"`
	RouterName            string    `json:"RouterName"`
	ServiceAddr           string    `json:"ServiceAddr"`
	ServiceName           string    `json:"ServiceName"`
	ServiceURL            string    `json:"ServiceURL"`
	StartLocal            time.Time `json:"StartLocal"`
	StartUTC              time.Time `json:"StartUTC"`
	TLSCipher             string    `json:"TLSCipher"`
	TLSVersion            string    `json:"TLSVersion"`
	EntryPointName        string    `json:"entryPointName"`
	Level                 string    `json:"level"`
	Msg                   string    `json:"msg"`
	Time                  time.Time `json:"time"`
}
