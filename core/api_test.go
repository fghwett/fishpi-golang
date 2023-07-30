package core

import (
	"fishpi/config"
	"fishpi/logger"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestApi(t *testing.T) {
	// 初始化日志程序
	loger := logger.NewConsoleLogger()

	// 读取配置文件
	conf, err := config.NewConfig(`../_tmp/config.yaml`)
	if err != nil {
		loger.Logf("读取配置文件失败 \n错误信息：%s", err)
		return
	}

	// 初始化FishPi API
	var api *Api
	if api, err = NewApi(conf.FishPi.ApiBase); err != nil {
		loger.Logf("FishPi地址信息填写失败 %s", err)
		return
	}

	fishPiSdk := NewSdk(api, conf.FishPi.ApiBase, conf.FishPi.ApiKey, conf.FishPi.Username, loger)

	//body, ee := fishPiSdk.PointTransfer("9811", 1, "(测试) [你在用什么笔记呢](https://fishpi.cn/article/1670463550914)有奖问答")
	//if ee != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(string(body))
	//
	//return
	reply, e := fishPiSdk.GetArticleInfo(&ArticleInfoData{
		ArticleId: `1670463550914`,
		Page:      1,
	})
	if e != nil {
		t.Fatal(e)
	}
	//t.Log(os.WriteFile(`../_tmp/article_info.json`, body, os.ModePerm))
	var cs comments

	cs.parse(reply)

	maxPage := reply.Data.Pagination.PaginationPageCount

	if maxPage > 1 {
		for i := 2; i <= maxPage; i++ {
			reply, e = fishPiSdk.GetArticleInfo(&ArticleInfoData{
				ArticleId: `1670463550914`,
				Page:      i,
			})
			if e != nil {
				t.Fatal(e)
			}
			cs.parse(reply)
		}
	}

	boolMap := make(map[string]bool)
	for _, c := range cs {
		if boolMap[c.PUsername] {
			continue
		}

		//if body, eee := fishPiSdk.PointTransfer(c.PUsername, 100, fmt.Sprintf("感谢你在 [你在用什么笔记呢](https://fishpi.cn%s) 的评论，获得100积分奖励", c.PSource)); eee != nil {
		//	t.Fatalf("%s : %s", c.PUsername, eee)
		//} else {
		//	t.Logf("success %s %s", c.PUsername, string(body))
		//	c.Send = true
		//}

		boolMap[c.PUsername] = true
	}

	t.Log(os.WriteFile(`../_tmp/comments.txt`, []byte(strings.Join(cs.String(), "\n\n")), os.ModePerm))

}

type comments []comment

func (c *comments) parse(reply *ArticleInfoReply) {
	for _, v := range reply.Data.Article.ArticleComments {
		if v.CommentOriginalCommentId != "" {
			continue
			has := false
			for i, v2 := range *c {
				if v2.OId == v.CommentOriginalCommentId {
					(*c)[i].Children = append((*c)[i].Children, comment{
						OId:       v.OId,
						POId:      v.CommentAuthorId,
						PStatus:   v.Commenter.UserStatus,
						PUsername: v.Commenter.UserName,
						PContent:  v.CommentContent,
						PTime:     v.CommentCreateTime,
					})
					has = true
					break
				}
			}
			if !has {
				*c = append(*c, comment{
					OId: v.CommentOriginalCommentId,
					Children: []comment{
						{
							OId:       v.OId,
							POId:      v.CommentAuthorId,
							PStatus:   v.Commenter.UserStatus,
							PUsername: v.Commenter.UserName,
							PContent:  v.CommentContent,
							PTime:     v.CommentCreateTime,
						},
					},
				})
			}
			continue
		}
		if v.Commenter.UserStatus == 4 {
			continue
		}

		if v.Commenter.UserName == "8888" {
			continue
		}
		//// 当c中有oid时更新数据
		//for i, v2 := range *c {
		//	if v2.OId == v.OId {
		//		(*c)[i].POId = v.CommentAuthorId
		//		(*c)[i].PStatus = v.Commenter.UserStatus
		//		(*c)[i].PUsername = v.Commenter.UserName
		//		(*c)[i].PContent = v.CommentContent
		//		(*c)[i].PTime = v.CommentCreateTime
		//		continue
		//	}
		//}
		*c = append(*c, comment{
			OId:       v.OId,
			POId:      v.CommentAuthorId,
			PStatus:   v.Commenter.UserStatus,
			PUsername: v.Commenter.UserName,
			PContent:  v.CommentContent,
			PSource:   v.CommentSharpURL,
			PTime:     v.CommentCreateTime,
		})
	}
}

func (c *comments) String() []string {
	var arr []string
	for i, v := range *c {
		arr = append(arr, fmt.Sprintf("%d: %s", i+1, strings.Join(v.String(), "\n\t")))
	}
	return arr
}

type comment struct {
	OId       string    `json:"o_id"`
	POId      string    `json:"po_id"`
	PStatus   int       `json:"p_status"`
	PUsername string    `json:"p_username"`
	PContent  string    `json:"p_content"`
	PSource   string    `json:"p_source"`
	PTime     string    `json:"p_time"`
	Send      bool      `json:"send"`
	Children  []comment `json:"children"`
}

func (c *comment) String() []string {
	var arr []string
	arr = append(arr, fmt.Sprintf("OId: %s", c.OId))
	arr = append(arr, fmt.Sprintf("POId: %s", c.POId))
	arr = append(arr, fmt.Sprintf("PStatus: %d", c.PStatus))
	arr = append(arr, fmt.Sprintf("PUsername: %s", c.PUsername))
	arr = append(arr, fmt.Sprintf("PContent: %s", c.PContent))
	arr = append(arr, fmt.Sprintf("PSource: %s", c.PSource))
	arr = append(arr, fmt.Sprintf("Send: %v", c.Send))
	arr = append(arr, fmt.Sprintf("PTime: %s", c.PTime))
	if c.Children != nil {
		children := comments(c.Children)
		arr = append(arr, fmt.Sprintf("\t Children: %s", strings.Join(children.String(), "\n\t")))
	}
	return arr
}
