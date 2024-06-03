package main

type OfficeConvertRequest struct {
	Input string `json:"input"`
}

type OfficeConvertReply struct {
	Outputs []string `json:"outputs"`
}
