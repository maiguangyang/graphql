package model

import(
  "strings"
  "regexp"
)

// 正则替换
func ReplaceAll(data, reg, target string) string {
  req, _ := regexp.Compile(reg);
  rep := req.ReplaceAllString(data, target);
  return rep
}

// 正则截取
func RegexpReplace(str, start string, end string) string {
  reg, _ := regexp.Compile(start + ".+?" + end)
  value := reg.FindString(str)

  reg = regexp.MustCompile(start)
  value = reg.ReplaceAllString(value, "")

  reg = regexp.MustCompile(end)
  value = reg.ReplaceAllString(value, "")
  return value
}

func ToLower(str string) string {
  return strings.ToLower(str)
}