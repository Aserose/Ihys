package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	tg2 "IhysBestowal/internal/service/webapi/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type buttonBack func(newline bool, newArgs argsBack, updateBasicArgs bool) tg2.Button
type argsBack func() (text, callbackData string, tap func(chatId int64, interMsgId int))

type TGMenu interface {
	Main(tr dto.Executor, msgId int, sourceData func(tr dto.Executor, msgId int) datastruct.AudioItems)
	Search(tr dto.Executor, msgId int, audio string)
}

type viewBuilderAuxiliary struct {
	resp         tgbotapi.MessageConfig
	tr           dto.Executor
	pageCapacity int
	msgId        int
	builder      tg2.TGMenu
	buttonBack   buttonBack
}

type menuService struct {
	api webapi.WebApiService
	iMiddleware
	iViewer
	builder tg2.TGMenu
	cfg     config.Buttons
}

func NewMenuService(api webapi.WebApiService, cfg config.Buttons) TGMenu {
	return menuService{
		api:     api,
		iViewer: newViewer(cfg, api),
		builder: api.ITelegram.NewMenuBuilder(),
		cfg:     cfg,
	}
}

func (ms menuService) Main(tr dto.Executor, msgId int, sourceData func(tr dto.Executor, msgId int) datastruct.AudioItems) {
	if ms.isNewMenu(msgId) {
		ms.deleteOldMenu(tr, msgId)
	}

	buttonBack := func(newline bool, newArgs argsBack, updateBasicArgs bool) tg2.Button {
		defaultArgs := func() (text, callbackData string, tap func(chatId int64, interMsgId int)) {
			return ms.cfg.MainMenu.Text, ms.cfg.MainMenu.CallbackData, func(chatId int64, interMsgId int) {
				ms.Main(tr, interMsgId, sourceData)
			}
		}
		build := func(args argsBack) tg2.Button {
			if newline {
				return ms.builder.NewLineMenuButton(args())
			}
			return ms.builder.NewMenuButton(args())
		}

		if newArgs != nil {
			return build(newArgs)
		}
		return build(defaultArgs)
	}

	ms.newMainMenu(
		viewBuilderAuxiliary{
			resp: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: tr.ChatID},
				Text:     "menu",
			},
			tr:           tr,
			pageCapacity: 7,
			msgId:        msgId,
			builder:      ms.builder,
			buttonBack:   buttonBack,
		},
	)

}

func (ms menuService) Search(tr dto.Executor, msgId int, requestSong string) {
	if ms.isNewMenu(msgId) {
		ms.deleteOldMenu(tr, msgId)
	}

	song := ms.findRequestAudio(requestSong)

	buttonBack := func(newline bool, newArgs argsBack, updateBasicArgs bool) tg2.Button {
		defaultArgs := func() (string, string, func(chatId int64, interMsgId int)) {
			return "back", ms.cfg.SearchMenu.CallbackData, func(chatId int64, interMsgId int) {
				ms.Search(tr, interMsgId, requestSong)
			}
		}
		build := func(args argsBack) tg2.Button {
			if newline {
				return ms.builder.NewLineMenuButton(args())
			}
			return ms.builder.NewMenuButton(args())
		}

		if newArgs != nil {
			return build(newArgs)
		}
		return build(defaultArgs)
	}

	ms.newSearchMenu(
		viewBuilderAuxiliary{
			resp: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{ChatID: tr.ChatID},
				Text:     song.Artist + " - " + song.Title,
			},
			tr:           tr,
			pageCapacity: 7,
			msgId:        msgId,
			builder:      ms.builder,
			buttonBack:   buttonBack,
		}, trimStr(song.Artist, 5)+"-"+trimStr(song.Title, 5),
		datastruct.AudioItems{
			From:  "",
			Items: []datastruct.AudioItem{song},
		})

}

func (ms menuService) findRequestAudio(requestAudio string) datastruct.AudioItem {
	return ms.api.GetAudio(requestAudio)
}

func (ms menuService) isNewMenu(msgId int) bool {
	return msgId == 0
}

func (ms menuService) deleteOldMenu(tr dto.Executor, msgId int) {
	if cmd, ok := tr.ExecCmd[ms.cfg.MainMenu.Delete]; ok {
		cmd(tr.ChatID, msgId)
	}
	if cmd, ok := tr.ExecCmd[ms.cfg.SearchMenu.Delete]; ok {
		cmd(tr.ChatID, msgId)
	}
}

