package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/tailscale/walk"
	. "github.com/tailscale/walk/declarative"
	"log"
	"strconv"
	"time"
)

type MyMainWindow struct {
	*walk.MainWindow
	mainLayout *walk.Composite
	// shortLayout     *walk.Composite
	longLayout      *walk.Composite
	explainLayout   *walk.Composite
	countDownLayout *walk.Composite

	// mainLayout
	mainTitleLabel  *walk.Label
	mainButtonLeft  *walk.PushButton
	mainButtonright *walk.PushButton

	// longLayout
	screenComboBox        *walk.ComboBox
	restComboBox          *walk.ComboBox
	longTitle             *walk.Label
	LongSelectEvery       *walk.Label
	LongSelectMinute1     *walk.Label
	LongSelectMinute2     *walk.Label
	LongSelectRest        *walk.Label
	LongSelectButtonLeft  *walk.PushButton
	LongSelectButtonRight *walk.PushButton

	// language
	languageSelectLayout *walk.Composite
	languageComboBox     *walk.ComboBox
	languageButton       *walk.PushButton

	// 说明界面
	explainTitle       *walk.Label
	explainText        *walk.TextLabel
	explainButtonLeft  *walk.PushButton
	explainButtonRight *walk.PushButton

	// 倒计时界面
	countDownTitle      *walk.Label
	countDownText       *walk.Label
	countDownButton     *walk.PushButton
	countDownBackButton *walk.PushButton

	// 当前模式与持续次数
	// mode  int // 1 为短频，2 为高频
	count int
	// 计时还是休息
	state int
	// 长时的选择
	screenTime int
	restTime   int
	// 语言
	language string
	// 通知图标
	notifyIcon string
}

func newMyMainWindow() *MyMainWindow {
	return &MyMainWindow{
		count:      0,
		notifyIcon: "./icon/app.ico",
		language:   "English",
	}
}

