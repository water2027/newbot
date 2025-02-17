package common

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"

	"testsse/config"
	"testsse/model"
	"testsse/sseapi"
)

// 目标群
var targetGroup []*openwechat.Group
var mu sync.Mutex

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
)

// 发送消息的通道
var PostChan = make(chan model.Post)

// 防止微信自动退出登录
func keepAlive(ctx context.Context, bot *openwechat.Self) {
	defer wg.Done()
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("keepAlive stop")
			return
		default:
		}
		randomMinutes := time.Duration(10+rand.Intn(11)) * time.Minute // 10-20分钟
		timer := time.NewTimer(randomMinutes)

		select {
		case <-ctx.Done():
			timer.Stop()
			log.Println("keepAlive stop")
			return
		case <-timer.C:
			heartBeat(bot)
		}
	}
}

func heartBeat(bot *openwechat.Self) {
	// 生成要发送的消息
	outMessage := fmt.Sprintf("防微信自动退出登录[%d]", time.Now().Unix())
	bot.SendTextToFriend(openwechat.NewFriendHelper(bot), outMessage)
}

// 打印二维码到控制台
func consoleQrcode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToSmallString(true))
}

// 需要开一个协程运行
func StartBot() {
	ctx, cancel = context.WithCancel(context.Background())

	wg.Add(3)

	config := config.GetBotConfig()

	go runBot(ctx)

	//保持登录

	//发送消息
	go sentPostToGroup(ctx, config.Str)
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Duration(config.TimeInterval) * time.Minute)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				sseapi.GetPosts(PostChan, &config)
			}
		}
	}(ctx)
	wg.Wait()
}

func runBot(ctx context.Context) {
	defer wg.Done()
	var err error
	bot := openwechat.DefaultBot(openwechat.Desktop)

	var errorChan = make(chan error, 1)

	//打印二维码到控制台
	bot.UUIDCallback = consoleQrcode
	//去掉心跳检测
	bot.SyncCheckCallback = nil
	//热登录，第一次会生成一个storage.json文件，用于存储登录信息。另一种方式是不扫码登录和每次都扫描登录，但是不扫码登录但是发现不行
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	if err = bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		log.Println("loginErr", err)
		return
	}

	//获取登录的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Println("getCurrentUserErr", err)
		return
	}

	//获取群组，必须是手机端通讯录里保存的群
	groups, err := self.Groups()
	if err != nil {
		log.Println("getGroupsErr", err)
		return
	}

	config := config.GetBotConfig()

	for _, groupName := range config.TargetGroupName {
		group := groups.GetByNickName(groupName)
		if group == nil {
			log.Println("groupNotFound")
			return
		}
		AddGroup(group)
	}

	//消息处理函数,接受到信息时触发
	//之后可以考虑接入大模型/连接数据库
	bot.MessageHandler = nil

	go func() {
		errorChan <- bot.Block()
	}()

	go keepAlive(ctx, self)

	select {
	case <-ctx.Done():
		log.Println("bot stop")
		bot.Logout()
	case err := <-errorChan:
		if err != nil {
			log.Println("Bot logout error:", err)
		}
		cancel()
	}
}

func sentPostToGroup(ctx context.Context, str string) {
	defer wg.Done()
	for {
		botConfig := config.GetBotConfig()
		urlstr := botConfig.Url
		select {
		case post := <-PostChan:
			if post.PostID > botConfig.StartNum {
				for _, group := range getGroup() {
					botConfig.StartNum = post.PostID
					url := fmt.Sprintf(urlstr, post.PostID)
					msg := fmt.Sprintf(str, post.Title, url)

					_, err := group.SendText(msg)
					time.Sleep(5 * time.Second)
					if err != nil {
						log.Println(err)
					}
					config.UpdateBotConfig(botConfig)
				}
			}
		case <-ctx.Done():
			fmt.Println("sendPost stop")
			return
		}
	}
}

func getGroup() []*openwechat.Group {
	mu.Lock()
	defer mu.Unlock()
	return targetGroup
}

func AddGroup(group *openwechat.Group) {
	mu.Lock()
	defer mu.Unlock()
	targetGroup = append(targetGroup, group)
}

func RemoveGroup(name string) {
	mu.Lock()
	defer mu.Unlock()
	for i, g := range targetGroup {
		if g.NickName == name {
			targetGroup = append(targetGroup[:i], targetGroup[i+1:]...)
			break
		}
	}
}

func StopBot() {
	fmt.Println("start stop")
	if cancel != nil {
		cancel()
	}
	wg.Wait()
	log.Println("stop over")
}

func RestartBot() {
	log.Println("stop!")
	StopBot()
	log.Println("start!")
	go StartBot()
}
