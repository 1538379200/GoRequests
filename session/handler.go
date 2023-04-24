package session

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type Handler struct {
	json string
}

func (h *Handler) Find(path string) gjson.Result {
	res := gjson.Get(h.json, path)
	return res
}

func (h *Handler) Json() string {
	return h.json
}

func (h *Handler) JsonFormat() string {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(h.json), &m)
	val, _ := json.MarshalIndent(m, "", "  ")
	return string(val)
}
