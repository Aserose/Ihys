package menu

import (
	"IhysBestowal/internal/config"
	"IhysBestowal/internal/datastruct"
	"IhysBestowal/internal/dto"
	"IhysBestowal/internal/repository"
	"IhysBestowal/internal/service/webapi"
)

const (
	backText     = "back"
	backCallback = "back"
	empty        = ``
)

type TGMenu interface {
	Main(p dto.Response)
	Search(tr dto.Response, query string)
}

type menuService struct {
	md middleware
	viewer
	cfg config.Keypads
}

func NewMenuService(api webapi.WebApiService, storage repository.TrackStorage, cfg config.Keypads) TGMenu {
	md := newMiddleware(api, storage)

	return menuService{
		md:     md,
		viewer: newViewer(cfg, md, api),
		cfg:    cfg,
	}
}

func (ms menuService) Search(p dto.Response, query string) {
	ms.newSearchMenu(p, datastruct.AudioItems{
		From:  empty,
		Items: []datastruct.AudioItem{ms.md.search(query)},
	})
}

func (ms menuService) newSearchMenu(p dto.Response, sourceData datastruct.AudioItems) {
	ms.viewer.openSongMenu(p, sourceData.Items[0])
}

func (ms menuService) Main(p dto.Response) {

	// TODO

}

// TODO
//func (ms menuService) newMainMenu(p dto.Param) {
//	respCfg := tgbotapi.MessageConfig{
//		BaseChat: tgbotapi.BaseChat{ChatId: p.ChatId},
//		Text:     "menu",
//	}
//
//	ms.md.tgBuilder.MenuBuild(respCfg, p,
//
//		ms.md.tgBuilder.NewSubMenu(ms.cfg.MainMenu.Text, ms.cfg.MainMenu.CallbackData,
//			ms.md.tgBuilder.NewLineMenuButton(
//				ms.cfg.LastFmBtn.Text,
//				ms.cfg.LastFmBtn.CallbackData,
//				func(p dto.Param) {
//					ms.viewer.openMenuWithControls(ms.md.getLastFMSimilar(ms.md.getVKSimilar(p)), p)
//				}),
//
//			ms.md.tgBuilder.NewSubMenu(
//				ms.cfg.VkSubMenu.Self.Text,
//				ms.cfg.VkSubMenu.Self.CallbackData,
//				ms.openVkSubmenu(p)...),
//
//			ms.md.tgBuilder.NewMenuButton(
//				ms.cfg.YaMusicBtn.Text,
//				ms.cfg.YaMusicBtn.CallbackData,
//				func(p dto.Param) {
//					ms.viewer.openMenuWithControls(ms.md.getYaMusicSimilar(ms.md.getVKSimilar(p)), p)
//				}),
//		),
//	)
//}
//
//func (ms menuService) openVkSubmenu(p dto.Param) []tg2.Button {
//	if ms.isVkAuthorized(p) {
//		return ms.newVkSubmenu(p)
//	}
//
//	return ms.newAuthVkSubmenu()
//}
//
//func (ms menuService) isVkAuthorized(p dto.Param) bool {
//	return ms.api.IVk.Auth().IsAuthorized(dto.TGUser{
//		UserId: p.UserId,
//		ChatId: p.ChatId,
//	})
//}
//
//func (ms menuService) newVkSubmenu(p dto.Param) []tg2.Button {
//	return []tg2.Button{
//
//		ms.builder.NewMenuButton(
//			ms.cfg.VkSubMenu.UserPlaylist.Text,
//			ms.cfg.VkSubMenu.UserPlaylist.CallbackData,
//			func(p dto.Param) {
//				ms.showPlaylists(aux.tr, callbackData,
//					paginatePlaylistItems(md.GetVKPlaylists(aux.tr, p.MsgId), aux.pageCapacity),
//					p.MsgId, aux.buttonBack)
//			}),
//
//
//		ms.builder.NewMenuButton(
//			ms.cfg.VkSubMenu.Recommendation.Text,
//			ms.cfg.VkSubMenu.Recommendation.CallbackData,
//			func(p dto.Param) {
//				ms.showAudioList(aux.tr, callbackData,
//					paginateAudioItems(md.GetVKSimilars(aux.tr, interMsgId), aux.pageCapacity),
//					interMsgId, aux.buttonBack)
//			}),
//	}
//}
//
//func (ms menuService) newAuthVkSubmenu() []tg2.Button {
//
//	return []tg2.Button{
//		ms.builder.NewMenuButton(
//			ms.cfg.VkSubMenu.Auth.Text,
//			ms.cfg.VkSubMenu.Auth.CallbackData,
//			func(p dto.Param) {
//				respCfg := tgbotapi.MessageConfig{
//					BaseChat: tgbotapi.BaseChat{ChatId: p.ChatId},
//					Text:     "follow the [link](" + ms.api.IVk.Auth().GetAuthURL() + ")" + ", allow access and send an accessToken by message.",
//					ParseMode: "Markdown",
//				}
//				ms.builder.MenuBuild(respCfg, p, aux.buttonBack(false, nil, false))
//			}),
//		aux.buttonBack(false, nil, false),
//	}
//}
