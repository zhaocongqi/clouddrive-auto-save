package cloud139

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zcq/clouddrive-auto-save/internal/core"
	"github.com/zcq/clouddrive-auto-save/internal/db"
)

const (
	BaseURL           = "https://yun.139.com"
	UserNjsURL        = "https://user-njs.yun.139.com"
	ShareKdNjsURL     = "https://share-kd-njs.yun.139.com"
	PersonalKdNjsURL  = "https://personal-kd-njs.yun.139.com"
	CatalogV1URL      = BaseURL + "/orchestration/personalCloud/catalog/v1.0"
)

type Cloud139 struct {
	account *db.Account
	client  *http.Client
}

func init() {
	core.RegisterDriver("139", func(account *db.Account) core.CloudDrive {
		return NewCloud139(account)
	})
}

func NewCloud139(account *db.Account) *Cloud139 {
	return &Cloud139{
		account: account,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// ─── 辅助工具 ─────────────────────────────────────────────────────────────────

func md5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// jsEncodeURIComponent 模拟 JS 的 encodeURIComponent，确保空格转为 %20 而非 +
func jsEncodeURIComponent(str string) string {
	t := url.QueryEscape(str)
	t = strings.ReplaceAll(t, "+", "%20")
	t = strings.ReplaceAll(t, "%21", "!")
	t = strings.ReplaceAll(t, "%27", "'")
	t = strings.ReplaceAll(t, "%28", "(")
	t = strings.ReplaceAll(t, "%29", ")")
	t = strings.ReplaceAll(t, "%2A", "*")
	t = strings.ReplaceAll(t, "%7E", "~")
	return t
}

func (c *Cloud139) getPhone() string {
	re := regexp.MustCompile(`1\d{10}`) // 匹配 11 位手机号

	// 1. 如果 AccountName 本身就是 11 位手机号，直接返回
	if re.MatchString(c.account.AccountName) {
		return re.FindString(c.account.AccountName)
	}

	// 2. 尝试从 AuthToken 中提取
	auth := c.account.AuthToken
	if auth != "" {
		token := auth
		if strings.HasPrefix(strings.ToLower(auth), "basic ") {
			token = auth[6:]
		}
		decoded, err := base64.StdEncoding.DecodeString(token)
		if err == nil {
			if match := re.FindString(string(decoded)); match != "" {
				return match
			}
		}
	}

	// 3. 尝试从 Cookie 中提取（通常 Cookie 中会包含手机号）
	if c.account.Cookie != "" {
		if match := re.FindString(c.account.Cookie); match != "" {
			return match
		}
	}

	return ""
}

func getNewSignHash(body interface{}, datetime, randomStr string) string {
	s := ""
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		s = jsEncodeURIComponent(string(jsonBytes))
		chars := strings.Split(s, "")
		sort.Strings(chars)
		s = strings.Join(chars, "")
	}
	r := md5Hash(base64.StdEncoding.EncodeToString([]byte(s)))
	c := md5Hash(datetime + ":" + randomStr)
	return strings.ToUpper(md5Hash(r + c))
}

func (c *Cloud139) computeMcloudSign(catalogID string) string {
	now := time.Now().In(time.FixedZone("CST", 8*3600))
	datetime := now.Format("2006-01-02 15:04:05")
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomStr := ""
	for i := 0; i < 16; i++ {
		randomStr += string(chars[rand.Intn(len(chars))])
	}

	getDiskBody := map[string]interface{}{
		"catalogID":       catalogID,
		"sortDirection":   1,
		"startNumber":     1,
		"endNumber":       100,
		"filterType":      0,
		"catalogSortType": 0,
		"contentSortType": 0,
		"commonAccountInfo": map[string]interface{}{
			"account":     c.getPhone(),
			"accountType": 1,
		},
	}
	hash := getNewSignHash(getDiskBody, datetime, randomStr)
	return fmt.Sprintf("%s,%s,%s", datetime, randomStr, hash)
}

// ─── HTTP 请求封装 ─────────────────────────────────────────────────────────────

func (c *Cloud139) doRequest(ctx context.Context, method, apiURL string, body interface{}, headers map[string]string) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// 1. 设置身份认证
	if c.account.AuthToken != "" {
		auth := c.account.AuthToken
		if !strings.HasPrefix(strings.ToLower(auth), "basic ") {
			auth = "Basic " + auth
		}
		req.Header.Set("Authorization", auth)
	} else if c.account.Cookie != "" {
		req.Header.Set("Cookie", c.account.Cookie)
	}

	// 2. 设置通用 Web 模拟 Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://yun.139.com/")
	req.Header.Set("Origin", "https://yun.139.com")
	req.Header.Set("x-yun-api-version", "v1")
	req.Header.Set("x-yun-app-channel", "10000034")
	req.Header.Set("x-yun-channel-source", "10000034")
	req.Header.Set("x-yun-client-info", "||9|7.17.2|chrome|143.0.0.0|ff559f01db65afce55f3b4e5d75be4cb||windows 10||zh-CN|||")
	req.Header.Set("x-yun-module-type", "100")
	req.Header.Set("x-yun-svc-type", "1")
	req.Header.Set("mcloud-channel", "1000101")
	req.Header.Set("mcloud-version", "7.17.2")
	req.Header.Set("mcloud-client", "10701")
	req.Header.Set("mcloud-route", "001")

	// 3. 针对不同域名的特殊处理
	if strings.Contains(apiURL, "personal-kd-njs") {
		req.Header.Set("INNER-HCY-ROUTER-HTTPS", "1")
		req.Header.Set("x-m4c-caller", "PC")
		req.Header.Set("x-m4c-src", "10002")
		req.Header.Set("x-inner-ntwk", "2")
		req.Header.Set("X-Deviceinfo", "||9|7.17.2|chrome|143.0.0.0|ff559f01db65afce55f3b4e5d75be4cb||windows 10||zh-CN|||")
		req.Header.Set("CMS-DEVICE", "default")
		req.Header.Set("x-huawei-channelSrc", "10000034")
		req.Header.Set("x-SvcType", "1")
	} else if strings.Contains(apiURL, "share-kd-njs") {
		req.Header.Set("caller", "web")
		req.Header.Set("x-m4c-caller", "PC")
		// 分享接口不需要 mcloud-sign，如果有则移除
		delete(headers, "mcloud-sign")
	}

	// 4. 覆盖/添加传入的自定义 Headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ─── CloudDrive 接口实现 ───────────────────────────────────────────────────────

func (c *Cloud139) GetInfo(ctx context.Context) (*db.Account, error) {
	headers := map[string]string{
		"caller":         "web",
		"x-m4c-caller":   "PC",
		"mcloud-version": "7.17.2",
		"mcloud-client":  "10701",
	}
	resp, err := c.doRequest(ctx, "POST", UserNjsURL+"/user/getUser", map[string]interface{}{}, headers)
	if err != nil {
		return nil, err
	}

	var resRaw map[string]interface{}
	if err := json.Unmarshal(resp, &resRaw); err != nil {
		return nil, err
	}

	// 打印调试日志，方便定位昵称获取失败原因
	// fmt.Printf("[139 Debug] getUser response: %s\n", string(resp))

	code := ""
	switch v := resRaw["code"].(type) {
	case string:
		code = v
	case float64:
		code = fmt.Sprintf("%.0f", v)
	}

	if code != "0000" && code != "0" && code != "" {
		return nil, fmt.Errorf("API error: %s", code)
	}

	// 139 的核心数据可能在顶层、data 或 result 节点下，进行深度探测
	data, _ := resRaw["data"].(map[string]interface{})
	if data == nil {
		data, _ = resRaw["result"].(map[string]interface{})
	}
	if data == nil {
		data = resRaw
	}

	nickname, _ := data["userName"].(string)
	if nickname == "" {
		nickname, _ = data["nickName"].(string)
	}
	// 关键修复：探测 userProfileInfo 节点
	if nickname == "" {
		if profile, ok := data["userProfileInfo"].(map[string]interface{}); ok {
			nickname, _ = profile["userName"].(string)
		}
	}

	phone, _ := data["loginName"].(string)
	if phone == "" {
		phone, _ = data["account"].(string)
	}
	if phone == "" {
		phone, _ = data["msisdn"].(string)
	}
	if profile, ok := data["userProfileInfo"].(map[string]interface{}); ok {
		if phone == "" {
			phone, _ = profile["msisdn"].(string)
		}
		if phone == "" {
			phone, _ = profile["loginAccount"].(string)
		}
		if phone == "" {
			phone, _ = profile["account"].(string)
		}
	}

	userDomainID, _ := data["userDomainId"].(string)

	if nickname == "" {
		nickname = phone
	}
	if nickname == "" {
		nickname = c.account.AccountName
	}
	if nickname == "" {
		nickname = "移动云盘用户"
	}

	c.account.Nickname = nickname
	c.account.Status = 1
	c.account.LastCheck = time.Now()
	
	// 关键修复：从获取到的任何包含手机号的字符串中精准提取 11 位数字
	rePhone := regexp.MustCompile(`1\d{10}`)
	if rePhone.MatchString(phone) {
		c.account.AccountName = rePhone.FindString(phone)
	} else if c.account.Cookie != "" && rePhone.MatchString(c.account.Cookie) {
		// 最后的尝试：从 Cookie 字符串中提取
		c.account.AccountName = rePhone.FindString(c.account.Cookie)
	} else if c.account.AccountName == "" {
		c.account.AccountName = nickname
	}

	// 尝试从 auth 节点获取会员等级
	auth, _ := resRaw["auth"].(map[string]interface{})
	vipLevels := map[int]string{0: "普通用户", 1: "白银会员", 2: "黄金会员", 3: "钻石会员"}
	if auth != nil {
		if ml, ok := auth["memberLevel"].(float64); ok && ml > 0 {
			c.account.VipName = vipLevels[int(ml)]
		}
	}

	// 2. 获取磁盘容量信息
	if userDomainID != "" {
		diskReq := map[string]interface{}{"userDomainId": userDomainID}
		diskResp, err := c.doRequest(ctx, "POST", UserNjsURL+"/user/disk/getPersonalDiskInfo", diskReq, headers)
		if err == nil {
			var diskRes struct {
				Data struct {
					DiskSize     string `json:"diskSize"`
					FreeDiskSize string `json:"freeDiskSize"`
				} `json:"data"`
			}
			if json.Unmarshal(diskResp, &diskRes) == nil {
				total, _ := strconv.ParseInt(diskRes.Data.DiskSize, 10, 64)
				free, _ := strconv.ParseInt(diskRes.Data.FreeDiskSize, 10, 64)
				c.account.CapacityTotal = total * 1024 * 1024 // MB to B
				c.account.CapacityUsed = (total - free) * 1024 * 1024
			}
		}
	}

	// 3. 获取会员等级 (如果第一步没拿到)
	if c.account.VipName == "" || c.account.VipName == "普通用户" {
		// 优先尝试从 AuthToken (Basic auth) 中提取手机号
		if c.account.AuthToken != "" {
			authVal := strings.TrimSpace(strings.TrimPrefix(c.account.AuthToken, "Basic "))
			authVal = strings.TrimPrefix(authVal, "basic ")
			if decoded, err := base64.StdEncoding.DecodeString(authVal); err == nil {
				parts := strings.Split(string(decoded), ":")
				if len(parts) >= 2 && len(parts[1]) == 11 && strings.HasPrefix(parts[1], "1") {
					phone = parts[1]
				}
			}
		}

		if phone == "" {
			phone = c.account.AccountName
		}
		benefitReq := map[string]interface{}{
			"isNeedBenefit": 1,
			"commonAccountInfo": map[string]interface{}{
				"account":     phone,
				"accountType": 1,
			},
		}
		// 使用精简 Headers (nil 则 doRequest 使用默认基础 Header)，避免被服务器风控或结构偏移
		benefitResp, err := c.doRequest(ctx, "POST", BaseURL+"/orchestration/group-rebuild/member/v1.0/queryUserBenefits", benefitReq, nil)
		if err == nil {
			var benefitRaw map[string]interface{}
			if json.Unmarshal(benefitResp, &benefitRaw) == nil {
				// 多路径探测：data.userSubMemberList -> data.userMemberList -> userSubMemberList
				var memberList []interface{}
				if data, ok := benefitRaw["data"].(map[string]interface{}); ok {
					if list, ok := data["userSubMemberList"].([]interface{}); ok {
						memberList = list
					} else if list, ok := data["userMemberList"].([]interface{}); ok {
						memberList = list
					}
				} else if list, ok := benefitRaw["userSubMemberList"].([]interface{}); ok {
					memberList = list
				} else if list, ok := benefitRaw["userMemberList"].([]interface{}); ok {
					memberList = list
				}

				if len(memberList) > 0 {
					if first, ok := memberList[0].(map[string]interface{}); ok {
						level := -1
						// 灵活类型探测
						v := first["memberLevel"]
						if v == nil {
							v = first["level"]
						}
						switch val := v.(type) {
						case float64:
							level = int(val)
						case string:
							level, _ = strconv.Atoi(val)
						case int:
							level = val
						}

						if level != -1 {
							if name, ok := vipLevels[level]; ok {
								c.account.VipName = name
							} else {
								c.account.VipName = fmt.Sprintf("会员L%d", level)
							}
						}
					}
				}
			}
		}
	}

	return c.account, nil
}

func (c *Cloud139) Login(ctx context.Context) error {
	_, err := c.GetInfo(ctx)
	return err
}

func (c *Cloud139) ListFiles(ctx context.Context, parentID string) ([]core.FileInfo, error) {
	if parentID == "" {
		parentID = "/"
	}
	sign := c.computeMcloudSign(parentID)
	headers := map[string]string{
		"mcloud-sign":            sign,
		"mcloud-version":         "7.17.2",
		"mcloud-channel":         "1000101",
		"mcloud-client":          "10701",
		"INNER-HCY-ROUTER-HTTPS": "1",
	}

	body := map[string]interface{}{
		"pageInfo": map[string]interface{}{
			"pageSize":   100,
			"pageCursor": nil,
		},
		"orderBy":         "updated_at",
		"orderDirection":  "DESC",
		"parentFileId":    parentID,
	}

	resp, err := c.doRequest(ctx, "POST", PersonalKdNjsURL+"/hcy/file/list", body, headers)
	if err != nil {
		return nil, err
	}

	var res struct {
		Code    string `json:"code"`
		Success bool   `json:"success"`
		Data    struct {
			Items []struct {
				FileID   string `json:"fileId"`
				Name     string `json:"name"`
				Type     string `json:"type"`
				Category string `json:"category"`
				Size     int64  `json:"size"`
				UpdateAt string `json:"updatedAt"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}

	var files []core.FileInfo
	for _, item := range res.Data.Items {
		isFolder := item.Type == "folder" || item.Category == "folder"
		updateTime, _ := time.Parse("2006-01-02 15:04:05", item.UpdateAt)
		files = append(files, core.FileInfo{
			ID:         item.FileID,
			Name:       item.Name,
			Path:       item.FileID, // 139 ID 即可代表 Path
			IsFolder:   isFolder,
			Size:       item.Size,
			UpdatedAt:  item.UpdateAt,
			UpdateTime: updateTime,
		})
	}
	return files, nil
}

func (c *Cloud139) CreateFolder(ctx context.Context, parentID, name string) (*core.FileInfo, error) {
	if parentID == "" {
		parentID = "/"
	}
	sign := c.computeMcloudSign(parentID)
	headers := map[string]string{
		"mcloud-sign": sign,
	}
	body := map[string]interface{}{
		"parentFileId": parentID,
		"name":         name,
		"type":         "folder",
	}
	resp, err := c.doRequest(ctx, "POST", PersonalKdNjsURL+"/hcy/file/create", body, headers)
	if err != nil {
		return nil, err
	}
	var res struct {
		Code    string `json:"code"`
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    struct {
			FileID   string `json:"fileId"`
			ID       string `json:"id"`
			FileName string `json:"fileName"`
		} `json:"data"`
	}
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}
	if !res.Success && res.Code != "0" && res.Code != "0000" && res.Code != "" {
		return nil, fmt.Errorf("139 CreateFolder error [%s]: %s", res.Code, res.Message)
	}

	finalID := res.Data.FileID
	if finalID == "" {
		finalID = res.Data.ID
	}

	return &core.FileInfo{
		ID:       finalID,
		Name:     name,
		Path:     finalID,
		IsFolder: true,
	}, nil
}

func (c *Cloud139) ParseShare(ctx context.Context, shareURL, extractCode string) ([]core.FileInfo, error) {
	linkID, passwd, err := c.parseShareLink(shareURL)
	if err != nil {
		return nil, err
	}
	if extractCode != "" {
		passwd = extractCode
	}

	info, err := c.getShareInfo(ctx, linkID, passwd, "root")
	if err != nil {
		return nil, err
	}

	cst := time.FixedZone("CST", 8*3600)

	var files []core.FileInfo
	// 139 分享信息中 caLst 是文件夹，coLst 是文件
	if caLst, ok := info["caLst"].([]interface{}); ok {
		for _, item := range caLst {
			if f, ok := item.(map[string]interface{}); ok {
				updateAt, _ := f["updateTime"].(string)
				updateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", updateAt, cst)
				
				// 优先使用官方返回的 path 字段，格式通常为 parentID/caID
				path, _ := f["path"].(string)
				if path == "" {
					caID := f["catalogID"]
					if caID == nil {
						caID = f["caID"]
					}
					path = fmt.Sprintf("%v/%v", f["parentCatalogID"], caID)
				}

				files = append(files, core.FileInfo{
					ID:         path,
					Name:       fmt.Sprintf("%v", f["catalogName"]),
					IsFolder:   true,
					UpdatedAt:  updateAt,
					UpdateTime: updateTime,
				})
			}
		}
	}
	if coLst, ok := info["coLst"].([]interface{}); ok {
		for _, item := range coLst {
			if f, ok := item.(map[string]interface{}); ok {
				updateAt, _ := f["updateTime"].(string)
				updateTime, _ := time.ParseInLocation("2006-01-02 15:04:05", updateAt, cst)
				size, _ := f["contentSize"].(float64)

				// 优先使用官方返回的 path 字段，格式通常为 parentID/coID
				path, _ := f["path"].(string)
				if path == "" {
					path = fmt.Sprintf("%v/%v", f["parentCatalogID"], f["contentID"])
				}

				files = append(files, core.FileInfo{
					ID:         path,
					Name:       fmt.Sprintf("%v", f["contentName"]),
					IsFolder:   false,
					Size:       int64(size),
					UpdatedAt:  updateAt,
					UpdateTime: updateTime,
				})
			}
		}
	}
	return files, nil
}

func (c *Cloud139) SaveFileTo(ctx context.Context, fileID, targetPath string) error {
	// 139 的 fileID 在分享场景下是 "parentID/fileID" 格式
	parts := strings.Split(fileID, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid fileID format for 139 share: %s", fileID)
	}

	// 提前解析出 linkID 和 passwd (由于 SaveFileTo 接口没传 shareURL，
	// 139 的这个接口设计可能需要驱动内部状态或者在调用处做特殊处理。
	// 这里为了符合接口，我们假设 fileID 中已经隐含了上下文，
	// 或者在实际 Worker 中，如果是 139，我们会调用 SaveLink 这种批量接口。)
	// 考虑到 139 的特殊性，我们在 SaveLink 中实现具体逻辑。
	return fmt.Errorf("139 driver prefers batch SaveLink operation")
}

func (c *Cloud139) SaveLink(ctx context.Context, shareURL, extractCode, targetPath string, fileIDs []string) error {
	phone := c.getPhone()
	if phone == "" {
		return fmt.Errorf("139 SaveLink error: 无法获取合法的 11 位手机号，请先在账号管理中校验账号")
	}

	linkID, passwd, err := c.parseShareLink(shareURL)
	if err != nil {
		return err
	}
	if extractCode != "" {
		passwd = extractCode
	}

	info, err := c.getShareInfo(ctx, linkID, passwd, "root")
	if err != nil {
		return err
	}

	targetID, err := c.prepareTargetPath(ctx, targetPath)
	if err != nil {
		return err
	}

	// 139 Orchestration API 根目录 ID 必须为 root
	if targetID == "/" || targetID == "" {
		targetID = "root"
	}

	idMap := make(map[string]bool)
	for _, id := range fileIDs {
		idMap[id] = true
	}

	coPathLst := []string{}
	if coLst, ok := info["coLst"].([]interface{}); ok {
		for _, item := range coLst {
			if f, ok := item.(map[string]interface{}); ok {
				path, _ := f["path"].(string)
				if path == "" {
					path = fmt.Sprintf("%v/%v", f["parentCatalogID"], f["contentID"])
				}
				if len(fileIDs) == 0 || idMap[path] {
					coPathLst = append(coPathLst, path)
				}
			}
		}
	}

	caPathLst := []string{}
	if caLst, ok := info["caLst"].([]interface{}); ok {
		for _, item := range caLst {
			if f, ok := item.(map[string]interface{}); ok {
				path, _ := f["path"].(string)
				if path == "" {
					caID := f["catalogID"]
					if caID == nil {
						caID = f["caID"]
					}
					path = fmt.Sprintf("%v/%v", f["parentCatalogID"], caID)
				}
				if len(fileIDs) == 0 || idMap[path] {
					caPathLst = append(caPathLst, path)
				}
			}
		}
	}

	if len(coPathLst) == 0 && len(caPathLst) == 0 {
		return nil
	}

	saveBody := map[string]interface{}{
		"createOuterLinkBatchOprTaskReq": map[string]interface{}{
			"msisdn":       phone,
			"ownerAccount": "",
			"taskType":     1,
			"linkID":       linkID,
			"needPassword": passwd != "",
			"taskInfo": map[string]interface{}{
				"linkID":          linkID,
				"needPassword":    passwd != "",
				"contentInfoList": coPathLst,
				"catalogInfoList": caPathLst,
				"newCatalogID":    targetID,
			},
		},
	}

	resp, err := c.doRequest(ctx, "POST", ShareKdNjsURL+"/yun-share/richlifeApp/devapp/IBatchOprTask/createOuterLinkBatchOprTask", saveBody, nil)
	if err != nil {
		return err
	}

	var saveRes struct {
		Code interface{} `json:"code"`
		Desc string      `json:"desc"`
	}
	if err := json.Unmarshal(resp, &saveRes); err == nil {
		codeStr := fmt.Sprintf("%v", saveRes.Code)
		if codeStr != "0" && codeStr != "0000" && codeStr != "" {
			return fmt.Errorf("139 SaveLink error [%s]: %s", codeStr, saveRes.Desc)
		}
	}

	return nil
}

func (c *Cloud139) DeleteFile(ctx context.Context, fileID string) error {
	sign := c.computeMcloudSign("/")
	headers := map[string]string{
		"mcloud-sign":            sign,
		"INNER-HCY-ROUTER-HTTPS": "1",
	}
	body := map[string]interface{}{
		"fileIds": []string{fileID},
	}
	_, err := c.doRequest(ctx, "POST", PersonalKdNjsURL+"/hcy/recyclebin/batchTrash", body, headers)
	return err
}

func (c *Cloud139) parseShareLink(input string) (string, string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", "", fmt.Errorf("empty share link")
	}

	// 1. 尝试作为纯 linkID 处理 (无 "/" 且长度在 4-32 之间)
	reBare := regexp.MustCompile(`^[a-zA-Z0-9_-]{4,32}$`)
	if reBare.MatchString(trimmed) && !strings.Contains(trimmed, "/") {
		return trimmed, "", nil
	}

	// 2. 补全协议并解析 URL
	urlStr := trimmed
	if !strings.HasPrefix(strings.ToLower(urlStr), "http") {
		urlStr = "https://" + urlStr
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse url: %v", err)
	}

	// 3. 提取常规 Query 参数
	linkID := u.Query().Get("linkID")
	if linkID == "" {
		linkID = u.Query().Get("linkId")
	}

	// 4. 处理 caiyun.139.com/m/i?{linkID} 格式 (linkID 作为 key)
	if linkID == "" && u.RawQuery != "" && !strings.Contains(u.RawQuery, "=") {
		if reBare.MatchString(u.RawQuery) {
			linkID = u.RawQuery
		}
	}

	// 5. 处理 Hash 路径格式: #/w/i/{linkID}
	if linkID == "" && u.Fragment != "" {
		// 尝试 Fragment 中的 query
		if strings.Contains(u.Fragment, "?") {
			parts := strings.Split(u.Fragment, "?")
			if len(parts) > 1 {
				fQuery, _ := url.ParseQuery(parts[1])
				linkID = fQuery.Get("linkID")
				if linkID == "" {
					linkID = fQuery.Get("linkId")
				}
			}
		}
		// 尝试 Fragment 路径末段
		if linkID == "" {
			fPath := u.Fragment
			if strings.Contains(fPath, "?") {
				fPath = strings.Split(fPath, "?")[0]
			}
			parts := strings.Split(strings.Trim(fPath, "/"), "/")
			if len(parts) > 0 {
				candidate := parts[len(parts)-1]
				if reBare.MatchString(candidate) {
					linkID = candidate
				}
			}
		}
	}

	// 6. 兜底处理 Path 末段
	if linkID == "" {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 {
			candidate := parts[len(parts)-1]
			if reBare.MatchString(candidate) {
				linkID = candidate
			}
		}
	}

	if linkID == "" {
		return "", "", fmt.Errorf("linkID not found in: %s", input)
	}

	passwd := u.Query().Get("passwd")
	if passwd == "" {
		passwd = u.Query().Get("pwd")
	}

	return linkID, passwd, nil
}

func (c *Cloud139) getShareInfo(ctx context.Context, linkID, passwd, pCaID string) (map[string]interface{}, error) {
	headers := map[string]string{
		"caller":         "web",
		"x-m4c-caller":   "PC",
		"mcloud-client":  "10701",
		"mcloud-version": "7.17.2",
		"mcloud-channel": "1000101",
	}
	body := map[string]interface{}{
		"getOutLinkInfoReq": map[string]interface{}{
			"account": c.getPhone(),
			"linkID":  linkID,
			"passwd":  passwd,
			"pCaID":   pCaID,
			"caSrt":   0,
			"coSrt":   0,
			"srtDr":   1,
			"bNum":    1,
			"eNum":    200,
		},
	}
	resp, err := c.doRequest(ctx, "POST", ShareKdNjsURL+"/yun-share/richlifeApp/devapp/IOutLink/getOutLinkInfoV6", body, headers)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}

	// 检查 code 状态
	if code, ok := res["code"]; ok {
		codeStr := fmt.Sprintf("%v", code)
		if codeStr != "0" && codeStr != "0000" && codeStr != "" {
			return nil, fmt.Errorf("139 getShareInfo error [%s]: %v", codeStr, res["message"])
		}
	}

	// 优先返回 data 字段，否则返回整个 map
	if data, ok := res["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return res, nil
}

func (c *Cloud139) prepareTargetPath(ctx context.Context, path string) (string, error) {
	if path == "" || path == "/" {
		return "root", nil
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentID := "root"
	for _, part := range parts {
		files, err := c.ListFiles(ctx, currentID)
		if err != nil {
			return "", err
		}
		found := false
		for _, f := range files {
			if f.IsFolder && f.Name == part {
				currentID = f.ID
				found = true
				break
			}
		}
		if !found {
			newFolder, err := c.CreateFolder(ctx, currentID, part)
			if err != nil {
				return "", err
			}
			currentID = newFolder.ID
		}
	}
	return currentID, nil
}