func main() {
	app, err := walk.InitApp()
	if err != nil {
		log.Fatal(err)
	}

	mw := newMyMainWindow()

	// 长时的时间
	screenSelect := []string{"30", "35", "40", "45"}
	restSelect := []string{"5", "10"}
	// language
	languageSelect := []string{"English", "简体中文", "日本語"}

	err = MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "yasumiProject",               // 窗口标题
		MinSize:  Size{Width: 400, Height: 600}, // 最小尺寸
		Size:     Size{Width: 400, Height: 600}, // 尺寸
		MaxSize:  Size{Width: 450, Height: 650}, // 尺寸
		Layout:   VBox{},                        // 窗口布局
		//MenuItems: []MenuItem{},
		Children: []Widget{
			Composite{
				AssignTo: &mw.mainLayout,
				Visible:  false,
				Layout:   VBox{},
				Children: []Widget{
					VSpacer{},
					Label{
						AssignTo:  &mw.mainTitleLabel,
						Text:      "",
						Alignment: AlignHCenterVCenter,
						Font:      Font{PointSize: 14},
					},
					VSpacer{},
					PushButton{
						AssignTo:  &mw.mainButtonLeft,
						Text:      "",
						OnClicked: mw.shortLayoutToggle,
					},
					PushButton{
						AssignTo:  &mw.mainButtonright,
						Text:      "",
						OnClicked: mw.LongSelectLayout,
					},
					VSpacer{},
				},
			},
			Composite{
				AssignTo: &mw.languageSelectLayout,
				Visible:  true,
				Layout:   VBox{},
				Children: []Widget{
					VSpacer{},
					Label{
						Text: "Please select language\n" + "请选择语言\n" + "言語を選んでください",
						//MaxSize: Size{Height: 150, Width: 400},
					},
					//Label{
					//	Text: "Please select language",
					//},
					//Label{
					//	Text: "请选择语言",
					//},
					//Label{
					//	Text: "言語を選んでください",
					//},

					ComboBox{
						AssignTo:     &mw.languageComboBox,
						Model:        languageSelect,
						CurrentIndex: 1,
						//MaxSize:      Size{Height: 100, Width: 200},
						OnCurrentIndexChanged: func() {
							mw.language = languageSelect[mw.languageComboBox.CurrentIndex()]
							switch mw.language {
							case "English":
								mw.languageButton.SetText("next")
							case "简体中文":
								mw.languageButton.SetText("下一步")
							case "日本語":
								mw.languageButton.SetText("次へ")
							}
						},
					},
					VSpacer{},
					PushButton{
						AssignTo: &mw.languageButton,
						Text:     "下一步",
						OnClicked: func() {
							mw.mainLayoutToggle()
						},
					},
					VSpacer{},
				},
			},
			//Composite{
			//	AssignTo: &mw.shortLayout,
			//	Visible:  false,
			//	Layout:   VBox{},
			//	Children: []Widget{
			//		Label{
			//			Text: "short",
			//		},
			//		PushButton{
			//			Text: "button",
			//		},
			//	},
			//},
			Composite{
				AssignTo: &mw.longLayout,
				Visible:  false,
				Layout:   VBox{},
				Children: []Widget{
					VSpacer{},
					Label{
						AssignTo:      &mw.longTitle,
						Text:          "",
						TextAlignment: AlignCenter,
						Font:          Font{PointSize: 20, Bold: true},
					},
					VSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								AssignTo:  &mw.LongSelectEvery,
								Text:      "",
								Alignment: AlignHCenterVCenter,
							},
							ComboBox{
								AssignTo:     &mw.screenComboBox,
								Model:        screenSelect,
								CurrentIndex: 0,
							},
							Label{
								AssignTo:  &mw.LongSelectMinute1,
								Text:      "",
								Alignment: AlignHCenterVCenter,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								AssignTo:  &mw.LongSelectRest,
								Text:      "",
								Alignment: AlignHCenterVCenter,
							},
							ComboBox{
								AssignTo:     &mw.restComboBox,
								Model:        restSelect,
								CurrentIndex: 0,
							},
							Label{
								AssignTo:  &mw.LongSelectMinute2,
								Text:      "",
								Alignment: AlignHCenterVCenter,
							},
						},
					},
					VSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{
								AssignTo:  &mw.LongSelectButtonLeft,
								Text:      "",
								OnClicked: mw.mainLayoutToggle,
							},
							PushButton{
								AssignTo: &mw.LongSelectButtonRight,
								Text:     "",
								OnClicked: func() {
									screenTime, _ := strconv.Atoi(screenSelect[mw.screenComboBox.CurrentIndex()])
									restTime, _ := strconv.Atoi(restSelect[mw.restComboBox.CurrentIndex()])
									mw.screenTime = screenTime
									mw.restTime = restTime
									mw.longScreenCountDownLayout()
								},
							},
						},
					},
				},
			},
			// (短时)说明的界面
			Composite{
				AssignTo: &mw.explainLayout,
				Visible:  false,
				Layout:   VBox{},
				Children: []Widget{
					Label{
						AssignTo:      &mw.explainTitle,
						Text:          "",
						TextAlignment: AlignCenter,
						Font:          Font{PointSize: 20, Bold: true},
					},
					TextLabel{
						AssignTo: &mw.explainText,
						Text:     "",
						MinSize:  Size{Width: 350, Height: 200},
						Font:     Font{PointSize: 14},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{
								AssignTo:  &mw.explainButtonLeft,
								Text:      "",
								OnClicked: mw.mainLayoutToggle,
							},
							PushButton{
								AssignTo: &mw.explainButtonRight,
								Text:     "",
								OnClicked: func() {
									mw.shortScreenCountDownLayout()
									//if mw.mode == 1 {
									//	mw.shortScreenCountDownLayout()
									//} else {
									//	err := beeep.Notify("a", "b", "./icon/app.ico")
									//	if err != nil {
									//		return
									//	}
									//}
								},
							},
						},
					},
				},
			},
			// 倒计时的界面
			Composite{
				AssignTo: &mw.countDownLayout,
				Visible:  false,
				Layout:   VBox{},
				Children: []Widget{
					VSpacer{},
					Label{
						AssignTo:      &mw.countDownTitle,
						Text:          "",
						TextAlignment: AlignCenter,
						Font:          Font{PointSize: 20, Bold: true},
					},
					VSpacer{},
					Label{
						AssignTo:      &mw.countDownText,
						Text:          "",
						TextAlignment: AlignCenter,
						Font:          Font{PointSize: 14},
					},
					VSpacer{},
					PushButton{
						AssignTo: &mw.countDownButton,
						Text:     "",
						Enabled:  false,
						OnClicked: func() {
						},
					},
					PushButton{
						AssignTo:  &mw.countDownBackButton,
						Text:      "",
						Enabled:   true,
						OnClicked: mw.mainLayoutToggle,
					},
				},
			},
		},
	}.Create()
	if err != nil {
		return
	}

	mw.clearLayout()
	mw.SetSize(walk.Size{Width: 400, Height: 600})
	mw.languageSelectLayout.SetVisible(true)

	icon, err := walk.NewIconFromFile(mw.notifyIcon)
	if err != nil {
		return
	}
	mw.SetIcon(icon)

	app.Run()

}

