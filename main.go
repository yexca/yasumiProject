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
	// longLayout
	screenComboBox *walk.ComboBox
	restComboBox   *walk.ComboBox
	// 说明界面
	explainTitle *walk.Label
	explainText  *walk.TextLabel
	// 倒计时界面
	countDownTitle      *walk.Label
	countDownText       *walk.Label
	countDownButton     *walk.PushButton
	countDownBackButton *walk.PushButton
	// 当前模式与持续次数
	mode  int // 1 为短频，2 为高频
	count int
	// 计时还是休息
	state int
	// 长时的选择
	screenTime int
	restTime   int
}

func newMyMainWindow() *MyMainWindow {
	return &MyMainWindow{
		count: 0,
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

	err = MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "yasumiProject",               // 窗口标题
		MinSize:  Size{Width: 300, Height: 500}, // 最小尺寸
		Size:     Size{Width: 300, Height: 500}, // 尺寸
		Layout:   VBox{},                        // 窗口布局
		//MenuItems: []MenuItem{},
		Children: []Widget{
			Composite{
				AssignTo: &mw.mainLayout,
				Visible:  true,
				Layout:   VBox{},
				Children: []Widget{
					VSpacer{},
					Label{
						Text:      "请选择模式",
						Alignment: AlignHCenterVCenter,
						Font:      Font{PointSize: 14},
					},
					VSpacer{},
					PushButton{
						Text:      "短时间高频休息",
						OnClicked: mw.shortLayoutToggle,
					},
					PushButton{
						Text:      "长时间低频休息",
						OnClicked: mw.LongSelectLayout,
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
						Text:          "长时间低频休息",
						TextAlignment: AlignCenter,
						Font:          Font{PointSize: 20, Bold: true},
					},
					VSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text:      "学习",
								Alignment: AlignHCenterVCenter,
							},
							ComboBox{
								AssignTo:     &mw.screenComboBox,
								Model:        screenSelect,
								CurrentIndex: 0,
							},
							Label{
								Text:      "分钟",
								Alignment: AlignHCenterVCenter,
							},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{
								Text:      "休息",
								Alignment: AlignHCenterVCenter,
							},
							ComboBox{
								AssignTo:     &mw.restComboBox,
								Model:        restSelect,
								CurrentIndex: 0,
							},
							Label{
								Text:      "分钟",
								Alignment: AlignHCenterVCenter,
							},
						},
					},
					VSpacer{},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							PushButton{
								Text:      "返回",
								OnClicked: mw.mainLayoutToggle,
							},
							PushButton{
								Text: "开始",
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
			// 说明的界面
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
								Text:      "返回",
								OnClicked: mw.mainLayoutToggle,
							},
							PushButton{
								Text: "开始",
								OnClicked: func() {
									if mw.mode == 1 {
										mw.shortScreenCountDownLayout()
									} else {
										err := beeep.Notify("a", "b", "./icon/app.ico")
										if err != nil {
											return
										}
									}
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
						Text:     "计时中···",
						Enabled:  false,
						OnClicked: func() {
						},
					},
					PushButton{
						AssignTo:  &mw.countDownBackButton,
						Text:      "主界面",
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

	mw.mainLayoutToggle()

	icon, err := walk.NewIconFromFile("./icon/app.ico")
	if err != nil {
		return
	}
	mw.SetIcon(icon)

	app.Run()

	//{
	//	Label{Text: "请选择模式", Alignment: AlignHCenterVCenter}.Create(NewBuilder(dynamicComposite))
	//	PushButton{Text: "短时间高频休息"}.Create(NewBuilder(dynamicComposite))
	//
	//	//dynamicComposite.Layout().Spacing()
	//}

	//MainWindow{
	//	AssignTo: &mainWindow,
	//	Title:    "Hi walk",                     // 窗口标题
	//	MinSize:  Size{Width: 400, Height: 500}, // 最小尺寸
	//	Size:     Size{Width: 400, Height: 500}, // 尺寸
	//	Layout:   VBox{},                        // 窗口布局
	//	Children: []Widget{
	//		PushButton{
	//			Text: "打开新窗口",
	//			OnClicked: func() {
	//				// 调用函数打开新窗口
	//				openNewWindow(mainWindow)
	//			},
	//		},
	//	},
	//}.Create()

	//app.Run()
}

//	func openNewWindow(owner walk.Form) {
//		var newWindow *walk.MainWindow
//
//		MainWindow{
//			AssignTo: &newWindow,
//			Title:    "新窗口",
//			MinSize:  Size{Width: 300, Height: 200},
//			Layout:   VBox{},
//			Children: []Widget{
//				Label{
//					Text: "这是一个新窗口",
//				},
//				PushButton{
//					Text: "关闭窗口",
//					OnClicked: func() {
//						newWindow.Close()
//					},
//				},
//			},
//		}.Create()
//	}

// 清除所有控件显示
func (mw *MyMainWindow) clearLayout() {
	mw.mainLayout.SetVisible(false)
	//mw.shortLayout.SetVisible(false)
	mw.longLayout.SetVisible(false)
	mw.explainLayout.SetVisible(false)
	mw.countDownLayout.SetVisible(false)
}

// 主界面
func (mw *MyMainWindow) mainLayoutToggle() {
	mw.clearLayout()
	mw.mainLayout.SetVisible(true)
}

// 短时高频说明界面
func (mw *MyMainWindow) shortLayoutToggle() {
	mw.clearLayout()
	//mw.shortLayout.SetVisible(true)

	mw.explainTitle.SetText("短时间高频休息")
	mw.explainText.SetText("此模式基于 20-20-20 规则：每使用屏幕 20 分钟，看向距离 20 英尺（约 6 米） 的物体 20 秒，可以有效缓解眼睛疲劳。\n\n" +
		"但是考虑到点击造成的时间延迟，本程序采用每 20 分钟休息 1 分钟的模式，同时每三个 20 分钟将进行 5 分钟的休息")

	// 切换模式
	mw.mode = 1

	mw.explainLayout.SetVisible(true)
}

// 短时高频看屏幕倒计时
func (mw *MyMainWindow) shortScreenCountDownLayout() {
	mw.clearLayout()

	mw.countDownTitle.SetText("短时间高频休息")
	mw.countDownText.SetText("倒计时 20 分钟")

	// 计时界面按钮修改
	mw.countDownButton.SetText("开始学习")
	mw.countDownButton.Clicked().Detach(0)
	mw.countDownButton.Clicked().Attach(func() {
		go countDown(1200, mw)
	})
	mw.countDownButton.SetEnabled(true)

	mw.count += 1

	// 看屏幕状态
	mw.state = 1

	mw.countDownLayout.SetVisible(true)

}

// 短时高频休息倒计时
func (mw *MyMainWindow) shortRestCountDownLayout() {
	mw.clearLayout()

	err := mw.countDownTitle.SetText("短时间高频休息")
	if err != nil {
		return
	}
	err = mw.countDownText.SetText("时间到，休息啦")
	if err != nil {
		return
	}

	// 计时界面按钮修改
	err = mw.countDownButton.SetText("开始休息")
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
	mw.longLayout.SetVisible(true)
	//beeep.Notify("Error", "还没做捏", "./icon/app.ico")
}

// 长时看屏幕倒计时
func (mw *MyMainWindow) longScreenCountDownLayout() {
	mw.clearLayout()

	mw.countDownTitle.SetText("长时间低频休息")
	mw.countDownText.SetText(fmt.Sprintf("倒计时 %d 分钟", mw.screenTime))

	// 计时界面按钮修改
	mw.countDownButton.SetText("开始学习")
	mw.countDownButton.Clicked().Detach(0)
	mw.countDownButton.Clicked().Attach(func() {
		go countDown(mw.screenTime*60, mw)
	})
	mw.countDownButton.SetEnabled(true)

	// 看屏幕状态
	mw.state = 1

	mw.countDownLayout.SetVisible(true)

}

// 长时休息倒计时
func (mw *MyMainWindow) longRestCountDownLayout() {
	mw.clearLayout()

	err := mw.countDownTitle.SetText("长时间低频休息")
	if err != nil {
		return
	}
	err = mw.countDownText.SetText("时间到，休息啦")
	if err != nil {
		return
	}

	// 计时界面按钮修改
	err = mw.countDownButton.SetText("开始休息")
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
		beeep.Notify("学习时间", "开始计时，一起努力加油吧", "./icon/app.ico")
	} else if mw.state == 2 {
		beeep.Notify("休息时间", "开始休息，记得向 6 米外的地方至少看 20 秒哦", "./icon/app.ico")
	}

	mw.countDownButton.SetText("计时中···")
	mw.countDownButton.SetEnabled(false)

	for i := seconds; i >= 0; i-- {
		m := i / 60
		s := i % 60

		time.Sleep(1 * time.Second)

		mw.countDownText.SetText(fmt.Sprintf("剩余时间：%2d:%2d", m, s))
	}

	// 根据状态调用界面
	if mw.state == 1 {
		beeep.Notify("休息时间", "休息时间到了，快点开始休息吧", "./icon/app.ico")
		mw.shortRestCountDownLayout()
	} else if mw.state == 2 {
		beeep.Notify("学习时间", "休息结束，继续一起努力学习吧", "./icon/app.ico")
		mw.longScreenCountDownLayout()
	}

	// 返回按钮开启
	mw.countDownBackButton.SetEnabled(true)
}
