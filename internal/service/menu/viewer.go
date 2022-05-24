package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/service/webapi"
	tg2 "IhysBestowal/internal/service/webapi/tg"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type iViewer interface {
	showAudioList(tr dto.Executor, callbackData string, audioItems paginatedAudioItems, msgId int, back buttonBack)
	showPlaylists(tr dto.Executor, callbackData string, items paginatedPlaylistItems, msgId int, back buttonBack)
}

type viewer struct {
	api          webapi.WebApiService
	pageCapacity int
	cfg          config.Buttons
}

func newViewer(cfg config.Buttons, api webapi.WebApiService) iViewer {
	return viewer{
		api:          api,
		pageCapacity: 7,
		cfg:          cfg,
	}
}

func (ms viewer) showPlaylists(tr dto.Executor, callbackData string, paginatedItems paginatedPlaylistItems, msgId int, back buttonBack) {
	buttonBuilder := ms.api.ITelegram.NewMenuBuilder()
	resp := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{ChatID: tr.ChatID},
		Text:     "playlists",
	}
	viewControl := newViewController(paginatedItems.pageCount, back(true, nil, false))

	contentEnumeration := func(playlistItems []datastruct.PlaylistItem) []tg2.Button {
		playlistBtns := make([]tg2.Button, len(playlistItems))
		for i, playlist := range playlistItems {
			playlistBtns[i] = ms.newPlaylistButton(
				viewBuilderAuxiliary{
					resp:         resp,
					tr:           tr,
					pageCapacity: ms.pageCapacity,
					msgId:        msgId,
					builder:      buttonBuilder,
					buttonBack:   back,
				}, callbackData,
				playlist)
		}
		return playlistBtns
	}

	ms.setPageContent(viewControl, callbackData,
		func() []tg2.Button {
			return contentEnumeration(paginatedItems.items[viewControl.GetCurrentPageNumber()-1])
		},
		func(buttons []tg2.Button, interMsgId int) {
			buttonBuilder.MenuBuild(resp, interMsgId, tr.ExecCmd, buttons...)
		})(contentEnumeration(paginatedItems.items[viewControl.GetCurrentPageNumber()-1]), msgId)
}

func (ms viewer) showAudioList(tr dto.Executor, callbackData string, paginatedItems paginatedAudioItems, msgId int, back buttonBack) {
	buttonBuilder := ms.api.ITelegram.NewMenuBuilder()
	resp := tgbotapi.MessageConfig{BaseChat: tgbotapi.BaseChat{ChatID: tr.ChatID}, Text: paginatedItems.from}
	viewControl := newViewController(paginatedItems.pageCount, back(true, nil, false))

	contentEnumeration := func(items []datastruct.AudioItem) []tg2.Button {
		audioButtons := make([]tg2.Button, len(items))
		aux := viewBuilderAuxiliary{
			resp:       resp,
			tr:         tr,
			msgId:      msgId,
			builder:    buttonBuilder,
			buttonBack: back,
		}
		for i, music := range items {
			audioButtons[i] = ms.newAudioButton(aux, music)
		}
		return audioButtons
	}

	ms.setPageContent(
		viewControl, callbackData,
		func() []tg2.Button {
			return contentEnumeration(paginatedItems.items[viewControl.GetCurrentPageNumber()-1])
		},
		func(buttons []tg2.Button, interMsgId int) {
			buttonBuilder.MenuBuild(resp, interMsgId, tr.ExecCmd, buttons...)
		},
	)(contentEnumeration(paginatedItems.items[viewControl.GetCurrentPageNumber()-1]), msgId)
}

func (ms viewer) showSong(song datastruct.AudioItem, callbackData string, tr dto.Executor, msgId int) {
	builder := ms.api.ITelegram.NewMenuBuilder()
	songName := song.Artist + " - " + song.Title

	back := func(newline bool, newArgs argsBack, updateBasicArgs bool) tg2.Button {
		getArgs := func() (text, callbackData string, tap func(chatId int64, interMsgId int)) {
			return "back", callbackData + "back", func(chatId int64, interMsgId int) {
				ms.showSong(song, callbackData, tr, interMsgId)
			}
		}
		if newArgs != nil {
			if newline {
				return builder.NewLineMenuButton(newArgs())
			}
			return builder.NewMenuButton(newArgs())
		}
		if newline {
			return builder.NewLineMenuButton(getArgs())
		}
		return builder.NewMenuButton(getArgs())
	}

	resp := tgbotapi.NewMessage(tr.ChatID, " ")
	resp.ParseMode = "markdown"
	YTurl := ms.api.IYouTube.GetYTUrl(songName)
	songName = "\n" + "[" + songName + "]" + "(" + song.Url + ")"
	if YTurl != " " {
		resp.Text = songName + "\n\n[YouTube](" + YTurl + ")"
	} else {
		resp.Text = songName
	}

	ms.newSongMenu(
		viewBuilderAuxiliary{
			resp:       resp,
			tr:         tr,
			msgId:      msgId,
			builder:    builder,
			buttonBack: back,
		}, callbackData, song)
}