// 清除所有控件显示
func (mw *MyMainWindow) clearLayout() {
	mw.mainLayout.SetVisible(false)
	mw.languageSelectLayout.SetVisible(false)
	//mw.shortLayout.SetVisible(false)
	mw.longLayout.SetVisible(false)
	mw.explainLayout.SetVisible(false)
	mw.countDownLayout.SetVisible(false)
}

// 主界面
func (mw *MyMainWindow) mainLayoutToggle() {
	mw.clearLayout()
	mw.mainTitleLabel.SetText(GetText(mw.language, "modeSelect"))
	mw.mainButtonLeft.SetText(GetText(mw.language, "shortMode"))
	mw.mainButtonright.SetText(GetText(mw.language, "longMode"))
	mw.mainLayout.SetVisible(true)
}

// 短时高频说明界面
func (mw *MyMainWindow) shortLayoutToggle() {
	mw.clearLayout()
	//mw.shortLayout.SetVisible(true)

	mw.explainTitle.SetText(GetText(mw.language, "shortMode"))
	mw.explainText.SetText(GetText(mw.language, "shortExplainText"))
	mw.explainButtonLeft.SetText(GetText(mw.language, "back"))
	mw.explainButtonRight.SetText(GetText(mw.language, "next"))

	// 切换模式
	//mw.mode = 1

	mw.explainLayout.SetVisible(true)
}

// 短时高频看屏幕倒计时
func (mw *MyMainWindow) shortScreenCountDownLayout() {
	mw.clearLayout()

	mw.countDownTitle.SetText(GetText(mw.language, "shortMode"))
	mw.countDownText.SetText(GetText(mw.language, "countDown") + " 20 " + GetText(mw.language, "minute"))

	// 计时界面按钮修改
	mw.countDownButton.SetText(GetText(mw.language, "startStudyButton"))
	mw.countDownButton.Clicked().Detach(0)
	mw.countDownButton.Clicked().Attach(func() {
		go countDown(1200, mw)
	})
	mw.countDownButton.SetEnabled(true)

	mw.countDownBackButton.SetText(GetText(mw.language, "home"))

	mw.count += 1

	// 看屏幕状态
	mw.state = 1

	mw.countDownLayout.SetVisible(true)

}

// 短时高频休息倒计时
func (mw *MyMainWindow) shortRestCountDownLayout() {
	mw.clearLayout()

	err := mw.countDownTitle.SetText(GetText(mw.language, "shortMode"))
	if err != nil {
		return
	}
	err = mw.countDownText.SetText(GetText(mw.language, "restLabel"))
	if err != nil {
		return
	}

	// 计时界面按钮修改
	err = mw.countDownButton.SetText(GetText(mw.language, "startRestButton"))
	if err != nil {
		return
	}
	mw.countDownButton.Clicked().Detach(0)
	mw.countDownButton.Clicked().Attach(func() {
		if mw.count < 4 {
			go countDown(60, mw)
		} else {
			mw.count = 0
			go countDown(300, mw)
		}
	})
	mw.countDownButton.SetEnabled(true)

	// 休息状态
	mw.state = 2

	mw.countDownLayout.SetVisible(true)

}

