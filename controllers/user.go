package controllers

import (
	"hsproom/config"
	"hsproom/models"
	"hsproom/templates"
	"hsproom/utils"
	"net/http"
	"os"
	"strconv"
)

type userViewMember struct {
	*templates.DefaultMember
	UserInfo     *models.User
	UserPrograms *[]models.ProgramInfo
}

func userViewHandler(document http.ResponseWriter, request *http.Request) {

	var err error

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userView.tmpl"

	rawUid := request.URL.Query().Get("u")
	uid, err := strconv.Atoi(rawUid)

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "ユーザが見つかりませんでした。")
		return
	}

	var user models.User
	err = user.Load(uid)

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "ユーザが見つかりませんでした。")
		return
	}

	var programs []models.ProgramInfo

	_, err = models.GetProgramListByUser(models.ProgramColCreated, &programs, user.Name, true, 0, 4)

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	err = tmpl.Render(document, userViewMember{
		DefaultMember: &templates.DefaultMember{
			Title: user.Name + " のプロフィール - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		UserInfo:     &user,
		UserPrograms: &programs,
	})
	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "ページの表示に失敗しました。管理人へ問い合わせてください。")
	}
}

func userListHandler(document http.ResponseWriter, request *http.Request) {

}

func userLogoutHandler(document http.ResponseWriter, request *http.Request) {

	session, err := sessionStore.Get(request, "go-wiki")

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "ログアウトに失敗しました。")

		return
	}

	session.Values["User"] = nil
	session.Save(request, document)

	http.Redirect(document, request, "http://localhost:8080/", 301)
}

func userEditHandler(document http.ResponseWriter, request *http.Request) {

}

type userProgramsMember struct {
	*templates.DefaultMember
	Programs     []models.ProgramInfo
	ProgramCount int
	CurPage      int
	MaxPage      int
	Sort         string
	UserName     string
	UserId       int
}

func userProgramsHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userPrograms.tmpl"

	userId, err := strconv.Atoi(request.URL.Query().Get("u"))

	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	sort := request.URL.Query().Get("s")

	var sortKey models.ProgramColumn
	switch sort {
	case "c":
		sortKey = models.ProgramColCreated
	case "g":
		sortKey = models.ProgramColGood
	case "n":
		sortKey = models.ProgramColTitle
	default:
		sortKey = models.ProgramColCreated
	}

	page, err := strconv.Atoi(request.URL.Query().Get("p"))
	if err != nil {
		page = 0
	}

	userName, err := models.GetUserName(userId)
	if err != nil {
		utils.PromulgateDebug(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	var programs []models.ProgramInfo
	i, err := models.GetProgramListByUser(sortKey, &programs, userName, true, page*10, 10)

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	maxPage := i / 10
	if i%10 == 0 {
		maxPage--
	}

	err = tmpl.Render(document, userProgramsMember{
		DefaultMember: &templates.DefaultMember{
			Title: userName + " - " + config.SiteTitle,
			User:  getSessionUser(request),
		},
		Programs:     programs,
		ProgramCount: i,
		CurPage:      page,
		MaxPage:      maxPage,
		Sort:         sort,
		UserName:     userName,
		UserId:       userId,
	})
}

type userSettingsMember struct {
	*templates.DefaultMember
	UserInfo models.User
}

func userSettingsHandler(document http.ResponseWriter, request *http.Request) {

	var tmpl templates.Template
	tmpl.Layout = "default.tmpl"
	tmpl.Template = "userSettings.tmpl"

	userId := getSessionUser(request)

	if userId == 0 {
		utils.PromulgateDebugStr(os.Stdout, "匿名の管理画面へのアクセス")

		showError(document, request, "ログインが必要です。")
		// TODO: ログインさせる。
		return
	}

	var user models.User
	err := user.Load(userId)

	if err != nil {
		utils.PromulgateFatal(os.Stdout, err)

		showError(document, request, "エラーが発生しました。")
		return
	}

	err = tmpl.Render(document, userSettingsMember{
		DefaultMember: &templates.DefaultMember{
			Title: "管理画面 - " + config.SiteTitle,
			User:  userId,
		},
		UserInfo: user,
	})
}
