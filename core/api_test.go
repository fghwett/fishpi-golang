package core

import (
	"fishpi/config"
	"fishpi/logger"
	"os"
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
	body, e := fishPiSdk.GetArticleInfo(`1670463550914`)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(os.WriteFile(`../_tmp/article_info.json`, body, os.ModePerm))
}

type T struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Article struct {
			ArticleCreateTime  string `json:"articleCreateTime"`
			DiscussionViewable bool   `json:"discussionViewable"`
			ArticleToC         string `json:"articleToC"`
			ThankedCnt         int    `json:"thankedCnt"`
			ArticleComments    []struct {
				CommentNice              bool    `json:"commentNice"`
				CommentCreateTimeStr     string  `json:"commentCreateTimeStr"`
				CommentAuthorId          string  `json:"commentAuthorId"`
				CommentScore             float64 `json:"commentScore"`
				CommentCreateTime        string  `json:"commentCreateTime"`
				CommentAuthorURL         string  `json:"commentAuthorURL"`
				CommentVote              int     `json:"commentVote"`
				CommentRevisionCount     int     `json:"commentRevisionCount"`
				TimeAgo                  string  `json:"timeAgo"`
				CommentOriginalCommentId string  `json:"commentOriginalCommentId"`
				SysMetal                 []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Data        string `json:"data"`
					Attr        string `json:"attr"`
					Enabled     bool   `json:"enabled"`
				} `json:"sysMetal"`
				CommentGoodCnt     int    `json:"commentGoodCnt"`
				CommentVisible     int    `json:"commentVisible"`
				CommentOnArticleId string `json:"commentOnArticleId"`
				RewardedCnt        int    `json:"rewardedCnt"`
				CommentSharpURL    string `json:"commentSharpURL"`
				CommentAnonymous   int    `json:"commentAnonymous"`
				CommentReplyCnt    int    `json:"commentReplyCnt"`
				OId                string `json:"oId"`
				CommentContent     string `json:"commentContent"`
				CommentStatus      int    `json:"commentStatus"`
				Commenter          struct {
					UserOnlineFlag                bool   `json:"userOnlineFlag"`
					OnlineMinute                  int    `json:"onlineMinute"`
					UserPointStatus               int    `json:"userPointStatus"`
					UserFollowerStatus            int    `json:"userFollowerStatus"`
					UserGuideStep                 int    `json:"userGuideStep"`
					UserOnlineStatus              int    `json:"userOnlineStatus"`
					UserCurrentCheckinStreakStart int    `json:"userCurrentCheckinStreakStart"`
					ChatRoomPictureStatus         int    `json:"chatRoomPictureStatus"`
					UserTags                      string `json:"userTags"`
					UserCommentStatus             int    `json:"userCommentStatus"`
					UserTimezone                  string `json:"userTimezone"`
					UserURL                       string `json:"userURL"`
					UserForwardPageStatus         int    `json:"userForwardPageStatus"`
					UserUAStatus                  int    `json:"userUAStatus"`
					UserIndexRedirectURL          string `json:"userIndexRedirectURL"`
					UserLatestArticleTime         int64  `json:"userLatestArticleTime"`
					UserTagCount                  int    `json:"userTagCount"`
					UserNickname                  string `json:"userNickname"`
					UserListViewMode              int    `json:"userListViewMode"`
					UserLongestCheckinStreak      int    `json:"userLongestCheckinStreak"`
					UserAvatarType                int    `json:"userAvatarType"`
					UserSubMailSendTime           int64  `json:"userSubMailSendTime"`
					UserUpdateTime                int64  `json:"userUpdateTime"`
					UserSubMailStatus             int    `json:"userSubMailStatus"`
					UserJoinPointRank             int    `json:"userJoinPointRank"`
					UserLatestLoginTime           int64  `json:"userLatestLoginTime"`
					UserAppRole                   int    `json:"userAppRole"`
					UserAvatarViewMode            int    `json:"userAvatarViewMode"`
					UserStatus                    int    `json:"userStatus"`
					UserLongestCheckinStreakEnd   int    `json:"userLongestCheckinStreakEnd"`
					UserWatchingArticleStatus     int    `json:"userWatchingArticleStatus"`
					UserLatestCmtTime             int64  `json:"userLatestCmtTime"`
					UserProvince                  string `json:"userProvince"`
					UserCurrentCheckinStreak      int    `json:"userCurrentCheckinStreak"`
					UserNo                        int    `json:"userNo"`
					UserAvatarURL                 string `json:"userAvatarURL"`
					UserFollowingTagStatus        int    `json:"userFollowingTagStatus"`
					UserLanguage                  string `json:"userLanguage"`
					UserJoinUsedPointRank         int    `json:"userJoinUsedPointRank"`
					UserCurrentCheckinStreakEnd   int    `json:"userCurrentCheckinStreakEnd"`
					UserFollowingArticleStatus    int    `json:"userFollowingArticleStatus"`
					UserKeyboardShortcutsStatus   int    `json:"userKeyboardShortcutsStatus"`
					UserReplyWatchArticleStatus   int    `json:"userReplyWatchArticleStatus"`
					UserCommentViewMode           int    `json:"userCommentViewMode"`
					UserBreezemoonStatus          int    `json:"userBreezemoonStatus"`
					UserCheckinTime               int64  `json:"userCheckinTime"`
					UserUsedPoint                 int    `json:"userUsedPoint"`
					UserArticleStatus             int    `json:"userArticleStatus"`
					UserPoint                     int    `json:"userPoint"`
					UserCommentCount              int    `json:"userCommentCount"`
					UserIntro                     string `json:"userIntro"`
					UserMobileSkin                string `json:"userMobileSkin"`
					UserListPageSize              int    `json:"userListPageSize"`
					OId                           string `json:"oId"`
					UserName                      string `json:"userName"`
					UserGeoStatus                 int    `json:"userGeoStatus"`
					UserLongestCheckinStreakStart int    `json:"userLongestCheckinStreakStart"`
					UserSkin                      string `json:"userSkin"`
					UserNotifyStatus              int    `json:"userNotifyStatus"`
					UserFollowingUserStatus       int    `json:"userFollowingUserStatus"`
					UserArticleCount              int    `json:"userArticleCount"`
					UserRole                      string `json:"userRole"`
				} `json:"commenter"`
				CommentAuthorName                 string `json:"commentAuthorName"`
				CommentThankCnt                   int    `json:"commentThankCnt"`
				CommentBadCnt                     int    `json:"commentBadCnt"`
				Rewarded                          bool   `json:"rewarded"`
				CommentAuthorThumbnailURL         string `json:"commentAuthorThumbnailURL"`
				CommentAudioURL                   string `json:"commentAudioURL"`
				CommentQnAOffered                 int    `json:"commentQnAOffered"`
				CommentOriginalAuthorThumbnailURL string `json:"commentOriginalAuthorThumbnailURL,omitempty"`
				PaginationCurrentPageNum          int    `json:"paginationCurrentPageNum,omitempty"`
			} `json:"articleComments"`
			ArticleRewardPoint          int    `json:"articleRewardPoint"`
			ArticleRevisionCount        int    `json:"articleRevisionCount"`
			ArticleLatestCmtTime        string `json:"articleLatestCmtTime"`
			ArticleThumbnailURL         string `json:"articleThumbnailURL"`
			ArticleAuthorName           string `json:"articleAuthorName"`
			ArticleType                 int    `json:"articleType"`
			ArticleCreateTimeStr        string `json:"articleCreateTimeStr"`
			ArticleViewCount            int    `json:"articleViewCount"`
			ArticleCommentable          bool   `json:"articleCommentable"`
			ArticleAuthorThumbnailURL20 string `json:"articleAuthorThumbnailURL20"`
			ArticleOriginalContent      string `json:"articleOriginalContent"`
			ArticlePreviewContent       string `json:"articlePreviewContent"`
			ArticleContent              string `json:"articleContent"`
			ArticleAuthorIntro          string `json:"articleAuthorIntro"`
			ArticleCommentCount         int    `json:"articleCommentCount"`
			RewardedCnt                 int    `json:"rewardedCnt"`
			ArticleLatestCmterName      string `json:"articleLatestCmterName"`
			ArticleAnonymousView        int    `json:"articleAnonymousView"`
			CmtTimeAgo                  string `json:"cmtTimeAgo"`
			ArticleLatestCmtTimeStr     string `json:"articleLatestCmtTimeStr"`
			ArticleNiceComments         []struct {
				CommentCreateTimeStr     string  `json:"commentCreateTimeStr"`
				CommentAuthorId          string  `json:"commentAuthorId"`
				CommentScore             float64 `json:"commentScore"`
				CommentCreateTime        string  `json:"commentCreateTime"`
				CommentAuthorURL         string  `json:"commentAuthorURL"`
				CommentVote              int     `json:"commentVote"`
				TimeAgo                  string  `json:"timeAgo"`
				CommentOriginalCommentId string  `json:"commentOriginalCommentId"`
				SysMetal                 []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Data        string `json:"data"`
					Attr        string `json:"attr"`
					Enabled     bool   `json:"enabled"`
				} `json:"sysMetal"`
				CommentGoodCnt     int    `json:"commentGoodCnt"`
				CommentVisible     int    `json:"commentVisible"`
				CommentOnArticleId string `json:"commentOnArticleId"`
				RewardedCnt        int    `json:"rewardedCnt"`
				CommentThankLabel  string `json:"commentThankLabel"`
				CommentSharpURL    string `json:"commentSharpURL"`
				CommentAnonymous   int    `json:"commentAnonymous"`
				CommentReplyCnt    int    `json:"commentReplyCnt"`
				OId                string `json:"oId"`
				CommentContent     string `json:"commentContent"`
				CommentStatus      int    `json:"commentStatus"`
				Commenter          struct {
					UserOnlineFlag                bool   `json:"userOnlineFlag"`
					OnlineMinute                  int    `json:"onlineMinute"`
					UserPointStatus               int    `json:"userPointStatus"`
					UserFollowerStatus            int    `json:"userFollowerStatus"`
					UserGuideStep                 int    `json:"userGuideStep"`
					UserOnlineStatus              int    `json:"userOnlineStatus"`
					UserCurrentCheckinStreakStart int    `json:"userCurrentCheckinStreakStart"`
					ChatRoomPictureStatus         int    `json:"chatRoomPictureStatus"`
					UserTags                      string `json:"userTags"`
					UserCommentStatus             int    `json:"userCommentStatus"`
					UserTimezone                  string `json:"userTimezone"`
					UserURL                       string `json:"userURL"`
					UserForwardPageStatus         int    `json:"userForwardPageStatus"`
					UserUAStatus                  int    `json:"userUAStatus"`
					UserIndexRedirectURL          string `json:"userIndexRedirectURL"`
					UserLatestArticleTime         int64  `json:"userLatestArticleTime"`
					UserTagCount                  int    `json:"userTagCount"`
					UserNickname                  string `json:"userNickname"`
					UserListViewMode              int    `json:"userListViewMode"`
					UserLongestCheckinStreak      int    `json:"userLongestCheckinStreak"`
					UserAvatarType                int    `json:"userAvatarType"`
					UserSubMailSendTime           int64  `json:"userSubMailSendTime"`
					UserUpdateTime                int64  `json:"userUpdateTime"`
					UserSubMailStatus             int    `json:"userSubMailStatus"`
					UserJoinPointRank             int    `json:"userJoinPointRank"`
					UserLatestLoginTime           int64  `json:"userLatestLoginTime"`
					UserAppRole                   int    `json:"userAppRole"`
					UserAvatarViewMode            int    `json:"userAvatarViewMode"`
					UserStatus                    int    `json:"userStatus"`
					UserLongestCheckinStreakEnd   int    `json:"userLongestCheckinStreakEnd"`
					UserWatchingArticleStatus     int    `json:"userWatchingArticleStatus"`
					UserLatestCmtTime             int64  `json:"userLatestCmtTime"`
					UserProvince                  string `json:"userProvince"`
					UserCurrentCheckinStreak      int    `json:"userCurrentCheckinStreak"`
					UserNo                        int    `json:"userNo"`
					UserAvatarURL                 string `json:"userAvatarURL"`
					UserFollowingTagStatus        int    `json:"userFollowingTagStatus"`
					UserLanguage                  string `json:"userLanguage"`
					UserJoinUsedPointRank         int    `json:"userJoinUsedPointRank"`
					UserCurrentCheckinStreakEnd   int    `json:"userCurrentCheckinStreakEnd"`
					UserFollowingArticleStatus    int    `json:"userFollowingArticleStatus"`
					UserKeyboardShortcutsStatus   int    `json:"userKeyboardShortcutsStatus"`
					UserReplyWatchArticleStatus   int    `json:"userReplyWatchArticleStatus"`
					UserCommentViewMode           int    `json:"userCommentViewMode"`
					UserBreezemoonStatus          int    `json:"userBreezemoonStatus"`
					UserCheckinTime               int64  `json:"userCheckinTime"`
					UserUsedPoint                 int    `json:"userUsedPoint"`
					UserArticleStatus             int    `json:"userArticleStatus"`
					UserPoint                     int    `json:"userPoint"`
					UserCommentCount              int    `json:"userCommentCount"`
					UserIntro                     string `json:"userIntro"`
					UserMobileSkin                string `json:"userMobileSkin"`
					UserListPageSize              int    `json:"userListPageSize"`
					OId                           string `json:"oId"`
					UserName                      string `json:"userName"`
					UserGeoStatus                 int    `json:"userGeoStatus"`
					UserLongestCheckinStreakStart int    `json:"userLongestCheckinStreakStart"`
					UserSkin                      string `json:"userSkin"`
					UserNotifyStatus              int    `json:"userNotifyStatus"`
					UserFollowingUserStatus       int    `json:"userFollowingUserStatus"`
					UserArticleCount              int    `json:"userArticleCount"`
					UserRole                      string `json:"userRole"`
				} `json:"commenter"`
				PaginationCurrentPageNum  int    `json:"paginationCurrentPageNum"`
				CommentAuthorName         string `json:"commentAuthorName"`
				CommentThankCnt           int    `json:"commentThankCnt"`
				CommentBadCnt             int    `json:"commentBadCnt"`
				Rewarded                  bool   `json:"rewarded"`
				CommentAuthorThumbnailURL string `json:"commentAuthorThumbnailURL"`
				CommentAudioURL           string `json:"commentAudioURL"`
				CommentQnAOffered         int    `json:"commentQnAOffered"`
			} `json:"articleNiceComments"`
			Rewarded                     bool    `json:"rewarded"`
			ArticleHeat                  int     `json:"articleHeat"`
			ArticlePerfect               int     `json:"articlePerfect"`
			ArticleAuthorThumbnailURL210 string  `json:"articleAuthorThumbnailURL210"`
			ArticlePermalink             string  `json:"articlePermalink"`
			ArticleCity                  string  `json:"articleCity"`
			ArticleShowInList            int     `json:"articleShowInList"`
			IsMyArticle                  bool    `json:"isMyArticle"`
			ArticleIP                    string  `json:"articleIP"`
			ArticleEditorType            int     `json:"articleEditorType"`
			ArticleVote                  int     `json:"articleVote"`
			ArticleRandomDouble          float64 `json:"articleRandomDouble"`
			ArticleAuthorId              string  `json:"articleAuthorId"`
			ArticleBadCnt                int     `json:"articleBadCnt"`
			ArticleAuthorURL             string  `json:"articleAuthorURL"`
			IsWatching                   bool    `json:"isWatching"`
			ArticleGoodCnt               int     `json:"articleGoodCnt"`
			ArticleQnAOfferPoint         int     `json:"articleQnAOfferPoint"`
			ArticleStickRemains          int     `json:"articleStickRemains"`
			TimeAgo                      string  `json:"timeAgo"`
			ArticleUpdateTimeStr         string  `json:"articleUpdateTimeStr"`
			Offered                      bool    `json:"offered"`
			ArticleWatchCnt              int     `json:"articleWatchCnt"`
			ArticleTitleEmoj             string  `json:"articleTitleEmoj"`
			ArticleTitleEmojUnicode      string  `json:"articleTitleEmojUnicode"`
			ArticleAudioURL              string  `json:"articleAudioURL"`
			ArticleAuthorThumbnailURL48  string  `json:"articleAuthorThumbnailURL48"`
			Thanked                      bool    `json:"thanked"`
			ArticleImg1URL               string  `json:"articleImg1URL"`
			ArticlePushOrder             int     `json:"articlePushOrder"`
			ArticleCollectCnt            int     `json:"articleCollectCnt"`
			ArticleTitle                 string  `json:"articleTitle"`
			IsFollowing                  bool    `json:"isFollowing"`
			ArticleTags                  string  `json:"articleTags"`
			OId                          string  `json:"oId"`
			ArticleStick                 int     `json:"articleStick"`
			ArticleTagObjs               []struct {
				TagShowSideAd     int     `json:"tagShowSideAd"`
				TagIconPath       string  `json:"tagIconPath"`
				TagStatus         int     `json:"tagStatus"`
				TagBadCnt         int     `json:"tagBadCnt"`
				TagRandomDouble   float64 `json:"tagRandomDouble"`
				TagTitle          string  `json:"tagTitle"`
				OId               string  `json:"oId"`
				TagURI            string  `json:"tagURI"`
				TagAd             string  `json:"tagAd"`
				TagGoodCnt        int     `json:"tagGoodCnt"`
				TagCSS            string  `json:"tagCSS"`
				TagCommentCount   int     `json:"tagCommentCount"`
				TagFollowerCount  int     `json:"tagFollowerCount"`
				TagSeoTitle       string  `json:"tagSeoTitle"`
				TagLinkCount      int     `json:"tagLinkCount"`
				TagSeoDesc        string  `json:"tagSeoDesc"`
				TagReferenceCount int     `json:"tagReferenceCount"`
				TagSeoKeywords    string  `json:"tagSeoKeywords"`
				TagDescription    string  `json:"tagDescription"`
			} `json:"articleTagObjs"`
			ArticleAnonymous     int    `json:"articleAnonymous"`
			ArticleThankCnt      int    `json:"articleThankCnt"`
			ArticleRewardContent string `json:"articleRewardContent"`
			RedditScore          int    `json:"redditScore"`
			ArticleUpdateTime    string `json:"articleUpdateTime"`
			ArticleStatus        int    `json:"articleStatus"`
			ArticleAuthor        struct {
				UserOnlineFlag                bool   `json:"userOnlineFlag"`
				OnlineMinute                  int    `json:"onlineMinute"`
				UserPointStatus               int    `json:"userPointStatus"`
				UserFollowerStatus            int    `json:"userFollowerStatus"`
				UserGuideStep                 int    `json:"userGuideStep"`
				UserOnlineStatus              int    `json:"userOnlineStatus"`
				UserCurrentCheckinStreakStart int    `json:"userCurrentCheckinStreakStart"`
				ChatRoomPictureStatus         int    `json:"chatRoomPictureStatus"`
				UserTags                      string `json:"userTags"`
				SysMetal                      []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Data        string `json:"data"`
					Attr        string `json:"attr"`
					Enabled     bool   `json:"enabled"`
				} `json:"sysMetal"`
				UserCommentStatus             int    `json:"userCommentStatus"`
				UserTimezone                  string `json:"userTimezone"`
				UserURL                       string `json:"userURL"`
				UserForwardPageStatus         int    `json:"userForwardPageStatus"`
				UserUAStatus                  int    `json:"userUAStatus"`
				UserIndexRedirectURL          string `json:"userIndexRedirectURL"`
				UserLatestArticleTime         int64  `json:"userLatestArticleTime"`
				UserTagCount                  int    `json:"userTagCount"`
				UserNickname                  string `json:"userNickname"`
				UserListViewMode              int    `json:"userListViewMode"`
				UserLongestCheckinStreak      int    `json:"userLongestCheckinStreak"`
				UserAvatarType                int    `json:"userAvatarType"`
				UserSubMailSendTime           int64  `json:"userSubMailSendTime"`
				UserUpdateTime                int64  `json:"userUpdateTime"`
				UserSubMailStatus             int    `json:"userSubMailStatus"`
				UserJoinPointRank             int    `json:"userJoinPointRank"`
				UserLatestLoginTime           int64  `json:"userLatestLoginTime"`
				UserAppRole                   int    `json:"userAppRole"`
				UserAvatarViewMode            int    `json:"userAvatarViewMode"`
				UserStatus                    int    `json:"userStatus"`
				UserLongestCheckinStreakEnd   int    `json:"userLongestCheckinStreakEnd"`
				UserWatchingArticleStatus     int    `json:"userWatchingArticleStatus"`
				UserLatestCmtTime             int64  `json:"userLatestCmtTime"`
				UserProvince                  string `json:"userProvince"`
				UserCurrentCheckinStreak      int    `json:"userCurrentCheckinStreak"`
				UserNo                        int    `json:"userNo"`
				UserAvatarURL                 string `json:"userAvatarURL"`
				UserFollowingTagStatus        int    `json:"userFollowingTagStatus"`
				UserLanguage                  string `json:"userLanguage"`
				UserJoinUsedPointRank         int    `json:"userJoinUsedPointRank"`
				UserCurrentCheckinStreakEnd   int    `json:"userCurrentCheckinStreakEnd"`
				UserFollowingArticleStatus    int    `json:"userFollowingArticleStatus"`
				UserKeyboardShortcutsStatus   int    `json:"userKeyboardShortcutsStatus"`
				UserReplyWatchArticleStatus   int    `json:"userReplyWatchArticleStatus"`
				UserCommentViewMode           int    `json:"userCommentViewMode"`
				UserBreezemoonStatus          int    `json:"userBreezemoonStatus"`
				UserCheckinTime               int64  `json:"userCheckinTime"`
				UserUsedPoint                 int    `json:"userUsedPoint"`
				UserArticleStatus             int    `json:"userArticleStatus"`
				UserPoint                     int    `json:"userPoint"`
				UserCommentCount              int    `json:"userCommentCount"`
				UserIntro                     string `json:"userIntro"`
				UserMobileSkin                string `json:"userMobileSkin"`
				UserListPageSize              int    `json:"userListPageSize"`
				OId                           string `json:"oId"`
				UserName                      string `json:"userName"`
				UserGeoStatus                 int    `json:"userGeoStatus"`
				UserLongestCheckinStreakStart int    `json:"userLongestCheckinStreakStart"`
				UserSkin                      string `json:"userSkin"`
				UserNotifyStatus              int    `json:"userNotifyStatus"`
				UserFollowingUserStatus       int    `json:"userFollowingUserStatus"`
				UserArticleCount              int    `json:"userArticleCount"`
				UserRole                      string `json:"userRole"`
			} `json:"articleAuthor"`
		} `json:"article"`
		Pagination struct {
			PaginationPageCount int   `json:"paginationPageCount"`
			PaginationPageNums  []int `json:"paginationPageNums"`
		} `json:"pagination"`
	} `json:"data"`
}