func (ms menuService) newSearchMenu(aux viewBuilderAuxiliary, callbackData string, sourceData datastruct.AudioItems) {
	md := newMiddleware(ms.api, func(tr dto.Executor, msgId int) datastruct.AudioItems { return sourceData })
	aux.builder.MenuBuild(aux.resp, aux.msgId, aux.tr.ExecCmd,

		aux.builder.NewSubMenu(
			ms.cfg.SearchMenu.Text,
			ms.cfg.SearchMenu.CallbackData,

			aux.builder.NewLineMenuButton(
				ms.cfg.LastFm.Text,
				ms.cfg.LastFm.CallbackData,
				func(chatId int64, interMsgId int) {
					ms.showAudioList(aux.tr, callbackData,
						paginateAudioItems(md.GetLastFMSimiliars(aux.tr, interMsgId), aux.pageCapacity),
						interMsgId, aux.buttonBack)
				}),

			aux.builder.NewMenuButton(
				ms.cfg.YaMusic.Text,
				ms.cfg.YaMusic.CallbackData,
				func(chatId int64, interMsgId int) {
					ms.showAudioList(aux.tr, callbackData,
						paginateAudioItems(md.GetYaMusicSimiliars(aux.tr, interMsgId), aux.pageCapacity),
						interMsgId, aux.buttonBack)
				}),

			aux.builder.NewMenuButton(
				"all",
				"all",
				func(chatId int64, interMsgId int) {
					ms.showAudioList(aux.tr, callbackData,
						paginateAudioItems(ms.api.GetSimiliar(sourceData, true), aux.pageCapacity),
						interMsgId, aux.buttonBack)
				}),
		),
	)
}

func (ms menuService) newMainMenu(aux viewBuilderAuxiliary) {
	md := newMiddleware(ms.api, nil)
	callbackData := strconv.Itoa(int(aux.resp.ChatID)) + ms.cfg.MainMenu.CallbackData

	aux.builder.MenuBuild(aux.resp, aux.msgId, aux.tr.ExecCmd,

		aux.builder.NewSubMenu(ms.cfg.MainMenu.Text, ms.cfg.MainMenu.CallbackData,
			aux.builder.NewLineMenuButton(
				ms.cfg.LastFm.Text,
				ms.cfg.LastFm.CallbackData,
				func(chatId int64, interMsgId int) {
					ms.showAudioList(aux.tr, callbackData,
						paginateAudioItems(md.GetLastFMSimiliars(aux.tr, interMsgId), aux.pageCapacity),
						interMsgId, aux.buttonBack)
				}),

			aux.builder.NewSubMenu(
				ms.cfg.VkSubMenu.Self.Text,
				ms.cfg.VkSubMenu.Self.CallbackData,
				ms.openVkSubmenu(aux, callbackData)...),

			aux.builder.NewMenuButton(
				ms.cfg.YaMusic.Text,
				ms.cfg.YaMusic.CallbackData,
				func(chatId int64, interMsgId int) {
					ms.showAudioList(aux.tr, callbackData,
						paginateAudioItems(md.GetYaMusicSimiliars(aux.tr, interMsgId), aux.pageCapacity),
						interMsgId, aux.buttonBack)
				}),
		),
	)
}

func (ms menuService) openVkSubmenu(aux viewBuilderAuxiliary, callbackData string) []tg2.Button {
	if ms.isVkAuthorized(aux) {
		return ms.newDefaulVkSubmenu(aux, callbackData)
	}

	return ms.newAuthVkSubmenu(aux)
}

func (ms menuService) isVkAuthorized(aux viewBuilderAuxiliary) bool {
	return ms.api.IVk.Auth().IsAuthorized(aux.tr.TGUser)
}

func (ms menuService) newDefaulVkSubmenu(aux viewBuilderAuxiliary, callbackData string) []tg2.Button {
	md := newMiddleware(ms.api, nil)
	return []tg2.Button{

		aux.builder.NewMenuButton(
			ms.cfg.VkSubMenu.UserPlaylist.Text,
			ms.cfg.VkSubMenu.UserPlaylist.CallbackData,
			func(chatId int64, interMsgId int) {
				ms.showPlaylists(aux.tr, callbackData,
					paginatePlaylistItems(md.GetVKPlaylists(aux.tr, interMsgId), aux.pageCapacity),
					interMsgId, aux.buttonBack)
			}),

		aux.builder.NewMenuButton(
			ms.cfg.VkSubMenu.Recommendation.Text,
			ms.cfg.VkSubMenu.Recommendation.CallbackData,
			func(chatId int64, interMsgId int) {
				ms.showAudioList(aux.tr, callbackData,
					paginateAudioItems(md.GetVKSimiliars(aux.tr, interMsgId), aux.pageCapacity),
					interMsgId, aux.buttonBack)
			}),
	}
}

func (ms menuService) newAuthVkSubmenu(aux viewBuilderAuxiliary) []tg2.Button {
	return []tg2.Button{
		aux.builder.NewMenuButton(
			ms.cfg.VkSubMenu.Auth.Text,
			ms.cfg.VkSubMenu.Auth.CallbackData,
			func(chatId int64, interMsgId int) {
				aux.resp.Text = "follow the [link](" + ms.api.IVk.Auth().GetAuthURL() + ")" + ", allow access and send an accessToken by message."
				aux.resp.ParseMode = "Markdown"
				aux.builder.MenuBuild(aux.resp, interMsgId, aux.tr.ExecCmd, aux.buttonBack(true, nil, false))
			}),
		aux.buttonBack(true, nil, false),
	}
}