func (ms viewer) newSongMenu(aux viewBuilderAuxiliary, callbackData string, song datastruct.AudioItem) {
	aux.builder.MenuBuild(aux.resp, aux.msgId, aux.tr.ExecCmd,

		aux.builder.NewMenuButton(
			ms.cfg.SongMenu.Delete.Text,
			ms.cfg.SongMenu.Delete.CallbackData,
			func(chatID int64, interMsgId int) {
				dltcfg := tgbotapi.NewDeleteMessage(chatID, interMsgId)
				delete(aux.tr.ExecCmd, callbackData)
				ms.api.Send(dltcfg)
			}),

		aux.builder.NewMenuButton(
			ms.cfg.SongMenu.Similiars.Text,
			ms.cfg.SongMenu.Similiars.CallbackData+callbackData,
			func(chatID int64, interMsgId int) {
				ms.showAudioList(
					aux.tr, callbackData,
					paginateAudioItems(
						ms.api.GetSimiliar(datastruct.AudioItems{Items: []datastruct.AudioItem{song}}, true),
						ms.pageCapacity),
					interMsgId,
					aux.buttonBack)
			}),

		aux.builder.NewMenuButton(
			ms.cfg.SongMenu.Best.Text,
			ms.cfg.SongMenu.Best.CallbackData+callbackData,
			func(chatID int64, interMsgId int) {
				ms.showAudioList(
					aux.tr, callbackData,
					paginateAudioItems(ms.api.GetTopSongs(song.Artist), ms.pageCapacity),
					interMsgId,
					aux.buttonBack)
			}))
}

func (ms viewer) newPlaylistButton(aux viewBuilderAuxiliary, callbackData string, pl datastruct.PlaylistItem) tg2.Button {
	return aux.builder.NewLineMenuButton(
		pl.Title,
		fmt.Sprint("playlist", strconv.Itoa(pl.ID)),
		func(chatId int64, interMsgId int) {
			aux.resp.Text = pl.Title
			md := newMiddleware(ms.api,
				func(tr dto.Executor, msgId int) datastruct.AudioItems {
					res, err := ms.api.IVk.GetPlaylistSongs(tr.TGUser, pl.ID, pl.OwnerId)
					if err != nil {
						ms.api.Send(tgbotapi.NewEditMessageText(tr.ChatID, msgId, err.Error()))
					}
					return res
				})

			aux.builder.MenuBuild(aux.resp, interMsgId, aux.tr.ExecCmd,
				aux.builder.NewMenuButton(
					ms.cfg.UserPlaylist.Text,
					ms.cfg.UserPlaylist.CallbackData,
					func(chatId int64, interMsgId int) {
						ms.showPlaylists(aux.tr, callbackData,
							paginatePlaylistItems(md.GetVKPlaylists(aux.tr, interMsgId), ms.pageCapacity),
							interMsgId, aux.buttonBack)
					}),

				aux.builder.NewMenuButton(
					ms.cfg.LastFm.Text,
					ms.cfg.LastFm.CallbackData,
					func(chatId int64, interMsgId int) {
						ms.showAudioList(aux.tr, callbackData,
							paginateAudioItems(md.GetLastFMSimiliars(aux.tr, interMsgId), ms.pageCapacity),
							interMsgId, aux.buttonBack)
					}),

				aux.builder.NewMenuButton(
					ms.cfg.YaMusic.Text,
					ms.cfg.YaMusic.CallbackData,
					func(chatId int64, interMsgId int) {
						ms.showAudioList(aux.tr, callbackData,
							paginateAudioItems(md.GetYaMusicSimiliars(aux.tr, interMsgId), ms.pageCapacity),
							interMsgId, aux.buttonBack)
					}),
			)
		})
}
func (ms viewer) newAudioButton(aux viewBuilderAuxiliary, song datastruct.AudioItem) tg2.Button {
	songName := song.Artist + " - " + song.Title
	callbackData := trimStr(song.Artist, 5) + "-" + trimStr(song.Title, 5)

	return aux.builder.NewLineMenuButton(
		songName,
		callbackData,
		func(chatId int64, interMsgId int) {
			ms.showSong(song, callbackData, aux.tr, 0)
		})
}

func (ms viewer) setPageContent(viewControl iViewController, callbackData string, contentEnumeration func() []tg2.Button, build func(buttons []tg2.Button, interMsgId int)) (setPageContent func(buttons []tg2.Button, interMsgId int)) {
	return func(buttons []tg2.Button, interMsgId int) {
		ms.addPageControls(
			viewControl, callbackData, buttons,
			func() { setPageContent(contentEnumeration(), interMsgId) },
			build, interMsgId)
	}
}

func (ms viewer) addPageControls(viewControl iViewController, callbackData string, buttons []tg2.Button, setPageContentWrap func(), menuBuild func(buttons []tg2.Button, interMsgId int), interMsgId int) {
	builder := ms.api.ITelegram.NewMenuBuilder()

	switch viewControl.IsFirstPage() {
	case true:
		switch viewControl.IsLastPage() {
		case true:
			buttons = append(buttons, viewControl.GetControlButtons(callbackData, false, false, setPageContentWrap, builder)...)
		case false:
			buttons = append(buttons, viewControl.GetControlButtons(callbackData, false, true, setPageContentWrap, builder)...)
		}
	case false:
		switch viewControl.IsLastPage() {
		case true:
			buttons = append(buttons, viewControl.GetControlButtons(callbackData, true, false, setPageContentWrap, builder)...)
		case false:
			buttons = append(buttons, viewControl.GetControlButtons(callbackData, true, true, setPageContentWrap, builder)...)
		}
	}

	menuBuild(buttons, interMsgId)
}

func trimStr(str string, limit uint8) string {
	if uint8(len(str)) > limit {
		return string([]rune(str)[:limit])
	}
	return str
}