// 长时选择界面
func (mw *MyMainWindow) LongSelectLayout() {
	mw.clearLayout()
	mw.longTitle.SetText(GetText(mw.language, "longMode"))
	mw.LongSelectEvery.SetText(GetText(mw.language, "every"))
	mw.LongSelectMinute1.SetText(GetText(mw.language, "minute"))
	mw.LongSelectMinute2.SetText(GetText(mw.language, "minute"))
	mw.LongSelectRest.SetText(GetText(mw.language, "rest"))
	mw.LongSelectButtonLeft.SetText(GetText(mw.language, "back"))
	mw.LongSelectButtonRight.SetText(GetText(mw.language, "next"))
	mw.longLayout.SetVisible(true)
	//beeep.Notify("Error", "还没做捏", "./icon/app.ico")
}

// 长时看屏幕倒计时
func (mw *MyMainWindow) longScreenCountDownLayout() {
	mw.clearLayout()

	mw.countDownTitle.SetText(GetText(mw.language, "longMode"))
	mw.countDownText.SetText(GetText(mw.language, "countDown") + fmt.Sprintf(" %d ", mw.screenTime) + GetText(mw.language, "minute"))

	// 计时界面按钮修改
	mw.countDownButton.SetText(GetText(mw.language, "startStudyButton"))
	mw.countDownButton.Clicked().Detach(0)
	mw.countDownButton.Clicked().Attach(func() {
		go countDown(mw.screenTime*60, mw)
	})
	mw.countDownButton.SetEnabled(true)

	mw.countDownBackButton.SetText(GetText(mw.language, "home"))

	// 看屏幕状态
	mw.state = 1

	mw.countDownLayout.SetVisible(true)

}

// 长时休息倒计时
func (mw *MyMainWindow) longRestCountDownLayout() {
	mw.clearLayout()

	err := mw.countDownTitle.SetText(GetText(mw.language, "longMode"))
	if err != nil {
		return
	}
	err = mw.countDownText.SetText(GetText(mw.language, "restLabel"))
	if err != nil {
		return
	}

	// 计时界面按钮修改
	err = mw.countDownButton.SetText(GetText(mw.language, "startRestButton"))
	if err != nil {
		return
	}
	mw.countDownButton.Clicked().Detach(0)
	mw.countDownButton.Clicked().Attach(func() {
		go countDown(mw.restTime*60, mw)
	})
	mw.countDownButton.SetEnabled(true)

	// 休息状态
	mw.state = 2

	mw.countDownLayout.SetVisible(true)

}

// 倒计时函数
func countDown(seconds int, mw *MyMainWindow) {
	// 返回按钮关闭
	mw.countDownBackButton.SetEnabled(false)

	// 根据状态调用通知
	if mw.state == 1 {
		beeep.Notify(GetText(mw.language, "studyStartNotifyTitle"), GetText(mw.language, "studyStartNotifyContent"), mw.notifyIcon)
	} else if mw.state == 2 {
		beeep.Notify(GetText(mw.language, "restStartNotifyTitle"), GetText(mw.language, "restStartNotifyContent"), mw.notifyIcon)
	}

	mw.countDownButton.SetText(GetText(mw.language, "Timing"))
	mw.countDownButton.SetEnabled(false)

	for i := seconds; i >= 0; i-- {
		m := i / 60
		s := i % 60

		time.Sleep(1 * time.Second)

		mw.countDownText.SetText(GetText(mw.language, "remain") + fmt.Sprintf("：%2d:%2d", m, s))
	}

	// 根据状态调用界面
	if mw.state == 1 {
		beeep.Notify(GetText(mw.language, "restStartNotifyTitle"), GetText(mw.language, "studyEndNotifyContent"), mw.notifyIcon)
		mw.shortRestCountDownLayout()
	} else if mw.state == 2 {
		beeep.Notify(GetText(mw.language, "studyStartNotifyTitle"), GetText(mw.language, "restEndNotifyContent"), mw.notifyIcon)
		mw.longScreenCountDownLayout()
	}

	// 返回按钮开启
	mw.countDownBackButton.SetEnabled(true)
}
